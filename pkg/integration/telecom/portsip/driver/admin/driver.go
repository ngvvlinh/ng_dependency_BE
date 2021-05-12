package admin

import (
	"context"
	"strings"
	"time"

	telecomtypes "o.o/backend/com/etelecom/provider/types"
	cm "o.o/backend/pkg/common"
	portsipclient "o.o/backend/pkg/integration/telecom/portsip/client"
)

var _ telecomtypes.TelecomAdminDriver = &PortsipAdminDriver{}

type PortsipAdminDriver struct {
	client *portsipclient.Client
}

func NewAdminAccount(cfg portsipclient.PortsipAccountCfg) *PortsipAdminDriver {
	client := portsipclient.New(cfg)
	return &PortsipAdminDriver{
		client: client,
	}
}

func (d *PortsipAdminDriver) GenerateToken(ctx context.Context) (*telecomtypes.GenerateTokenResponse, error) {
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

func (d *PortsipAdminDriver) CreateTenant(ctx context.Context, req *telecomtypes.CreateTenantRequest) (*telecomtypes.CreateTenantResponse, error) {
	cmd := &portsipclient.CreateTenantRequest{
		Name:     req.Name,
		Domain:   req.Domain,
		Password: req.Password,
		Enabled:  req.Enable,
		Profile: &portsipclient.TenantProfile{
			FirstName:                     req.Profile.FirstName,
			LastName:                      req.Profile.LastName,
			Email:                         req.Profile.Email,
			Region:                        "Vietnam",
			Timezone:                      "Asia/Ho_Chi_Minh",
			Currency:                      "VND",
			EnableExtensionChangePassword: true,
			EnableExtensionVideoRecording: true,
			EnableExtensionAudioRecording: true,
		},
		Capability: &portsipclient.TenantCapability{
			MaxExtensions:           100,
			MaxConcurrentCalls:      100,
			MaxRingGroups:           10,
			MaxVirtualReceptionists: 10,
			MaxCallQueues:           10,
			MaxConferenceRooms:      10,
		},
		Quota: &portsipclient.TenantQuata{
			MaxRecordingsQuota:       0,
			MaxVoicemailQuota:        0,
			MaxCallReportQuota:       0,
			AutoCleanRecordingsDays:  30,
			AutoCleanVoicemailDays:   30,
			AutoCleanCallReportsDays: 30,
		},
	}
	resp, err := d.client.CreateTenant(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &telecomtypes.CreateTenantResponse{ID: resp.ID.String()}, nil
}

func (d *PortsipAdminDriver) AddHotlineToTenantInTrunkProvider(ctx context.Context, req *telecomtypes.AddHotlineToTenantInTrunkProviderRequest) error {
	if req.TrunkProviderID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing trunk provider ID")
	}
	if req.Hotline == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing hotline number")
	}
	if req.TenantID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing tenant id")
	}

	query := &portsipclient.GetTrunkProviderRequest{
		ID: req.TrunkProviderID,
	}
	trunkProvider, err := d.client.GetTrunkProvider(ctx, query)
	if err != nil {
		return err
	}
	pools := trunkProvider.DidPool

	var updatePools []*portsipclient.TrunkProviderDidPool
	var added bool
	for _, pool := range pools {
		if pool.TenantID == req.TenantID {
			if pool.NumberMask == "" {
				continue
			}
			if strings.Contains(pool.NumberMask, req.Hotline) {
				// tenant & hotline existed in trunk provider
				// do nothing
				return nil
			}
			hotlines := strings.Split(pool.NumberMask, ";")
			hotlines = append(hotlines, req.Hotline)

			updatePools = append(updatePools, &portsipclient.TrunkProviderDidPool{
				TenantID:   req.TenantID,
				NumberMask: strings.Join(hotlines, ";"),
			})
			added = true
		} else {
			updatePools = append(updatePools, pool)
		}
	}

	if !added {
		updatePools = append(updatePools, &portsipclient.TrunkProviderDidPool{
			TenantID:   req.TenantID,
			NumberMask: req.Hotline,
		})
	}

	if len(updatePools) == 0 {
		return nil
	}
	update := &portsipclient.UpdateTrunkProviderRequest{
		ID:      req.TrunkProviderID,
		DidPool: updatePools,
	}
	return d.client.UpdateTrunkProvider(ctx, update)
}

func (d *PortsipAdminDriver) RemoveHotlineOutOfTenantInTrunkProvider(ctx context.Context, req *telecomtypes.RemoveHotlineOutOfTenantInTrunkProviderRequest) error {
	if req.TrunkProviderID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing trunk provider ID")
	}
	if req.Hotline == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing hotline number")
	}
	if req.TenantID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing tenant id")
	}

	query := &portsipclient.GetTrunkProviderRequest{
		ID: req.TrunkProviderID,
	}
	trunkProvider, err := d.client.GetTrunkProvider(ctx, query)
	if err != nil {
		return err
	}
	pools := trunkProvider.DidPool

	var updatePools []*portsipclient.TrunkProviderDidPool
	var removed bool
	for _, pool := range pools {
		if pool.TenantID == req.TenantID {
			// NumberMask is list hotlines, split by semicolon
			if pool.NumberMask == "" {
				return nil
			}

			hotlines := strings.Split(pool.NumberMask, ";")
			newHotlines := make([]string, 0, len(hotlines))
			for _, hotline := range hotlines {
				if hotline == req.Hotline {
					removed = true
					continue
				}
				newHotlines = append(newHotlines, hotline)
			}

			updatePools = append(updatePools, &portsipclient.TrunkProviderDidPool{
				TenantID:   req.TenantID,
				NumberMask: strings.Join(newHotlines, ";"),
			})
		} else {
			updatePools = append(updatePools, pool)
		}
	}

	if !removed || len(updatePools) == 0 {
		return nil
	}
	update := &portsipclient.UpdateTrunkProviderRequest{
		ID:      req.TrunkProviderID,
		DidPool: updatePools,
	}
	return d.client.UpdateTrunkProvider(ctx, update)
}
