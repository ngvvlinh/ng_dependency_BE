// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	etelecom "o.o/api/etelecom"
	etelecommodel "o.o/backend/com/etelecom/model"
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
	s.Register((*etelecommodel.CallLog)(nil), (*etelecom.CallLog)(nil), func(arg, out interface{}) error {
		Convert_etelecommodel_CallLog_etelecom_CallLog(arg.(*etelecommodel.CallLog), out.(*etelecom.CallLog))
		return nil
	})
	s.Register(([]*etelecommodel.CallLog)(nil), (*[]*etelecom.CallLog)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecommodel_CallLogs_etelecom_CallLogs(arg.([]*etelecommodel.CallLog))
		*out.(*[]*etelecom.CallLog) = out0
		return nil
	})
	s.Register((*etelecom.CallLog)(nil), (*etelecommodel.CallLog)(nil), func(arg, out interface{}) error {
		Convert_etelecom_CallLog_etelecommodel_CallLog(arg.(*etelecom.CallLog), out.(*etelecommodel.CallLog))
		return nil
	})
	s.Register(([]*etelecom.CallLog)(nil), (*[]*etelecommodel.CallLog)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_CallLogs_etelecommodel_CallLogs(arg.([]*etelecom.CallLog))
		*out.(*[]*etelecommodel.CallLog) = out0
		return nil
	})
	s.Register((*etelecom.CreateCallLogFromCDRArgs)(nil), (*etelecom.CallLog)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateCallLogFromCDRArgs_etelecom_CallLog(arg.(*etelecom.CreateCallLogFromCDRArgs), out.(*etelecom.CallLog))
		return nil
	})
	s.Register((*etelecommodel.Extension)(nil), (*etelecom.Extension)(nil), func(arg, out interface{}) error {
		Convert_etelecommodel_Extension_etelecom_Extension(arg.(*etelecommodel.Extension), out.(*etelecom.Extension))
		return nil
	})
	s.Register(([]*etelecommodel.Extension)(nil), (*[]*etelecom.Extension)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecommodel_Extensions_etelecom_Extensions(arg.([]*etelecommodel.Extension))
		*out.(*[]*etelecom.Extension) = out0
		return nil
	})
	s.Register((*etelecom.Extension)(nil), (*etelecommodel.Extension)(nil), func(arg, out interface{}) error {
		Convert_etelecom_Extension_etelecommodel_Extension(arg.(*etelecom.Extension), out.(*etelecommodel.Extension))
		return nil
	})
	s.Register(([]*etelecom.Extension)(nil), (*[]*etelecommodel.Extension)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_Extensions_etelecommodel_Extensions(arg.([]*etelecom.Extension))
		*out.(*[]*etelecommodel.Extension) = out0
		return nil
	})
	s.Register((*etelecom.CreateExtensionArgs)(nil), (*etelecom.Extension)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateExtensionArgs_etelecom_Extension(arg.(*etelecom.CreateExtensionArgs), out.(*etelecom.Extension))
		return nil
	})
	s.Register((*etelecommodel.ExtensionExternalData)(nil), (*etelecom.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		Convert_etelecommodel_ExtensionExternalData_etelecom_ExtensionExternalData(arg.(*etelecommodel.ExtensionExternalData), out.(*etelecom.ExtensionExternalData))
		return nil
	})
	s.Register(([]*etelecommodel.ExtensionExternalData)(nil), (*[]*etelecom.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecommodel_ExtensionExternalDatas_etelecom_ExtensionExternalDatas(arg.([]*etelecommodel.ExtensionExternalData))
		*out.(*[]*etelecom.ExtensionExternalData) = out0
		return nil
	})
	s.Register((*etelecom.ExtensionExternalData)(nil), (*etelecommodel.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		Convert_etelecom_ExtensionExternalData_etelecommodel_ExtensionExternalData(arg.(*etelecom.ExtensionExternalData), out.(*etelecommodel.ExtensionExternalData))
		return nil
	})
	s.Register(([]*etelecom.ExtensionExternalData)(nil), (*[]*etelecommodel.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_ExtensionExternalDatas_etelecommodel_ExtensionExternalDatas(arg.([]*etelecom.ExtensionExternalData))
		*out.(*[]*etelecommodel.ExtensionExternalData) = out0
		return nil
	})
	s.Register((*etelecommodel.Hotline)(nil), (*etelecom.Hotline)(nil), func(arg, out interface{}) error {
		Convert_etelecommodel_Hotline_etelecom_Hotline(arg.(*etelecommodel.Hotline), out.(*etelecom.Hotline))
		return nil
	})
	s.Register(([]*etelecommodel.Hotline)(nil), (*[]*etelecom.Hotline)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecommodel_Hotlines_etelecom_Hotlines(arg.([]*etelecommodel.Hotline))
		*out.(*[]*etelecom.Hotline) = out0
		return nil
	})
	s.Register((*etelecom.Hotline)(nil), (*etelecommodel.Hotline)(nil), func(arg, out interface{}) error {
		Convert_etelecom_Hotline_etelecommodel_Hotline(arg.(*etelecom.Hotline), out.(*etelecommodel.Hotline))
		return nil
	})
	s.Register(([]*etelecom.Hotline)(nil), (*[]*etelecommodel.Hotline)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_Hotlines_etelecommodel_Hotlines(arg.([]*etelecom.Hotline))
		*out.(*[]*etelecommodel.Hotline) = out0
		return nil
	})
	s.Register((*etelecom.CreateHotlineArgs)(nil), (*etelecom.Hotline)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateHotlineArgs_etelecom_Hotline(arg.(*etelecom.CreateHotlineArgs), out.(*etelecom.Hotline))
		return nil
	})
}

