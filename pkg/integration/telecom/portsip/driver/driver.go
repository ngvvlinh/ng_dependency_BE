package driver

import (
	"context"
	"time"

	"o.o/api/etelecom/call_state"
	telecomtypes "o.o/backend/com/etelecom/provider/types"
	cm "o.o/backend/pkg/common"
	portsipclient "o.o/backend/pkg/integration/telecom/portsip/client"
)

var _ telecomtypes.TelecomDriver = &VHTDriver{}

type VHTDriver struct {
	client *portsipclient.Client
}

func New(env string, cfg portsipclient.VHTAccountCfg) *VHTDriver {
	client := portsipclient.New(env, cfg)
	return &VHTDriver{
		client: client,
	}
}

func (v *VHTDriver) GetClient() *portsipclient.Client {
	return v.client
}

func (v *VHTDriver) Ping(ctx context.Context) error {
	return nil
}

func (v *VHTDriver) GenerateToken(ctx context.Context) (*telecomtypes.GenerateTokenResponse, error) {
	loginResp, err := v.client.Login(ctx)
	if err != nil {
		return nil, err
	}

	token := loginResp.AccessToken.String()
	expiresIn := loginResp.Expires.Int()
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	v.client.UpdateToken(token)

	return &telecomtypes.GenerateTokenResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		ExpiresIn:   expiresIn,
	}, nil
}

func (v *VHTDriver) CreateExtension(ctx context.Context, req *telecomtypes.CreateExtensionRequest) (*telecomtypes.CreateExtensionResponse, error) {
	if req.ExtensionPassword == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "password cannot be empty")
	}
	if req.Hotline == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "hotline cannot be empty")
	}
	if req.ExtensionNumber == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Extension number cannot be empty")
	}

	createExtensionReq := &portsipclient.CreateExtensionsRequest{
		ExtensionNumber:   req.ExtensionNumber,
		Password:          req.ExtensionPassword,
		WebAccessPassword: req.ExtensionPassword,
		Options: &portsipclient.OptionsCreateExtension{
			EnableAudioRecordCalls: true,
			EnableVideoRecordCalls: false,
			EnableExtension:        true,
			OutboundCallerID:       req.Hotline,
		},
		ForwardRules: &portsipclient.ForwardRulesCreateExtension{
			Available: &portsipclient.AvailableForwardRules{
				NoAnswerTimeval:     20,
				NoAnswerAction:      "CONNECT_TO_VOICE_MAIL",
				NoAnswerActionValue: "",
				BusyAction:          "CONNECT_TO_VOICE",
				BusyActionValue:     "",
			},
			Offline: &portsipclient.OfflineForwardRules{
				OfficeHoursAction:             "CONNECT_TO_VOICE",
				OfficeHoursActionValue:        "",
				OutsideOfficeHoursAction:      "CONNECT_TO_VOICE",
				OutsideOfficeHoursActionValue: "",
			},
			Dnd: &portsipclient.DndForwardRules{
				OfficeHoursAction:             "CONNECT_TO_VOICE",
				OfficeHoursActionValue:        "",
				OutsideOfficeHoursAction:      "CONNECT_TO_VOICE",
				OutsideOfficeHoursActionValue: "",
			},
			Away: &portsipclient.AwayForwardRules{
				OfficeHoursAction:             "CONNECT_TO_VOICE",
				OfficeHoursActionValue:        "",
				OutsideOfficeHoursAction:      "CONNECT_TO_VOICE",
				OutsideOfficeHoursActionValue: "",
			},
		},
	}
	if req.Profile != nil {
		createExtensionReq.Profile = &portsipclient.ProfileCreateExtension{
			FirstName:   req.Profile.FirstName,
			LastName:    req.Profile.LastName,
			Email:       req.Profile.Email,
			MobilePhone: req.Profile.Phone,
			Description: req.Profile.Description,
		}
	}

	createExtensionResp, err := v.client.CreateExtension(ctx, createExtensionReq)
	if err != nil {
		return nil, err
	}

	return &telecomtypes.CreateExtensionResponse{
		ID: createExtensionResp.ID.String(),
	}, nil
}

func (v *VHTDriver) GetCallLogs(ctx context.Context, req *telecomtypes.GetCallLogsRequest) (res *telecomtypes.GetCallLogsResponse, _ error) {
	getCallLogReq := &portsipclient.GetCallLogsRequest{
		ScrollID: req.ScrollID,
	}
	if !req.StartedAt.IsZero() {
		getCallLogReq.StartTime = req.StartedAt.Unix()
	}
	if !req.EndedAt.IsZero() {
		getCallLogReq.EndTime = req.EndedAt.Unix()
	}

	getCallLogsResp, err := v.client.GetCallLogs(ctx, getCallLogReq)
	if err != nil {
		return nil, err
	}

	res = &telecomtypes.GetCallLogsResponse{
		ScrollID: getCallLogsResp.ScrollID.String(),
	}

	for _, callLog := range getCallLogsResp.Sessions {
		tartgets := getCallTargets(callLog)
		callState := callLog.CallStatus.ToCallState()
		hotlineNumber := callLog.DidCid
		if hotlineNumber == "" {
			hotlineNumber = callLog.OutboundCallerID
		}
		callLogRes := &telecomtypes.CallLog{
			CallID:        callLog.CallID.String(),
			CallStatus:    string(callLog.CallStatus),
			Caller:        callLog.Caller.String(),
			Callee:        callLog.Callee.String(),
			Direction:     callLog.Direction.String(),
			StartedAt:     callLog.StartTime.ToTime(),
			EndedAt:       callLog.EndedTime.ToTime(),
			Duration:      callLog.TalkDuration.Int(),
			CallTargets:   tartgets,
			CallState:     callState,
			HotlineNumber: hotlineNumber.String(),
			SessionID:     callLog.SessionID.String(),
		}
		callLogRes.AudioURLs = append(callLogRes.AudioURLs, callLog.RecordingFileURL.String())
		res.CallLogs = append(res.CallLogs, callLogRes)
	}

	return res, nil
}

func getCallTargets(session *portsipclient.SessionCallLog) (res []*telecomtypes.CallTarget) {
	callTargetsMap := make(map[string]*telecomtypes.CallTarget)

	for _, target := range session.CallTargets {
		targetNumber := target.TargetNumber.String()
		callState := target.Status.ToCallState()

		if _, ok := callTargetsMap[targetNumber]; !ok || callState == call_state.Answered {
			// make sure target number is unique
			callTargetsMap[targetNumber] = &telecomtypes.CallTarget{
				TargetNumber: targetNumber,
				TalkDuration: target.TalkDuration.Int(),
				CallState:    callState,
				AnsweredTime: target.AnsweredTime.ToTime(),
				EndedTime:    target.EndedTime.ToTime(),
			}
		}
	}

	for _, target := range callTargetsMap {
		res = append(res, target)
	}
	return res
}
