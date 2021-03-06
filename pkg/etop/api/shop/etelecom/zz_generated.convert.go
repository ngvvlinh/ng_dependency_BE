// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package etelecom

import (
	time "time"

	etelecom "o.o/api/etelecom"
	etelecomtypes "o.o/api/top/int/etelecom/types"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*etelecom.CallLog)(nil), (*etelecomtypes.CallLog)(nil), func(arg, out interface{}) error {
		Convert_etelecom_CallLog_etelecomtypes_CallLog(arg.(*etelecom.CallLog), out.(*etelecomtypes.CallLog))
		return nil
	})
	s.Register(([]*etelecom.CallLog)(nil), (*[]*etelecomtypes.CallLog)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_CallLogs_etelecomtypes_CallLogs(arg.([]*etelecom.CallLog))
		*out.(*[]*etelecomtypes.CallLog) = out0
		return nil
	})
	s.Register((*etelecomtypes.CallLog)(nil), (*etelecom.CallLog)(nil), func(arg, out interface{}) error {
		Convert_etelecomtypes_CallLog_etelecom_CallLog(arg.(*etelecomtypes.CallLog), out.(*etelecom.CallLog))
		return nil
	})
	s.Register(([]*etelecomtypes.CallLog)(nil), (*[]*etelecom.CallLog)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecomtypes_CallLogs_etelecom_CallLogs(arg.([]*etelecomtypes.CallLog))
		*out.(*[]*etelecom.CallLog) = out0
		return nil
	})
	s.Register((*etelecom.CreateOrUpdateCallLogFromCDRArgs)(nil), (*etelecomtypes.CallLog)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateOrUpdateCallLogFromCDRArgs_etelecomtypes_CallLog(arg.(*etelecom.CreateOrUpdateCallLogFromCDRArgs), out.(*etelecomtypes.CallLog))
		return nil
	})
	s.Register((*etelecom.CallTarget)(nil), (*etelecomtypes.CallTarget)(nil), func(arg, out interface{}) error {
		Convert_etelecom_CallTarget_etelecomtypes_CallTarget(arg.(*etelecom.CallTarget), out.(*etelecomtypes.CallTarget))
		return nil
	})
	s.Register(([]*etelecom.CallTarget)(nil), (*[]*etelecomtypes.CallTarget)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_CallTargets_etelecomtypes_CallTargets(arg.([]*etelecom.CallTarget))
		*out.(*[]*etelecomtypes.CallTarget) = out0
		return nil
	})
	s.Register((*etelecomtypes.CallTarget)(nil), (*etelecom.CallTarget)(nil), func(arg, out interface{}) error {
		Convert_etelecomtypes_CallTarget_etelecom_CallTarget(arg.(*etelecomtypes.CallTarget), out.(*etelecom.CallTarget))
		return nil
	})
	s.Register(([]*etelecomtypes.CallTarget)(nil), (*[]*etelecom.CallTarget)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecomtypes_CallTargets_etelecom_CallTargets(arg.([]*etelecomtypes.CallTarget))
		*out.(*[]*etelecom.CallTarget) = out0
		return nil
	})
	s.Register((*etelecom.CreateExtensionArgs)(nil), (*etelecomtypes.Extension)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateExtensionArgs_etelecomtypes_Extension(arg.(*etelecom.CreateExtensionArgs), out.(*etelecomtypes.Extension))
		return nil
	})
	s.Register((*etelecom.Extension)(nil), (*etelecomtypes.Extension)(nil), func(arg, out interface{}) error {
		Convert_etelecom_Extension_etelecomtypes_Extension(arg.(*etelecom.Extension), out.(*etelecomtypes.Extension))
		return nil
	})
	s.Register(([]*etelecom.Extension)(nil), (*[]*etelecomtypes.Extension)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_Extensions_etelecomtypes_Extensions(arg.([]*etelecom.Extension))
		*out.(*[]*etelecomtypes.Extension) = out0
		return nil
	})
	s.Register((*etelecomtypes.Extension)(nil), (*etelecom.Extension)(nil), func(arg, out interface{}) error {
		Convert_etelecomtypes_Extension_etelecom_Extension(arg.(*etelecomtypes.Extension), out.(*etelecom.Extension))
		return nil
	})
	s.Register(([]*etelecomtypes.Extension)(nil), (*[]*etelecom.Extension)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecomtypes_Extensions_etelecom_Extensions(arg.([]*etelecomtypes.Extension))
		*out.(*[]*etelecom.Extension) = out0
		return nil
	})
	s.Register((*etelecom.ExtensionExternalData)(nil), (*etelecomtypes.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		Convert_etelecom_ExtensionExternalData_etelecomtypes_ExtensionExternalData(arg.(*etelecom.ExtensionExternalData), out.(*etelecomtypes.ExtensionExternalData))
		return nil
	})
	s.Register(([]*etelecom.ExtensionExternalData)(nil), (*[]*etelecomtypes.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_ExtensionExternalDatas_etelecomtypes_ExtensionExternalDatas(arg.([]*etelecom.ExtensionExternalData))
		*out.(*[]*etelecomtypes.ExtensionExternalData) = out0
		return nil
	})
	s.Register((*etelecomtypes.ExtensionExternalData)(nil), (*etelecom.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		Convert_etelecomtypes_ExtensionExternalData_etelecom_ExtensionExternalData(arg.(*etelecomtypes.ExtensionExternalData), out.(*etelecom.ExtensionExternalData))
		return nil
	})
	s.Register(([]*etelecomtypes.ExtensionExternalData)(nil), (*[]*etelecom.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecomtypes_ExtensionExternalDatas_etelecom_ExtensionExternalDatas(arg.([]*etelecomtypes.ExtensionExternalData))
		*out.(*[]*etelecom.ExtensionExternalData) = out0
		return nil
	})
	s.Register((*etelecom.CreateHotlineArgs)(nil), (*etelecomtypes.Hotline)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateHotlineArgs_etelecomtypes_Hotline(arg.(*etelecom.CreateHotlineArgs), out.(*etelecomtypes.Hotline))
		return nil
	})
	s.Register((*etelecom.Hotline)(nil), (*etelecomtypes.Hotline)(nil), func(arg, out interface{}) error {
		Convert_etelecom_Hotline_etelecomtypes_Hotline(arg.(*etelecom.Hotline), out.(*etelecomtypes.Hotline))
		return nil
	})
	s.Register(([]*etelecom.Hotline)(nil), (*[]*etelecomtypes.Hotline)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_Hotlines_etelecomtypes_Hotlines(arg.([]*etelecom.Hotline))
		*out.(*[]*etelecomtypes.Hotline) = out0
		return nil
	})
	s.Register((*etelecomtypes.Hotline)(nil), (*etelecom.Hotline)(nil), func(arg, out interface{}) error {
		Convert_etelecomtypes_Hotline_etelecom_Hotline(arg.(*etelecomtypes.Hotline), out.(*etelecom.Hotline))
		return nil
	})
	s.Register(([]*etelecomtypes.Hotline)(nil), (*[]*etelecom.Hotline)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecomtypes_Hotlines_etelecom_Hotlines(arg.([]*etelecomtypes.Hotline))
		*out.(*[]*etelecom.Hotline) = out0
		return nil
	})
	s.Register((*etelecom.Tenant)(nil), (*etelecomtypes.Tenant)(nil), func(arg, out interface{}) error {
		Convert_etelecom_Tenant_etelecomtypes_Tenant(arg.(*etelecom.Tenant), out.(*etelecomtypes.Tenant))
		return nil
	})
	s.Register(([]*etelecom.Tenant)(nil), (*[]*etelecomtypes.Tenant)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_Tenants_etelecomtypes_Tenants(arg.([]*etelecom.Tenant))
		*out.(*[]*etelecomtypes.Tenant) = out0
		return nil
	})
	s.Register((*etelecomtypes.Tenant)(nil), (*etelecom.Tenant)(nil), func(arg, out interface{}) error {
		Convert_etelecomtypes_Tenant_etelecom_Tenant(arg.(*etelecomtypes.Tenant), out.(*etelecom.Tenant))
		return nil
	})
	s.Register(([]*etelecomtypes.Tenant)(nil), (*[]*etelecom.Tenant)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecomtypes_Tenants_etelecom_Tenants(arg.([]*etelecomtypes.Tenant))
		*out.(*[]*etelecom.Tenant) = out0
		return nil
	})
}