//-- convert o.o/api/etelecom.CallLog --//

func Convert_etelecommodel_CallLog_etelecom_CallLog(arg *etelecommodel.CallLog, out *etelecom.CallLog) *etelecom.CallLog {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.CallLog{}
	}
	convert_etelecommodel_CallLog_etelecom_CallLog(arg, out)
	return out
}

func convert_etelecommodel_CallLog_etelecom_CallLog(arg *etelecommodel.CallLog, out *etelecom.CallLog) {
	out.ID = arg.ID                                 // simple assign
	out.ExternalID = arg.ExternalID                 // simple assign
	out.AccountID = arg.AccountID                   // simple assign
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
}

func Convert_etelecommodel_CallLogs_etelecom_CallLogs(args []*etelecommodel.CallLog) (outs []*etelecom.CallLog) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.CallLog, len(args))
	outs = make([]*etelecom.CallLog, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecommodel_CallLog_etelecom_CallLog(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecom_CallLog_etelecommodel_CallLog(arg *etelecom.CallLog, out *etelecommodel.CallLog) *etelecommodel.CallLog {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecommodel.CallLog{}
	}
	convert_etelecom_CallLog_etelecommodel_CallLog(arg, out)
	return out
}

func convert_etelecom_CallLog_etelecommodel_CallLog(arg *etelecom.CallLog, out *etelecommodel.CallLog) {
	out.ID = arg.ID                                 // simple assign
	out.ExternalID = arg.ExternalID                 // simple assign
	out.AccountID = arg.AccountID                   // simple assign
	out.StartedAt = arg.StartedAt                   // simple assign
	out.EndedAt = arg.EndedAt                       // simple assign
	out.Duration = arg.Duration                     // simple assign
	out.Caller = arg.Caller                         // simple assign
	out.Callee = arg.Callee                         // simple assign
	out.AudioURLs = arg.AudioURLs                   // simple assign
	out.ExternalDirection = arg.ExternalDirection   // simple assign
	out.Direction = arg.Direction                   // simple assign
	out.ExtensionID = arg.ExtensionID               // simple assign
	out.HotlineID = arg.HotlineID                   // simple assign
	out.ExternalCallStatus = arg.ExternalCallStatus // simple assign
	out.ContactID = arg.ContactID                   // simple assign
	out.CreatedAt = arg.CreatedAt                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                   // simple assign
	out.CallState = arg.CallState                   // simple assign
	out.CallStatus = arg.CallStatus                 // simple assign
	out.DurationPostage = arg.DurationPostage       // simple assign
	out.Postage = arg.Postage                       // simple assign
}

func Convert_etelecom_CallLogs_etelecommodel_CallLogs(args []*etelecom.CallLog) (outs []*etelecommodel.CallLog) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecommodel.CallLog, len(args))
	outs = make([]*etelecommodel.CallLog, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_CallLog_etelecommodel_CallLog(args[i], &tmps[i])
	}
	return outs
}

