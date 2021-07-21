package driver

import (
	"context"
	"time"

	telecomtypes "o.o/backend/com/etelecom/provider/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	portsipclient "o.o/backend/pkg/integration/telecom/portsip/client"
)

var _ telecomtypes.TelecomDriver = &PortsipDriver{}

type PortsipDriver struct {
	client *portsipclient.Client
}

func New(cfg portsipclient.PortsipAccountCfg) *PortsipDriver {
	client := portsipclient.New(cfg)
	return &PortsipDriver{
		client: client,
	}
}

func (d *PortsipDriver) GetClient() *portsipclient.Client {
	return d.client
}

func (d *PortsipDriver) Ping(ctx context.Context) error {
	return nil
}

func (d *PortsipDriver) GenerateToken(ctx context.Context) (*telecomtypes.GenerateTokenResponse, error) {
	loginResp, err := d.client.Login(ctx)
	if err != nil {
		return nil, err
	}

	token := loginResp.AccessToken.String()
	expiresIn := loginResp.Expires.Int()
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	d.client.UpdateToken(token)

	return &telecomtypes.GenerateTokenResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		ExpiresIn:   expiresIn,
	}, nil
}

func (d *PortsipDriver) CreateExtension(ctx context.Context, req *telecomtypes.CreateExtensionRequest) (*telecomtypes.CreateExtensionResponse, error) {
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
			EnableVideoRecordCalls: true,
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
		createExtensionReq.Profile = &portsipclient.ExtensionProfile{
			FirstName:   req.Profile.FirstName,
			LastName:    req.Profile.LastName,
			Email:       req.Profile.Email,
			MobilePhone: req.Profile.Phone,
			Description: req.Profile.Description,
		}
	}

	createExtensionResp, err := d.client.CreateExtension(ctx, createExtensionReq)
	if err != nil {
		return nil, err
	}

	return &telecomtypes.CreateExtensionResponse{
		ID: createExtensionResp.ID.String(),
	}, nil
}

const PortsipDefaultGroupName = "__DEFAULT__"
const OutboundRuleDefaultName = "eB2B Outbound Rule"

func (d *PortsipDriver) CreateOutboundRule(ctx context.Context, args *telecomtypes.CreateOutboundRuleRequest) error {
	reqListOutboundRules := &portsipclient.CommonListRequest{
		Pagination: 1,
		Pagesize:   1000,
	}
	respListOutboundRule, err := d.client.ListOutboundRules(ctx, reqListOutboundRules)
	if err != nil {
		return err
	}
	for _, rule := range respListOutboundRule.Rules {
		if rule.Name == OutboundRuleDefaultName {
			// outbound rule default was created
			return nil
		}
	}

	//  Portsip create outbound rule need:
	//  Extension group default
	reqExtensionGroups := &portsipclient.CommonListRequest{
		Pagination: 1,
		Pagesize:   1000,
	}
	respExtensionGroups, err := d.client.GetExtensionGroups(ctx, reqExtensionGroups)
	if err != nil {
		return err
	}

	tenantExtensionGroupDefaultID := ""
	for _, group := range respExtensionGroups.Groups {
		if group.GroupName == PortsipDefaultGroupName {
			tenantExtensionGroupDefaultID = group.ID.String()
			break
		}
	}
	if tenantExtensionGroupDefaultID == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "Can not create outbound rules. Portsip Extension Group default not found")
	}

	// create outbound
	reqOutboundRule := &portsipclient.CreateOutboundRuleRequest{
		Name:         OutboundRuleDefaultName,
		NumberPrefix: "",
		FromExtensionGroups: &portsipclient.ExtensionGroup{
			ID:        httpreq.String(tenantExtensionGroupDefaultID),
			GroupName: PortsipDefaultGroupName,
		},
		Routes: []*portsipclient.OutboundRuleRoute{
			{
				ID: args.TrunkProviderID,
			},
		},
	}
	_, err = d.client.CreateOutboundRule(ctx, reqOutboundRule)
	if err != nil {
		return err
	}
	return nil
}