//-- convert o.o/api/top/int/etelecom/types.CallLog --//

func Convert_etelecom_CallLog_etelecomtypes_CallLog(arg *etelecom.CallLog, out *etelecomtypes.CallLog) *etelecomtypes.CallLog {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.CallLog{}
	}
	convert_etelecom_CallLog_etelecomtypes_CallLog(arg, out)
	return out
}

func convert_etelecom_CallLog_etelecomtypes_CallLog(arg *etelecom.CallLog, out *etelecomtypes.CallLog) {
	out.ID = arg.ID                                 // simple assign
	out.ExternalID = arg.ExternalID                 // simple assign
	out.AccountID = arg.AccountID                   // simple assign
	out.HotlineID = arg.HotlineID                   // simple assign
	out.OwnerID = arg.OwnerID                       // simple assign
	out.UserID = arg.UserID                         // simple assign
	out.StartedAt = arg.StartedAt                   // simple assign
	out.EndedAt = arg.EndedAt                       // simple assign
	out.Duration = arg.Duration                     // simple assign
	out.Caller = arg.Caller                         // simple assign
	out.Callee = arg.Callee                         // simple assign
	out.AudioURLs = arg.AudioURLs                   // simple assign
	out.ExternalDirection = arg.ExternalDirection   // simple assign
	out.Direction = arg.Direction                   // simple assign
	out.ExtensionID = arg.ExtensionID               // simple assign
	out.ExternalCallStatus = arg.ExternalCallStatus // simple assign
	out.ContactID = arg.ContactID                   // simple assign
	out.CreatedAt = arg.CreatedAt                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                   // simple assign
	out.CallState = arg.CallState                   // simple assign
	out.CallStatus = arg.CallStatus                 // simple assign
	out.DurationPostage = arg.DurationPostage       // simple assign
	out.Postage = arg.Postage                       // simple assign
	out.Note = arg.Note                             // simple assign
	out.CallTargets = Convert_etelecom_CallTargets_etelecomtypes_CallTargets(arg.CallTargets)
}