func Apply_etelecom_CreateCallLogFromCDRArgs_etelecom_CallLog(arg *etelecom.CreateCallLogFromCDRArgs, out *etelecom.CallLog) *etelecom.CallLog {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.CallLog{}
	}
	apply_etelecom_CreateCallLogFromCDRArgs_etelecom_CallLog(arg, out)
	return out
}

func apply_etelecom_CreateCallLogFromCDRArgs_etelecom_CallLog(arg *etelecom.CreateCallLogFromCDRArgs, out *etelecom.CallLog) {
	out.ID = 0                                      // zero value
	out.ExternalID = arg.ExternalID                 // simple assign
	out.AccountID = 0                               // zero value
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
	out.ExtensionID = 0                             // zero value
	out.HotlineID = 0                               // zero value
	out.ContactID = 0                               // zero value
	out.CreatedAt = time.Time{}                     // zero value
	out.UpdatedAt = time.Time{}                     // zero value
	out.DurationPostage = 0                         // zero value
	out.Postage = 0                                 // zero value
}

//-- convert o.o/api/etelecom.Extension --//

func Convert_etelecommodel_Extension_etelecom_Extension(arg *etelecommodel.Extension, out *etelecom.Extension) *etelecom.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Extension{}
	}
	convert_etelecommodel_Extension_etelecom_Extension(arg, out)
	return out
}

func convert_etelecommodel_Extension_etelecom_Extension(arg *etelecommodel.Extension, out *etelecom.Extension) {
	out.ID = arg.ID                               // simple assign
	out.UserID = arg.UserID                       // simple assign
	out.AccountID = arg.AccountID                 // simple assign
	out.HotlineID = arg.HotlineID                 // simple assign
	out.ExtensionNumber = arg.ExtensionNumber     // simple assign
	out.ExtensionPassword = arg.ExtensionPassword // simple assign
	out.ExternalData = Convert_etelecommodel_ExtensionExternalData_etelecom_ExtensionExternalData(arg.ExternalData, nil)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
}

func Convert_etelecommodel_Extensions_etelecom_Extensions(args []*etelecommodel.Extension) (outs []*etelecom.Extension) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.Extension, len(args))
	outs = make([]*etelecom.Extension, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecommodel_Extension_etelecom_Extension(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecom_Extension_etelecommodel_Extension(arg *etelecom.Extension, out *etelecommodel.Extension) *etelecommodel.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecommodel.Extension{}
	}
	convert_etelecom_Extension_etelecommodel_Extension(arg, out)
	return out
}

func convert_etelecom_Extension_etelecommodel_Extension(arg *etelecom.Extension, out *etelecommodel.Extension) {
	out.ID = arg.ID                               // simple assign
	out.UserID = arg.UserID                       // simple assign
	out.AccountID = arg.AccountID                 // simple assign
	out.HotlineID = arg.HotlineID                 // simple assign
	out.ExtensionNumber = arg.ExtensionNumber     // simple assign
	out.ExtensionPassword = arg.ExtensionPassword // simple assign
	out.ExternalData = Convert_etelecom_ExtensionExternalData_etelecommodel_ExtensionExternalData(arg.ExternalData, nil)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
}

