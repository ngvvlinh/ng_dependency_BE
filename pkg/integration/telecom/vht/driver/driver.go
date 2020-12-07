package driver

import (
	"context"
	"time"

	telecomtypes "o.o/backend/com/etelecom/provider/types"
	cm "o.o/backend/pkg/common"
	vhtclient "o.o/backend/pkg/integration/telecom/vht/client"
)

var _ telecomtypes.TelecomDriver = &VHTDriver{}

type VHTDriver struct {
	client *vhtclient.Client
}

func New(env string, cfg vhtclient.VHTAccountCfg) *VHTDriver {
	client := vhtclient.New(env, cfg)
	return &VHTDriver{
		client: client,
	}
}

func (v *VHTDriver) GetClient() *vhtclient.Client {
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

	createExtensionReq := &vhtclient.CreateExtensionsRequest{
		ExtensionNumber:   req.ExtensionNumber,
		Password:          req.ExtensionPassword,
		WebAccessPassword: req.ExtensionPassword,
		Options: &vhtclient.OptionsCreateExtension{
			EnableAudioRecordCalls: true,
			EnableVideoRecordCalls: false,
			EnableExtension:        true,
			OutboundCallerID:       req.Hotline,
		},
		ForwardRules: &vhtclient.ForwardRulesCreateExtension{
			Available: &vhtclient.AvailableForwardRules{
				NoAnswerTimeval:     20,
				NoAnswerAction:      "CONNECT_TO_VOICE_MAIL",
				NoAnswerActionValue: "",
				BusyAction:          "CONNECT_TO_VOICE",
				BusyActionValue:     "",
			},
			Offline: &vhtclient.OfflineForwardRules{
				OfficeHoursAction:             "CONNECT_TO_VOICE",
				OfficeHoursActionValue:        "",
				OutsideOfficeHoursAction:      "CONNECT_TO_VOICE",
				OutsideOfficeHoursActionValue: "",
			},
			Dnd: &vhtclient.DndForwardRules{
				OfficeHoursAction:             "CONNECT_TO_VOICE",
				OfficeHoursActionValue:        "",
				OutsideOfficeHoursAction:      "CONNECT_TO_VOICE",
				OutsideOfficeHoursActionValue: "",
			},
			Away: &vhtclient.AwayForwardRules{
				OfficeHoursAction:             "CONNECT_TO_VOICE",
				OfficeHoursActionValue:        "",
				OutsideOfficeHoursAction:      "CONNECT_TO_VOICE",
				OutsideOfficeHoursActionValue: "",
			},
		},
	}
	if req.Profile != nil {
		createExtensionReq.Profile = &vhtclient.ProfileCreateExtension{
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
	getCallLogReq := &vhtclient.GetCallLogsRequest{
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
		callLogRes := &telecomtypes.CallLog{
			CallID:     callLog.CallID.String(),
			CallStatus: string(callLog.CallStatus),
			Caller:     callLog.Caller.String(),
			Callee:     callLog.Callee.String(),
			Direction:  callLog.Direction.String(),
			StartedAt:  callLog.StartTime.ToTime(),
			EndedAt:    callLog.EndTime.ToTime(),
			Duration:   callLog.TalkDuration.Int(),
		}
		for _, audioURL := range callLog.AudioURLs {
			callLogRes.AudioURLs = append(callLogRes.AudioURLs, audioURL.URL.String())
		}
		res.CallLogs = append(res.CallLogs, callLogRes)
	}

	return res, nil
}