func Convert_etelecom_CallLogs_etelecomtypes_CallLogs(args []*etelecom.CallLog) (outs []*etelecomtypes.CallLog) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecomtypes.CallLog, len(args))
	outs = make([]*etelecomtypes.CallLog, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_CallLog_etelecomtypes_CallLog(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecomtypes_CallLog_etelecom_CallLog(arg *etelecomtypes.CallLog, out *etelecom.CallLog) *etelecom.CallLog {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.CallLog{}
	}
	convert_etelecomtypes_CallLog_etelecom_CallLog(arg, out)
	return out
}

func convert_etelecomtypes_CallLog_etelecom_CallLog(arg *etelecomtypes.CallLog, out *etelecom.CallLog) {
	out.ID = arg.ID                                 // simple assign
	out.ExternalID = arg.ExternalID                 // simple assign
	out.AccountID = arg.AccountID                   // simple assign
	out.OwnerID = arg.OwnerID                       // simple assign
	out.UserID = arg.UserID                         // simple assign
	out.StartedAt = arg.StartedAt                   // simple assign
	out.EndedAt = arg.EndedAt                       // simple assign
	out.Duration = arg.Duration                     // simple assign
	out.Caller = arg.Caller                         // simple assign
	out.Callee = arg.Callee                         // simple assign
	out.AudioURLs = arg.AudioURLs                   // simple assign
	out.ExternalDirection = arg.ExternalDirection   // simple assign
	out.ExternalCallStatus = arg.ExternalCallStatus // simple assign
	out.CallState = arg.CallState                   // simple assign
	out.CallStatus = arg.CallStatus                 // simple assign
	out.Direction = arg.Direction                   // simple assign
	out.ExtensionID = arg.ExtensionID               // simple assign
	out.HotlineID = arg.HotlineID                   // simple assign
	out.ContactID = arg.ContactID                   // simple assign
	out.CreatedAt = arg.CreatedAt                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                   // simple assign
	out.DurationPostage = arg.DurationPostage       // simple assign
	out.Postage = arg.Postage                       // simple assign
	out.ExternalSessionID = ""                      // zero value
	out.Note = arg.Note                             // simple assign
	out.CallTargets = Convert_etelecomtypes_CallTargets_etelecom_CallTargets(arg.CallTargets)
}

func Convert_etelecomtypes_CallLogs_etelecom_CallLogs(args []*etelecomtypes.CallLog) (outs []*etelecom.CallLog) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.CallLog, len(args))
	outs = make([]*etelecom.CallLog, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecomtypes_CallLog_etelecom_CallLog(args[i], &tmps[i])
	}
	return outs
}

