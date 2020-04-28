package convert

import (
	"o.o/api/supporting/crm/vht"
	"o.o/backend/com/supporting/crm/vht/model"
)

func ConvertToModel(v *vht.VhtCallLog) *model.VhtCallHistory {
	return &model.VhtCallHistory{
		CdrID:           v.CdrID,
		CallID:          v.CallID,
		SipCallID:       v.SipCallID,
		SdkCallID:       v.SdkCallID,
		Cause:           v.Cause,
		Q850Cause:       v.Q850Cause,
		FromExtension:   v.FromExtension,
		ToExtension:     v.ToExtension,
		FromNumber:      v.FromNumber,
		ToNumber:        v.ToNumber,
		Duration:        v.Duration,
		Direction:       v.Direction,
		TimeStarted:     v.TimeStarted,
		TimeConnected:   v.TimeConnected,
		TimeEnded:       v.TimeEnded,
		RecordingPath:   v.RecordingPath,
		RecordingURL:    v.RecordingUrl,
		RecordFileSize:  v.RecordFileSize,
		EtopAccountID:   v.EtopAccountID,
		VtigerAccountID: v.VtigerAccountID,
	}
}

func ConvertFromModel(m *model.VhtCallHistory) *vht.VhtCallLog {
	return &vht.VhtCallLog{
		CdrID:           m.CdrID,
		CallID:          m.CallID,
		SipCallID:       m.SipCallID,
		SdkCallID:       m.SdkCallID,
		Cause:           m.Cause,
		Q850Cause:       m.Q850Cause,
		FromExtension:   m.FromExtension,
		ToExtension:     m.ToExtension,
		FromNumber:      m.FromNumber,
		ToNumber:        m.ToNumber,
		Duration:        m.Duration,
		Direction:       m.Direction,
		TimeStarted:     m.TimeStarted,
		TimeConnected:   m.TimeConnected,
		TimeEnded:       m.TimeEnded,
		RecordingPath:   m.RecordingPath,
		RecordingUrl:    m.RecordingURL,
		RecordFileSize:  m.RecordFileSize,
		EtopAccountID:   m.EtopAccountID,
		VtigerAccountID: m.VtigerAccountID,
	}
}
