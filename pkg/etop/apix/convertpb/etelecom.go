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