func Apply_etelecom_CreateOrUpdateCallLogFromCDRArgs_etelecomtypes_CallLog(arg *etelecom.CreateOrUpdateCallLogFromCDRArgs, out *etelecomtypes.CallLog) *etelecomtypes.CallLog {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.CallLog{}
	}
	apply_etelecom_CreateOrUpdateCallLogFromCDRArgs_etelecomtypes_CallLog(arg, out)
	return out
}

func apply_etelecom_CreateOrUpdateCallLogFromCDRArgs_etelecomtypes_CallLog(arg *etelecom.CreateOrUpdateCallLogFromCDRArgs, out *etelecomtypes.CallLog) {
	out.ID = 0                                      // zero value
	out.ExternalID = arg.ExternalID                 // simple assign
	out.AccountID = 0                               // zero value
	out.HotlineID = arg.HotlineID                   // simple assign
	out.OwnerID = arg.OwnerID                       // simple assign
	out.UserID = arg.UserID                         // simple assign
	out.StartedAt = arg.StartedAt                   // simple assign
	out.EndedAt = arg.EndedAt                       // simple assign
	out.Duration = arg.Duration                     // simple assign
	out.Caller = arg.Caller                         // simple assign
	out.Callee = arg.Callee                         // simple assign
	out.AudioURLs = arg.AudioURLs                   // simple assign
	out.ExternalDirection = arg.ExternalDirection   // simple assign
	out.Direction = arg.Direction                   // simple assign
	out.ExtensionID = arg.ExtensionID               // simple assign
	out.ExternalCallStatus = arg.ExternalCallStatus // simple assign
	out.ContactID = 0                               // zero value
	out.CreatedAt = time.Time{}                     // zero value
	out.UpdatedAt = time.Time{}                     // zero value
	out.CallState = arg.CallState                   // simple assign
	out.CallStatus = 0                              // zero value
	out.DurationPostage = 0                         // zero value
	out.Postage = 0                                 // zero value
	out.Note = ""                                   // zero value
	out.CallTargets = Convert_etelecom_CallTargets_etelecomtypes_CallTargets(arg.CallTargets)
}

//-- convert o.o/api/top/int/etelecom/types.CallTarget --//

func Convert_etelecom_CallTarget_etelecomtypes_CallTarget(arg *etelecom.CallTarget, out *etelecomtypes.CallTarget) *etelecomtypes.CallTarget {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.CallTarget{}
	}
	convert_etelecom_CallTarget_etelecomtypes_CallTarget(arg, out)
	return out
}

func convert_etelecom_CallTarget_etelecomtypes_CallTarget(arg *etelecom.CallTarget, out *etelecomtypes.CallTarget) {
	out.AddTime = arg.AddTime           // simple assign
	out.AnsweredTime = arg.AnsweredTime // simple assign
	out.EndReason = arg.EndReason       // simple assign
	out.EndedTime = arg.EndedTime       // simple assign
	out.FailCode = arg.FailCode         // simple assign
	out.RingDuration = arg.RingDuration // simple assign
	out.RingTime = arg.RingTime         // simple assign
	out.Status = arg.Status             // simple assign
	out.TargetNumber = arg.TargetNumber // simple assign
	out.TrunkName = arg.TrunkName       // simple assign
}

func Convert_etelecom_CallTargets_etelecomtypes_CallTargets(args []*etelecom.CallTarget) (outs []*etelecomtypes.CallTarget) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecomtypes.CallTarget, len(args))
	outs = make([]*etelecomtypes.CallTarget, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_CallTarget_etelecomtypes_CallTarget(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecomtypes_CallTarget_etelecom_CallTarget(arg *etelecomtypes.CallTarget, out *etelecom.CallTarget) *etelecom.CallTarget {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.CallTarget{}
	}
	convert_etelecomtypes_CallTarget_etelecom_CallTarget(arg, out)
	return out
}

func convert_etelecomtypes_CallTarget_etelecom_CallTarget(arg *etelecomtypes.CallTarget, out *etelecom.CallTarget) {
	out.AddTime = arg.AddTime           // simple assign
	out.AnsweredTime = arg.AnsweredTime // simple assign
	out.EndReason = arg.EndReason       // simple assign
	out.EndedTime = arg.EndedTime       // simple assign
	out.FailCode = arg.FailCode         // simple assign
	out.RingDuration = arg.RingDuration // simple assign
	out.RingTime = arg.RingTime         // simple assign
	out.Status = arg.Status             // simple assign
	out.TargetNumber = arg.TargetNumber // simple assign
	out.TrunkName = arg.TrunkName       // simple assign
}

func Convert_etelecomtypes_CallTargets_etelecom_CallTargets(args []*etelecomtypes.CallTarget) (outs []*etelecom.CallTarget) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.CallTarget, len(args))
	outs = make([]*etelecom.CallTarget, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecomtypes_CallTarget_etelecom_CallTarget(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/top/int/etelecom/types.Extension --//

func Apply_etelecom_CreateExtensionArgs_etelecomtypes_Extension(arg *etelecom.CreateExtensionArgs, out *etelecomtypes.Extension) *etelecomtypes.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.Extension{}
	}
	apply_etelecom_CreateExtensionArgs_etelecomtypes_Extension(arg, out)
	return out
}

