package service

import (
	"etop.vn/backend/pb/common"
	"etop.vn/backend/pb/services/crmservice"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/crm-service/model"
	"etop.vn/backend/pkg/services/crm-service/sqlstore"
)

// Config represents configuration for vht service
type Config struct {
	ServiceURL string `yaml:"service_url"`
	UserName   string `yaml:"user_name"`
	PassWord   string `yaml:"pass_word"`
}

type VhtService struct {
	vhtCallHistories sqlstore.VhtCallHistoriesFactory
}

func NewVhtService(db cmsql.Database) *VhtService {
	return &VhtService{
		vhtCallHistories: sqlstore.NewVhtCallHistoryStore(db),
	}
}

func (s *VhtService) ConvertProto2Model(vht *crmservice.VHTCallLog) *model.VhtCallHistory {
	return &model.VhtCallHistory{
		CdrID:           vht.CdrId,
		CallID:          vht.CallId,
		SipCallID:       vht.SipCallId,
		SdkCallID:       vht.SdkCallId,
		Cause:           vht.Cause,
		Q850Cause:       vht.Q850Cause,
		FromExtension:   vht.FromExtension,
		ToExtension:     vht.ToExtension,
		FromNumber:      vht.FromNumber,
		ToNumber:        vht.ToNumber,
		Duration:        vht.Duration,
		Direction:       vht.Direction,
		TimeStarted:     common.PbTimeToModel(vht.TimeStarted),
		TimeConnected:   common.PbTimeToModel(vht.TimeConnected),
		TimeEnded:       common.PbTimeToModel(vht.TimeEnded),
		RecordingPath:   vht.RecordingPath,
		RecordingURL:    vht.RecordingUrl,
		RecordFileSize:  vht.RecordFileSize,
		EtopAccountID:   vht.EtopAccountId,
		VtigerAccountID: vht.VtigerAccountId,
	}
}

func (s *VhtService) ConvertModel2Proto(vht *model.VhtCallHistory) *crmservice.VHTCallLog {
	return &crmservice.VHTCallLog{
		CdrId:           vht.CdrID,
		CallId:          vht.CallID,
		SipCallId:       vht.SipCallID,
		SdkCallId:       vht.SdkCallID,
		Cause:           vht.Cause,
		Q850Cause:       vht.Q850Cause,
		FromExtension:   vht.FromExtension,
		ToExtension:     vht.ToExtension,
		FromNumber:      vht.FromNumber,
		ToNumber:        vht.ToNumber,
		Duration:        vht.Duration,
		Direction:       vht.Direction,
		TimeStarted:     common.PbTime(vht.TimeStarted),
		TimeConnected:   common.PbTime(vht.TimeConnected),
		TimeEnded:       common.PbTime(vht.TimeEnded),
		RecordingPath:   vht.RecordingPath,
		RecordingUrl:    vht.RecordingURL,
		RecordFileSize:  vht.RecordFileSize,
		EtopAccountId:   vht.EtopAccountID,
		VtigerAccountId: vht.VtigerAccountID,
	}
}
