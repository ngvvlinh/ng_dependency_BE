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

func (v *VHTDriver) GetCallLogs(ctx context.Context) (res []*telecomtypes.CallLog, _ error) {
	getCallLogsRequest, err := v.client.GetCallLogs(ctx)
	if err != nil {
		return nil, err
	}

	for _, callLog := range getCallLogsRequest.Sessions {
		callLogRes := &telecomtypes.CallLog{
			CallID:       callLog.CallID.String(),
			CallStatus:   callLog.CallStatus.String(),
			Callee:       callLog.Callee.String(),
			CalleeDomain: callLog.CalleeDomain.String(),
			StartTime:    callLog.StartTime.String(),
			EndTime:      callLog.EndTime.String(),
			TaskDuration: callLog.TaskDuration.String(),
		}
		for _, audioURL := range callLog.AudioURLs {
			callLogRes.AudioURLs = append(callLogRes.AudioURLs, audioURL.String())
		}
		res = append(res, callLogRes)
	}

	return res, nil
}