func apply_etelecom_CreateExtensionArgs_etelecomtypes_Extension(arg *etelecom.CreateExtensionArgs, out *etelecomtypes.Extension) {
	out.ID = 0                                    // zero value
	out.UserID = arg.UserID                       // simple assign
	out.AccountID = arg.AccountID                 // simple assign
	out.ExtensionNumber = arg.ExtensionNumber     // simple assign
	out.ExtensionPassword = arg.ExtensionPassword // simple assign
	out.TenantID = 0                              // zero value
	out.TenantDomain = ""                         // zero value
	out.HotlineID = arg.HotlineID                 // simple assign
	out.CreatedAt = time.Time{}                   // zero value
	out.UpdatedAt = time.Time{}                   // zero value
	out.ExpiresAt = arg.ExpiresAt                 // simple assign
	out.SubscriptionID = arg.SubscriptionID       // simple assign
}

func Convert_etelecom_Extension_etelecomtypes_Extension(arg *etelecom.Extension, out *etelecomtypes.Extension) *etelecomtypes.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.Extension{}
	}
	convert_etelecom_Extension_etelecomtypes_Extension(arg, out)
	return out
}

func convert_etelecom_Extension_etelecomtypes_Extension(arg *etelecom.Extension, out *etelecomtypes.Extension) {
	out.ID = arg.ID                               // simple assign
	out.UserID = arg.UserID                       // simple assign
	out.AccountID = arg.AccountID                 // simple assign
	out.ExtensionNumber = arg.ExtensionNumber     // simple assign
	out.ExtensionPassword = arg.ExtensionPassword // simple assign
	out.TenantID = arg.TenantID                   // simple assign
	out.TenantDomain = arg.TenantDomain           // simple assign
	out.HotlineID = arg.HotlineID                 // simple assign
	out.CreatedAt = arg.CreatedAt                 // simple assign
	out.UpdatedAt = arg.UpdatedAt                 // simple assign
	out.ExpiresAt = arg.ExpiresAt                 // simple assign
	out.SubscriptionID = arg.SubscriptionID       // simple assign
}

func Convert_etelecom_Extensions_etelecomtypes_Extensions(args []*etelecom.Extension) (outs []*etelecomtypes.Extension) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecomtypes.Extension, len(args))
	outs = make([]*etelecomtypes.Extension, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_Extension_etelecomtypes_Extension(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecomtypes_Extension_etelecom_Extension(arg *etelecomtypes.Extension, out *etelecom.Extension) *etelecom.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Extension{}
	}
	convert_etelecomtypes_Extension_etelecom_Extension(arg, out)
	return out
}

func convert_etelecomtypes_Extension_etelecom_Extension(arg *etelecomtypes.Extension, out *etelecom.Extension) {
	out.ID = arg.ID                               // simple assign
	out.UserID = arg.UserID                       // simple assign
	out.AccountID = arg.AccountID                 // simple assign
	out.HotlineID = arg.HotlineID                 // simple assign
	out.ExtensionNumber = arg.ExtensionNumber     // simple assign
	out.ExtensionPassword = arg.ExtensionPassword // simple assign
	out.TenantDomain = arg.TenantDomain           // simple assign
	out.TenantID = arg.TenantID                   // simple assign
	out.ExternalData = nil                        // zero value
	out.CreatedAt = arg.CreatedAt                 // simple assign
	out.UpdatedAt = arg.UpdatedAt                 // simple assign
	out.DeletedAt = time.Time{}                   // zero value
	out.SubscriptionID = arg.SubscriptionID       // simple assign
	out.ExpiresAt = arg.ExpiresAt                 // simple assign
}