func Convert_etelecom_Extensions_etelecommodel_Extensions(args []*etelecom.Extension) (outs []*etelecommodel.Extension) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecommodel.Extension, len(args))
	outs = make([]*etelecommodel.Extension, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_Extension_etelecommodel_Extension(args[i], &tmps[i])
	}
	return outs
}

func Apply_etelecom_CreateExtensionArgs_etelecom_Extension(arg *etelecom.CreateExtensionArgs, out *etelecom.Extension) *etelecom.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Extension{}
	}
	apply_etelecom_CreateExtensionArgs_etelecom_Extension(arg, out)
	return out
}

func apply_etelecom_CreateExtensionArgs_etelecom_Extension(arg *etelecom.CreateExtensionArgs, out *etelecom.Extension) {
	out.ID = 0                                    // zero value
	out.UserID = arg.UserID                       // simple assign
	out.AccountID = arg.AccountID                 // simple assign
	out.HotlineID = arg.HotlineID                 // simple assign
	out.ExtensionNumber = arg.ExtensionNumber     // simple assign
	out.ExtensionPassword = arg.ExtensionPassword // simple assign
	out.ExternalData = nil                        // zero value
	out.CreatedAt = time.Time{}                   // zero value
	out.UpdatedAt = time.Time{}                   // zero value
	out.DeletedAt = time.Time{}                   // zero value
}

//-- convert o.o/api/etelecom.ExtensionExternalData --//

func Convert_etelecommodel_ExtensionExternalData_etelecom_ExtensionExternalData(arg *etelecommodel.ExtensionExternalData, out *etelecom.ExtensionExternalData) *etelecom.ExtensionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.ExtensionExternalData{}
	}
	convert_etelecommodel_ExtensionExternalData_etelecom_ExtensionExternalData(arg, out)
	return out
}

func convert_etelecommodel_ExtensionExternalData_etelecom_ExtensionExternalData(arg *etelecommodel.ExtensionExternalData, out *etelecom.ExtensionExternalData) {
	out.ID = arg.ID // simple assign
}

func Convert_etelecommodel_ExtensionExternalDatas_etelecom_ExtensionExternalDatas(args []*etelecommodel.ExtensionExternalData) (outs []*etelecom.ExtensionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.ExtensionExternalData, len(args))
	outs = make([]*etelecom.ExtensionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecommodel_ExtensionExternalData_etelecom_ExtensionExternalData(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecom_ExtensionExternalData_etelecommodel_ExtensionExternalData(arg *etelecom.ExtensionExternalData, out *etelecommodel.ExtensionExternalData) *etelecommodel.ExtensionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecommodel.ExtensionExternalData{}
	}
	convert_etelecom_ExtensionExternalData_etelecommodel_ExtensionExternalData(arg, out)
	return out
}

func convert_etelecom_ExtensionExternalData_etelecommodel_ExtensionExternalData(arg *etelecom.ExtensionExternalData, out *etelecommodel.ExtensionExternalData) {
	out.ID = arg.ID // simple assign
}

