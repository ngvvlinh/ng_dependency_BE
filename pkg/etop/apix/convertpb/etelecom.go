package convertpb

import (
	"o.o/api/etelecom"
	externaltypes "o.o/api/top/external/types"
)

func Convert_core_Extension_To_api_ExtensionInfo(in *etelecom.Extension) *externaltypes.ExtensionInfo {
	if in == nil {
		return nil
	}
	out := &externaltypes.ExtensionInfo{
		ExtensionNumber:   in.ExtensionNumber,
		ExtensionPassword: in.ExtensionPassword,
		TenantDomain:      in.TenantDomain,
		ExpiresAt:         in.ExpiresAt,
	}
	return out
}

func Convert_core_Calllog_To_api_ShopCalllog(in *etelecom.CallLog) *externaltypes.ShopCallLog {
	if in == nil {
		return nil
	}
	out := &externaltypes.ShopCallLog{
		ID:                in.ID,
		ExternalSessionID: in.ExternalID,
		UserID:            in.UserID,
		StartedAt:         in.StartedAt,
		EndedAt:           in.EndedAt,
		Duration:          in.Duration,
		Caller:            in.Caller,
		Callee:            in.Callee,
		RecordingURLs:     in.AudioURLs,
		Direction:         in.Direction,
		ExtensionID:       in.ExtensionID,
		ContactID:         in.ContactID,
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         in.UpdatedAt,
		CallState:         in.CallState,
		CallStatus:        in.CallStatus,
		Note:              in.Note,
	}
	return out
}

func Convert_core_Calllogs_To_api_ShopCalllogs(in []*etelecom.CallLog) []*externaltypes.ShopCallLog {
	if in == nil {
		return nil
	}
	outs := make([]*externaltypes.ShopCallLog, len(in))
	for i, v := range in {
		outs[i] = Convert_core_Calllog_To_api_ShopCalllog(v)
	}
	return outs
}