func Convert_etelecomtypes_Extensions_etelecom_Extensions(args []*etelecomtypes.Extension) (outs []*etelecom.Extension) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.Extension, len(args))
	outs = make([]*etelecom.Extension, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecomtypes_Extension_etelecom_Extension(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/top/int/etelecom/types.ExtensionExternalData --//

func Convert_etelecom_ExtensionExternalData_etelecomtypes_ExtensionExternalData(arg *etelecom.ExtensionExternalData, out *etelecomtypes.ExtensionExternalData) *etelecomtypes.ExtensionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.ExtensionExternalData{}
	}
	convert_etelecom_ExtensionExternalData_etelecomtypes_ExtensionExternalData(arg, out)
	return out
}

func convert_etelecom_ExtensionExternalData_etelecomtypes_ExtensionExternalData(arg *etelecom.ExtensionExternalData, out *etelecomtypes.ExtensionExternalData) {
	out.ID = 0 // types do not match
}

func Convert_etelecom_ExtensionExternalDatas_etelecomtypes_ExtensionExternalDatas(args []*etelecom.ExtensionExternalData) (outs []*etelecomtypes.ExtensionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecomtypes.ExtensionExternalData, len(args))
	outs = make([]*etelecomtypes.ExtensionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_ExtensionExternalData_etelecomtypes_ExtensionExternalData(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecomtypes_ExtensionExternalData_etelecom_ExtensionExternalData(arg *etelecomtypes.ExtensionExternalData, out *etelecom.ExtensionExternalData) *etelecom.ExtensionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.ExtensionExternalData{}
	}
	convert_etelecomtypes_ExtensionExternalData_etelecom_ExtensionExternalData(arg, out)
	return out
}

func convert_etelecomtypes_ExtensionExternalData_etelecom_ExtensionExternalData(arg *etelecomtypes.ExtensionExternalData, out *etelecom.ExtensionExternalData) {
	out.ID = "" // types do not match
}

func Convert_etelecomtypes_ExtensionExternalDatas_etelecom_ExtensionExternalDatas(args []*etelecomtypes.ExtensionExternalData) (outs []*etelecom.ExtensionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.ExtensionExternalData, len(args))
	outs = make([]*etelecom.ExtensionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecomtypes_ExtensionExternalData_etelecom_ExtensionExternalData(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/top/int/etelecom/types.Hotline --//

func Apply_etelecom_CreateHotlineArgs_etelecomtypes_Hotline(arg *etelecom.CreateHotlineArgs, out *etelecomtypes.Hotline) *etelecomtypes.Hotline {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.Hotline{}
	}
	apply_etelecom_CreateHotlineArgs_etelecomtypes_Hotline(arg, out)
	return out
}

func apply_etelecom_CreateHotlineArgs_etelecomtypes_Hotline(arg *etelecom.CreateHotlineArgs, out *etelecomtypes.Hotline) {
	out.ID = 0                          // zero value
	out.OwnerID = arg.OwnerID           // simple assign
	out.Name = arg.Name                 // simple assign
	out.Hotline = arg.Hotline           // simple assign
	out.Network = arg.Network           // simple assign
	out.ConnectionID = arg.ConnectionID // simple assign
	out.ConnectionMethod = 0            // zero value
	out.CreatedAt = time.Time{}         // zero value
	out.UpdatedAt = time.Time{}         // zero value
	out.Status = arg.Status             // simple assign
	out.Description = arg.Description   // simple assign
	out.IsFreeCharge = arg.IsFreeCharge // simple assign
}

func Convert_etelecom_Hotline_etelecomtypes_Hotline(arg *etelecom.Hotline, out *etelecomtypes.Hotline) *etelecomtypes.Hotline {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.Hotline{}
	}
	convert_etelecom_Hotline_etelecomtypes_Hotline(arg, out)
	return out
}

func convert_etelecom_Hotline_etelecomtypes_Hotline(arg *etelecom.Hotline, out *etelecomtypes.Hotline) {
	out.ID = arg.ID                             // simple assign
	out.OwnerID = arg.OwnerID                   // simple assign
	out.Name = arg.Name                         // simple assign
	out.Hotline = arg.Hotline                   // simple assign
	out.Network = arg.Network                   // simple assign
	out.ConnectionID = arg.ConnectionID         // simple assign
	out.ConnectionMethod = arg.ConnectionMethod // simple assign
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.Status = arg.Status                     // simple assign
	out.Description = arg.Description           // simple assign
	out.IsFreeCharge = arg.IsFreeCharge         // simple assign
}