func Convert_etelecom_ExtensionExternalDatas_etelecommodel_ExtensionExternalDatas(args []*etelecom.ExtensionExternalData) (outs []*etelecommodel.ExtensionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecommodel.ExtensionExternalData, len(args))
	outs = make([]*etelecommodel.ExtensionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_ExtensionExternalData_etelecommodel_ExtensionExternalData(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/etelecom.Hotline --//

func Convert_etelecommodel_Hotline_etelecom_Hotline(arg *etelecommodel.Hotline, out *etelecom.Hotline) *etelecom.Hotline {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Hotline{}
	}
	convert_etelecommodel_Hotline_etelecom_Hotline(arg, out)
	return out
}

func convert_etelecommodel_Hotline_etelecom_Hotline(arg *etelecommodel.Hotline, out *etelecom.Hotline) {
	out.ID = arg.ID                             // simple assign
	out.OwnerID = arg.OwnerID                   // simple assign
	out.Name = arg.Name                         // simple assign
	out.Hotline = arg.Hotline                   // simple assign
	out.Network = 0                             // types do not match
	out.ConnectionID = arg.ConnectionID         // simple assign
	out.ConnectionMethod = arg.ConnectionMethod // simple assign
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.DeletedAt = arg.DeletedAt               // simple assign
	out.Status = arg.Status                     // simple assign
	out.Description = arg.Description           // simple assign
	out.IsFreeCharge = arg.IsFreeCharge         // simple assign
}

func Convert_etelecommodel_Hotlines_etelecom_Hotlines(args []*etelecommodel.Hotline) (outs []*etelecom.Hotline) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.Hotline, len(args))
	outs = make([]*etelecom.Hotline, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecommodel_Hotline_etelecom_Hotline(args[i], &tmps[i])
	}
	return outs
}

func Convert_etelecom_Hotline_etelecommodel_Hotline(arg *etelecom.Hotline, out *etelecommodel.Hotline) *etelecommodel.Hotline {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecommodel.Hotline{}
	}
	convert_etelecom_Hotline_etelecommodel_Hotline(arg, out)
	return out
}

func convert_etelecom_Hotline_etelecommodel_Hotline(arg *etelecom.Hotline, out *etelecommodel.Hotline) {
	out.ID = arg.ID                             // simple assign
	out.OwnerID = arg.OwnerID                   // simple assign
	out.Name = arg.Name                         // simple assign
	out.Hotline = arg.Hotline                   // simple assign
	out.Network = ""                            // types do not match
	out.ConnectionID = arg.ConnectionID         // simple assign
	out.ConnectionMethod = arg.ConnectionMethod // simple assign
	out.CreatedAt = arg.CreatedAt               // simple assign
	out.UpdatedAt = arg.UpdatedAt               // simple assign
	out.DeletedAt = arg.DeletedAt               // simple assign
	out.Status = arg.Status                     // simple assign
	out.Description = arg.Description           // simple assign
	out.IsFreeCharge = arg.IsFreeCharge         // simple assign
}

func Convert_etelecom_Hotlines_etelecommodel_Hotlines(args []*etelecom.Hotline) (outs []*etelecommodel.Hotline) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecommodel.Hotline, len(args))
	outs = make([]*etelecommodel.Hotline, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_Hotline_etelecommodel_Hotline(args[i], &tmps[i])
	}
	return outs
}

func Apply_etelecom_CreateHotlineArgs_etelecom_Hotline(arg *etelecom.CreateHotlineArgs, out *etelecom.Hotline) *etelecom.Hotline {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Hotline{}
	}
	apply_etelecom_CreateHotlineArgs_etelecom_Hotline(arg, out)
	return out
}

func apply_etelecom_CreateHotlineArgs_etelecom_Hotline(arg *etelecom.CreateHotlineArgs, out *etelecom.Hotline) {
	out.ID = 0                          // zero value
	out.OwnerID = arg.OwnerID           // simple assign
	out.Name = arg.Name                 // simple assign
	out.Hotline = arg.Hotline           // simple assign
	out.Network = arg.Network           // simple assign
	out.ConnectionID = arg.ConnectionID // simple assign
	out.ConnectionMethod = 0            // zero value
	out.CreatedAt = time.Time{}         // zero value
	out.UpdatedAt = time.Time{}         // zero value
	out.DeletedAt = time.Time{}         // zero value
	out.Status = arg.Status             // simple assign
	out.Description = arg.Description   // simple assign
	out.IsFreeCharge = arg.IsFreeCharge // simple assign
}