func Convert_etelecom_Hotlines_etelecomtypes_Hotlines(args []*etelecom.Hotline) (outs []*etelecomtypes.Hotline) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecomtypes.Hotline, len(args))
	outs = make([]*etelecomtypes.Hotline, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_Hotline_etelecomtypes_Hotline(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecomtypes_Hotline_etelecom_Hotline(arg *etelecomtypes.Hotline, out *etelecom.Hotline) *etelecom.Hotline {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Hotline{}
	}
	convert_etelecomtypes_Hotline_etelecom_Hotline(arg, out)
	return out
}

func convert_etelecomtypes_Hotline_etelecom_Hotline(arg *etelecomtypes.Hotline, out *etelecom.Hotline) {
	out.ID = arg.ID                             // simple assign
	out.OwnerID = arg.OwnerID                   // simple assign
	out.Name = arg.Name                         // simple assign
	out.Hotline = arg.Hotline                   // simple assign
	out.Network = arg.Network                   // simple assign
	out.ConnectionID = arg.ConnectionID         // simple assign
	out.ConnectionMethod = arg.ConnectionMethod // simple assign
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.DeletedAt = time.Time{}                 // zero value
	out.Status = arg.Status                     // simple assign
	out.Description = arg.Description           // simple assign
	out.IsFreeCharge = arg.IsFreeCharge         // simple assign
	out.TenantID = 0                            // zero value
}

func Convert_etelecomtypes_Hotlines_etelecom_Hotlines(args []*etelecomtypes.Hotline) (outs []*etelecom.Hotline) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.Hotline, len(args))
	outs = make([]*etelecom.Hotline, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecomtypes_Hotline_etelecom_Hotline(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/top/int/etelecom/types.Tenant --//

func Convert_etelecom_Tenant_etelecomtypes_Tenant(arg *etelecom.Tenant, out *etelecomtypes.Tenant) *etelecomtypes.Tenant {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecomtypes.Tenant{}
	}
	convert_etelecom_Tenant_etelecomtypes_Tenant(arg, out)
	return out
}

func convert_etelecom_Tenant_etelecomtypes_Tenant(arg *etelecom.Tenant, out *etelecomtypes.Tenant) {
	out.ID = arg.ID                             // simple assign
	out.OwnerID = arg.OwnerID                   // simple assign
	out.Name = arg.Name                         // simple assign
	out.Domain = arg.Domain                     // simple assign
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.Status = arg.Status                     // simple assign
	out.ConnectionID = arg.ConnectionID         // simple assign
	out.ConnectionMethod = arg.ConnectionMethod // simple assign
}

func Convert_etelecom_Tenants_etelecomtypes_Tenants(args []*etelecom.Tenant) (outs []*etelecomtypes.Tenant) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecomtypes.Tenant, len(args))
	outs = make([]*etelecomtypes.Tenant, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_Tenant_etelecomtypes_Tenant(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecomtypes_Tenant_etelecom_Tenant(arg *etelecomtypes.Tenant, out *etelecom.Tenant) *etelecom.Tenant {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Tenant{}
	}
	convert_etelecomtypes_Tenant_etelecom_Tenant(arg, out)
	return out
}

func convert_etelecomtypes_Tenant_etelecom_Tenant(arg *etelecomtypes.Tenant, out *etelecom.Tenant) {
	out.ID = arg.ID                             // simple assign
	out.OwnerID = arg.OwnerID                   // simple assign
	out.Name = arg.Name                         // simple assign
	out.Domain = arg.Domain                     // simple assign
	out.Password = ""                           // zero value
	out.ExternalData = nil                      // zero value
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.DeletedAt = time.Time{}                 // zero value
	out.Status = arg.Status                     // simple assign
	out.ConnectionID = arg.ConnectionID         // simple assign
	out.ConnectionMethod = arg.ConnectionMethod // simple assign
}

func Convert_etelecomtypes_Tenants_etelecom_Tenants(args []*etelecomtypes.Tenant) (outs []*etelecom.Tenant) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.Tenant, len(args))
	outs = make([]*etelecom.Tenant, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecomtypes_Tenant_etelecom_Tenant(args[i], &tmps[i])
	}
	return outs
}
