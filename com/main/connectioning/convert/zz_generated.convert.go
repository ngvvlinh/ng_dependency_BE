// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	connectioning "etop.vn/api/main/connectioning"
	connectioningmodel "etop.vn/backend/com/main/connectioning/model"
	conversion "etop.vn/backend/pkg/common/conversion"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*connectioningmodel.Connection)(nil), (*connectioning.Connection)(nil), func(arg, out interface{}) error {
		Convert_connectioningmodel_Connection_connectioning_Connection(arg.(*connectioningmodel.Connection), out.(*connectioning.Connection))
		return nil
	})
	s.Register(([]*connectioningmodel.Connection)(nil), (*[]*connectioning.Connection)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioningmodel_Connections_connectioning_Connections(arg.([]*connectioningmodel.Connection))
		*out.(*[]*connectioning.Connection) = out0
		return nil
	})
	s.Register((*connectioning.Connection)(nil), (*connectioningmodel.Connection)(nil), func(arg, out interface{}) error {
		Convert_connectioning_Connection_connectioningmodel_Connection(arg.(*connectioning.Connection), out.(*connectioningmodel.Connection))
		return nil
	})
	s.Register(([]*connectioning.Connection)(nil), (*[]*connectioningmodel.Connection)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioning_Connections_connectioningmodel_Connections(arg.([]*connectioning.Connection))
		*out.(*[]*connectioningmodel.Connection) = out0
		return nil
	})
	s.Register((*connectioning.CreateConnectionArgs)(nil), (*connectioning.Connection)(nil), func(arg, out interface{}) error {
		Apply_connectioning_CreateConnectionArgs_connectioning_Connection(arg.(*connectioning.CreateConnectionArgs), out.(*connectioning.Connection))
		return nil
	})
	s.Register((*connectioning.UpdateConnectionAffiliateAccountArgs)(nil), (*connectioning.Connection)(nil), func(arg, out interface{}) error {
		Apply_connectioning_UpdateConnectionAffiliateAccountArgs_connectioning_Connection(arg.(*connectioning.UpdateConnectionAffiliateAccountArgs), out.(*connectioning.Connection))
		return nil
	})
	s.Register((*connectioning.UpdateConnectionArgs)(nil), (*connectioning.Connection)(nil), func(arg, out interface{}) error {
		Apply_connectioning_UpdateConnectionArgs_connectioning_Connection(arg.(*connectioning.UpdateConnectionArgs), out.(*connectioning.Connection))
		return nil
	})
	s.Register((*connectioningmodel.ConnectionService)(nil), (*connectioning.ConnectionService)(nil), func(arg, out interface{}) error {
		Convert_connectioningmodel_ConnectionService_connectioning_ConnectionService(arg.(*connectioningmodel.ConnectionService), out.(*connectioning.ConnectionService))
		return nil
	})
	s.Register(([]*connectioningmodel.ConnectionService)(nil), (*[]*connectioning.ConnectionService)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioningmodel_ConnectionServices_connectioning_ConnectionServices(arg.([]*connectioningmodel.ConnectionService))
		*out.(*[]*connectioning.ConnectionService) = out0
		return nil
	})
	s.Register((*connectioning.ConnectionService)(nil), (*connectioningmodel.ConnectionService)(nil), func(arg, out interface{}) error {
		Convert_connectioning_ConnectionService_connectioningmodel_ConnectionService(arg.(*connectioning.ConnectionService), out.(*connectioningmodel.ConnectionService))
		return nil
	})
	s.Register(([]*connectioning.ConnectionService)(nil), (*[]*connectioningmodel.ConnectionService)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioning_ConnectionServices_connectioningmodel_ConnectionServices(arg.([]*connectioning.ConnectionService))
		*out.(*[]*connectioningmodel.ConnectionService) = out0
		return nil
	})
	s.Register((*connectioningmodel.EtopAffiliateAccount)(nil), (*connectioning.EtopAffiliateAccount)(nil), func(arg, out interface{}) error {
		Convert_connectioningmodel_EtopAffiliateAccount_connectioning_EtopAffiliateAccount(arg.(*connectioningmodel.EtopAffiliateAccount), out.(*connectioning.EtopAffiliateAccount))
		return nil
	})
	s.Register(([]*connectioningmodel.EtopAffiliateAccount)(nil), (*[]*connectioning.EtopAffiliateAccount)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioningmodel_EtopAffiliateAccounts_connectioning_EtopAffiliateAccounts(arg.([]*connectioningmodel.EtopAffiliateAccount))
		*out.(*[]*connectioning.EtopAffiliateAccount) = out0
		return nil
	})
	s.Register((*connectioning.EtopAffiliateAccount)(nil), (*connectioningmodel.EtopAffiliateAccount)(nil), func(arg, out interface{}) error {
		Convert_connectioning_EtopAffiliateAccount_connectioningmodel_EtopAffiliateAccount(arg.(*connectioning.EtopAffiliateAccount), out.(*connectioningmodel.EtopAffiliateAccount))
		return nil
	})
	s.Register(([]*connectioning.EtopAffiliateAccount)(nil), (*[]*connectioningmodel.EtopAffiliateAccount)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioning_EtopAffiliateAccounts_connectioningmodel_EtopAffiliateAccounts(arg.([]*connectioning.EtopAffiliateAccount))
		*out.(*[]*connectioningmodel.EtopAffiliateAccount) = out0
		return nil
	})
	s.Register((*connectioningmodel.ShopConnection)(nil), (*connectioning.ShopConnection)(nil), func(arg, out interface{}) error {
		Convert_connectioningmodel_ShopConnection_connectioning_ShopConnection(arg.(*connectioningmodel.ShopConnection), out.(*connectioning.ShopConnection))
		return nil
	})
	s.Register(([]*connectioningmodel.ShopConnection)(nil), (*[]*connectioning.ShopConnection)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioningmodel_ShopConnections_connectioning_ShopConnections(arg.([]*connectioningmodel.ShopConnection))
		*out.(*[]*connectioning.ShopConnection) = out0
		return nil
	})
	s.Register((*connectioning.ShopConnection)(nil), (*connectioningmodel.ShopConnection)(nil), func(arg, out interface{}) error {
		Convert_connectioning_ShopConnection_connectioningmodel_ShopConnection(arg.(*connectioning.ShopConnection), out.(*connectioningmodel.ShopConnection))
		return nil
	})
	s.Register(([]*connectioning.ShopConnection)(nil), (*[]*connectioningmodel.ShopConnection)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioning_ShopConnections_connectioningmodel_ShopConnections(arg.([]*connectioning.ShopConnection))
		*out.(*[]*connectioningmodel.ShopConnection) = out0
		return nil
	})
	s.Register((*connectioning.CreateShopConnectionArgs)(nil), (*connectioning.ShopConnection)(nil), func(arg, out interface{}) error {
		Apply_connectioning_CreateShopConnectionArgs_connectioning_ShopConnection(arg.(*connectioning.CreateShopConnectionArgs), out.(*connectioning.ShopConnection))
		return nil
	})
	s.Register((*connectioning.UpdateShopConnectionExternalDataArgs)(nil), (*connectioning.ShopConnection)(nil), func(arg, out interface{}) error {
		Apply_connectioning_UpdateShopConnectionExternalDataArgs_connectioning_ShopConnection(arg.(*connectioning.UpdateShopConnectionExternalDataArgs), out.(*connectioning.ShopConnection))
		return nil
	})
	s.Register((*connectioningmodel.ShopConnectionExternalData)(nil), (*connectioning.ShopConnectionExternalData)(nil), func(arg, out interface{}) error {
		Convert_connectioningmodel_ShopConnectionExternalData_connectioning_ShopConnectionExternalData(arg.(*connectioningmodel.ShopConnectionExternalData), out.(*connectioning.ShopConnectionExternalData))
		return nil
	})
	s.Register(([]*connectioningmodel.ShopConnectionExternalData)(nil), (*[]*connectioning.ShopConnectionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioningmodel_ShopConnectionExternalDatas_connectioning_ShopConnectionExternalDatas(arg.([]*connectioningmodel.ShopConnectionExternalData))
		*out.(*[]*connectioning.ShopConnectionExternalData) = out0
		return nil
	})
	s.Register((*connectioning.ShopConnectionExternalData)(nil), (*connectioningmodel.ShopConnectionExternalData)(nil), func(arg, out interface{}) error {
		Convert_connectioning_ShopConnectionExternalData_connectioningmodel_ShopConnectionExternalData(arg.(*connectioning.ShopConnectionExternalData), out.(*connectioningmodel.ShopConnectionExternalData))
		return nil
	})
	s.Register(([]*connectioning.ShopConnectionExternalData)(nil), (*[]*connectioningmodel.ShopConnectionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_connectioning_ShopConnectionExternalDatas_connectioningmodel_ShopConnectionExternalDatas(arg.([]*connectioning.ShopConnectionExternalData))
		*out.(*[]*connectioningmodel.ShopConnectionExternalData) = out0
		return nil
	})
}

//-- convert etop.vn/api/main/connectioning.Connection --//

func Convert_connectioningmodel_Connection_connectioning_Connection(arg *connectioningmodel.Connection, out *connectioning.Connection) *connectioning.Connection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.Connection{}
	}
	convert_connectioningmodel_Connection_connectioning_Connection(arg, out)
	return out
}

func convert_connectioningmodel_Connection_connectioning_Connection(arg *connectioningmodel.Connection, out *connectioning.Connection) {
	out.ID = arg.ID                                 // simple assign
	out.Name = arg.Name                             // simple assign
	out.Status = arg.Status                         // simple assign
	out.PartnerID = arg.PartnerID                   // simple assign
	out.CreatedAt = arg.CreatedAt                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                   // simple assign
	out.DeletedAt = arg.DeletedAt                   // simple assign
	out.DriverConfig = arg.DriverConfig             // simple assign
	out.Driver = arg.Driver                         // simple assign
	out.ConnectionType = arg.ConnectionType         // simple assign
	out.ConnectionSubtype = arg.ConnectionSubtype   // simple assign
	out.ConnectionMethod = arg.ConnectionMethod     // simple assign
	out.ConnectionProvider = arg.ConnectionProvider // simple assign
	out.EtopAffiliateAccount = Convert_connectioningmodel_EtopAffiliateAccount_connectioning_EtopAffiliateAccount(arg.EtopAffiliateAccount, nil)
	out.Code = arg.Code         // simple assign
	out.ImageURL = arg.ImageURL // simple assign
	out.Services = Convert_connectioningmodel_ConnectionServices_connectioning_ConnectionServices(arg.Services)
	out.WLPartnerID = arg.WLPartnerID // simple assign
}

func Convert_connectioningmodel_Connections_connectioning_Connections(args []*connectioningmodel.Connection) (outs []*connectioning.Connection) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioning.Connection, len(args))
	outs = make([]*connectioning.Connection, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioningmodel_Connection_connectioning_Connection(args[i], &tmps[i])
	}
	return outs
}

func Convert_connectioning_Connection_connectioningmodel_Connection(arg *connectioning.Connection, out *connectioningmodel.Connection) *connectioningmodel.Connection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioningmodel.Connection{}
	}
	convert_connectioning_Connection_connectioningmodel_Connection(arg, out)
	return out
}

func convert_connectioning_Connection_connectioningmodel_Connection(arg *connectioning.Connection, out *connectioningmodel.Connection) {
	out.ID = arg.ID                                 // simple assign
	out.Name = arg.Name                             // simple assign
	out.Status = arg.Status                         // simple assign
	out.PartnerID = arg.PartnerID                   // simple assign
	out.CreatedAt = arg.CreatedAt                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                   // simple assign
	out.DeletedAt = arg.DeletedAt                   // simple assign
	out.DriverConfig = arg.DriverConfig             // simple assign
	out.Driver = arg.Driver                         // simple assign
	out.ConnectionType = arg.ConnectionType         // simple assign
	out.ConnectionSubtype = arg.ConnectionSubtype   // simple assign
	out.ConnectionMethod = arg.ConnectionMethod     // simple assign
	out.ConnectionProvider = arg.ConnectionProvider // simple assign
	out.EtopAffiliateAccount = Convert_connectioning_EtopAffiliateAccount_connectioningmodel_EtopAffiliateAccount(arg.EtopAffiliateAccount, nil)
	out.Code = arg.Code         // simple assign
	out.ImageURL = arg.ImageURL // simple assign
	out.Services = Convert_connectioning_ConnectionServices_connectioningmodel_ConnectionServices(arg.Services)
	out.WLPartnerID = arg.WLPartnerID // simple assign
}

func Convert_connectioning_Connections_connectioningmodel_Connections(args []*connectioning.Connection) (outs []*connectioningmodel.Connection) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioningmodel.Connection, len(args))
	outs = make([]*connectioningmodel.Connection, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioning_Connection_connectioningmodel_Connection(args[i], &tmps[i])
	}
	return outs
}

func Apply_connectioning_CreateConnectionArgs_connectioning_Connection(arg *connectioning.CreateConnectionArgs, out *connectioning.Connection) *connectioning.Connection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.Connection{}
	}
	apply_connectioning_CreateConnectionArgs_connectioning_Connection(arg, out)
	return out
}

func apply_connectioning_CreateConnectionArgs_connectioning_Connection(arg *connectioning.CreateConnectionArgs, out *connectioning.Connection) {
	out.ID = 0                                      // zero value
	out.Name = arg.Name                             // simple assign
	out.Status = 0                                  // zero value
	out.PartnerID = arg.PartnerID                   // simple assign
	out.CreatedAt = time.Time{}                     // zero value
	out.UpdatedAt = time.Time{}                     // zero value
	out.DeletedAt = time.Time{}                     // zero value
	out.DriverConfig = arg.DriverConfig             // simple assign
	out.Driver = arg.Driver                         // simple assign
	out.ConnectionType = arg.ConnectionType         // simple assign
	out.ConnectionSubtype = arg.ConnectionSubtype   // simple assign
	out.ConnectionMethod = arg.ConnectionMethod     // simple assign
	out.ConnectionProvider = arg.ConnectionProvider // simple assign
	out.EtopAffiliateAccount = nil                  // zero value
	out.Code = ""                                   // zero value
	out.ImageURL = ""                               // zero value
	out.Services = nil                              // zero value
	out.WLPartnerID = 0                             // zero value
}

func Apply_connectioning_UpdateConnectionAffiliateAccountArgs_connectioning_Connection(arg *connectioning.UpdateConnectionAffiliateAccountArgs, out *connectioning.Connection) *connectioning.Connection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.Connection{}
	}
	apply_connectioning_UpdateConnectionAffiliateAccountArgs_connectioning_Connection(arg, out)
	return out
}

func apply_connectioning_UpdateConnectionAffiliateAccountArgs_connectioning_Connection(arg *connectioning.UpdateConnectionAffiliateAccountArgs, out *connectioning.Connection) {
	out.ID = arg.ID                                     // simple assign
	out.Name = out.Name                                 // no change
	out.Status = out.Status                             // no change
	out.PartnerID = out.PartnerID                       // no change
	out.CreatedAt = out.CreatedAt                       // no change
	out.UpdatedAt = out.UpdatedAt                       // no change
	out.DeletedAt = out.DeletedAt                       // no change
	out.DriverConfig = out.DriverConfig                 // no change
	out.Driver = out.Driver                             // no change
	out.ConnectionType = out.ConnectionType             // no change
	out.ConnectionSubtype = out.ConnectionSubtype       // no change
	out.ConnectionMethod = out.ConnectionMethod         // no change
	out.ConnectionProvider = out.ConnectionProvider     // no change
	out.EtopAffiliateAccount = arg.EtopAffiliateAccount // simple assign
	out.Code = out.Code                                 // no change
	out.ImageURL = out.ImageURL                         // no change
	out.Services = out.Services                         // no change
	out.WLPartnerID = out.WLPartnerID                   // no change
}

func Apply_connectioning_UpdateConnectionArgs_connectioning_Connection(arg *connectioning.UpdateConnectionArgs, out *connectioning.Connection) *connectioning.Connection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.Connection{}
	}
	apply_connectioning_UpdateConnectionArgs_connectioning_Connection(arg, out)
	return out
}

func apply_connectioning_UpdateConnectionArgs_connectioning_Connection(arg *connectioning.UpdateConnectionArgs, out *connectioning.Connection) {
	out.ID = arg.ID                                     // simple assign
	out.Name = arg.Name                                 // simple assign
	out.Status = out.Status                             // no change
	out.PartnerID = out.PartnerID                       // identifier
	out.CreatedAt = out.CreatedAt                       // no change
	out.UpdatedAt = out.UpdatedAt                       // no change
	out.DeletedAt = out.DeletedAt                       // no change
	out.DriverConfig = arg.DriverConfig                 // simple assign
	out.Driver = out.Driver                             // no change
	out.ConnectionType = out.ConnectionType             // no change
	out.ConnectionSubtype = out.ConnectionSubtype       // no change
	out.ConnectionMethod = out.ConnectionMethod         // no change
	out.ConnectionProvider = out.ConnectionProvider     // no change
	out.EtopAffiliateAccount = out.EtopAffiliateAccount // no change
	out.Code = out.Code                                 // no change
	out.ImageURL = arg.ImageURL                         // simple assign
	out.Services = out.Services                         // no change
	out.WLPartnerID = out.WLPartnerID                   // no change
}

//-- convert etop.vn/api/main/connectioning.ConnectionService --//

func Convert_connectioningmodel_ConnectionService_connectioning_ConnectionService(arg *connectioningmodel.ConnectionService, out *connectioning.ConnectionService) *connectioning.ConnectionService {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.ConnectionService{}
	}
	convert_connectioningmodel_ConnectionService_connectioning_ConnectionService(arg, out)
	return out
}

func convert_connectioningmodel_ConnectionService_connectioning_ConnectionService(arg *connectioningmodel.ConnectionService, out *connectioning.ConnectionService) {
	out.ServiceID = arg.ServiceID // simple assign
	out.Name = arg.Name           // simple assign
}

func Convert_connectioningmodel_ConnectionServices_connectioning_ConnectionServices(args []*connectioningmodel.ConnectionService) (outs []*connectioning.ConnectionService) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioning.ConnectionService, len(args))
	outs = make([]*connectioning.ConnectionService, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioningmodel_ConnectionService_connectioning_ConnectionService(args[i], &tmps[i])
	}
	return outs
}

func Convert_connectioning_ConnectionService_connectioningmodel_ConnectionService(arg *connectioning.ConnectionService, out *connectioningmodel.ConnectionService) *connectioningmodel.ConnectionService {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioningmodel.ConnectionService{}
	}
	convert_connectioning_ConnectionService_connectioningmodel_ConnectionService(arg, out)
	return out
}

func convert_connectioning_ConnectionService_connectioningmodel_ConnectionService(arg *connectioning.ConnectionService, out *connectioningmodel.ConnectionService) {
	out.ServiceID = arg.ServiceID // simple assign
	out.Name = arg.Name           // simple assign
}

func Convert_connectioning_ConnectionServices_connectioningmodel_ConnectionServices(args []*connectioning.ConnectionService) (outs []*connectioningmodel.ConnectionService) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioningmodel.ConnectionService, len(args))
	outs = make([]*connectioningmodel.ConnectionService, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioning_ConnectionService_connectioningmodel_ConnectionService(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/connectioning.EtopAffiliateAccount --//

func Convert_connectioningmodel_EtopAffiliateAccount_connectioning_EtopAffiliateAccount(arg *connectioningmodel.EtopAffiliateAccount, out *connectioning.EtopAffiliateAccount) *connectioning.EtopAffiliateAccount {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.EtopAffiliateAccount{}
	}
	convert_connectioningmodel_EtopAffiliateAccount_connectioning_EtopAffiliateAccount(arg, out)
	return out
}

func convert_connectioningmodel_EtopAffiliateAccount_connectioning_EtopAffiliateAccount(arg *connectioningmodel.EtopAffiliateAccount, out *connectioning.EtopAffiliateAccount) {
	out.UserID = arg.UserID // simple assign
	out.Token = arg.Token   // simple assign
}

func Convert_connectioningmodel_EtopAffiliateAccounts_connectioning_EtopAffiliateAccounts(args []*connectioningmodel.EtopAffiliateAccount) (outs []*connectioning.EtopAffiliateAccount) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioning.EtopAffiliateAccount, len(args))
	outs = make([]*connectioning.EtopAffiliateAccount, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioningmodel_EtopAffiliateAccount_connectioning_EtopAffiliateAccount(args[i], &tmps[i])
	}
	return outs
}

func Convert_connectioning_EtopAffiliateAccount_connectioningmodel_EtopAffiliateAccount(arg *connectioning.EtopAffiliateAccount, out *connectioningmodel.EtopAffiliateAccount) *connectioningmodel.EtopAffiliateAccount {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioningmodel.EtopAffiliateAccount{}
	}
	convert_connectioning_EtopAffiliateAccount_connectioningmodel_EtopAffiliateAccount(arg, out)
	return out
}

func convert_connectioning_EtopAffiliateAccount_connectioningmodel_EtopAffiliateAccount(arg *connectioning.EtopAffiliateAccount, out *connectioningmodel.EtopAffiliateAccount) {
	out.UserID = arg.UserID // simple assign
	out.Token = arg.Token   // simple assign
}

func Convert_connectioning_EtopAffiliateAccounts_connectioningmodel_EtopAffiliateAccounts(args []*connectioning.EtopAffiliateAccount) (outs []*connectioningmodel.EtopAffiliateAccount) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioningmodel.EtopAffiliateAccount, len(args))
	outs = make([]*connectioningmodel.EtopAffiliateAccount, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioning_EtopAffiliateAccount_connectioningmodel_EtopAffiliateAccount(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/connectioning.ShopConnection --//

func Convert_connectioningmodel_ShopConnection_connectioning_ShopConnection(arg *connectioningmodel.ShopConnection, out *connectioning.ShopConnection) *connectioning.ShopConnection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.ShopConnection{}
	}
	convert_connectioningmodel_ShopConnection_connectioning_ShopConnection(arg, out)
	return out
}

func convert_connectioningmodel_ShopConnection_connectioning_ShopConnection(arg *connectioningmodel.ShopConnection, out *connectioning.ShopConnection) {
	out.ShopID = arg.ShopID                 // simple assign
	out.ConnectionID = arg.ConnectionID     // simple assign
	out.Token = arg.Token                   // simple assign
	out.TokenExpiresAt = arg.TokenExpiresAt // simple assign
	out.Status = arg.Status                 // simple assign
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
	out.DeletedAt = arg.DeletedAt           // simple assign
	out.IsGlobal = arg.IsGlobal             // simple assign
	out.ExternalData = Convert_connectioningmodel_ShopConnectionExternalData_connectioning_ShopConnectionExternalData(arg.ExternalData, nil)
}

func Convert_connectioningmodel_ShopConnections_connectioning_ShopConnections(args []*connectioningmodel.ShopConnection) (outs []*connectioning.ShopConnection) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioning.ShopConnection, len(args))
	outs = make([]*connectioning.ShopConnection, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioningmodel_ShopConnection_connectioning_ShopConnection(args[i], &tmps[i])
	}
	return outs
}

func Convert_connectioning_ShopConnection_connectioningmodel_ShopConnection(arg *connectioning.ShopConnection, out *connectioningmodel.ShopConnection) *connectioningmodel.ShopConnection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioningmodel.ShopConnection{}
	}
	convert_connectioning_ShopConnection_connectioningmodel_ShopConnection(arg, out)
	return out
}

func convert_connectioning_ShopConnection_connectioningmodel_ShopConnection(arg *connectioning.ShopConnection, out *connectioningmodel.ShopConnection) {
	out.ShopID = arg.ShopID                 // simple assign
	out.ConnectionID = arg.ConnectionID     // simple assign
	out.Token = arg.Token                   // simple assign
	out.TokenExpiresAt = arg.TokenExpiresAt // simple assign
	out.Status = arg.Status                 // simple assign
	out.ConnectionStates = nil              // zero value
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
	out.DeletedAt = arg.DeletedAt           // simple assign
	out.IsGlobal = arg.IsGlobal             // simple assign
	out.ExternalData = Convert_connectioning_ShopConnectionExternalData_connectioningmodel_ShopConnectionExternalData(arg.ExternalData, nil)
}

func Convert_connectioning_ShopConnections_connectioningmodel_ShopConnections(args []*connectioning.ShopConnection) (outs []*connectioningmodel.ShopConnection) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioningmodel.ShopConnection, len(args))
	outs = make([]*connectioningmodel.ShopConnection, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioning_ShopConnection_connectioningmodel_ShopConnection(args[i], &tmps[i])
	}
	return outs
}

func Apply_connectioning_CreateShopConnectionArgs_connectioning_ShopConnection(arg *connectioning.CreateShopConnectionArgs, out *connectioning.ShopConnection) *connectioning.ShopConnection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.ShopConnection{}
	}
	apply_connectioning_CreateShopConnectionArgs_connectioning_ShopConnection(arg, out)
	return out
}

func apply_connectioning_CreateShopConnectionArgs_connectioning_ShopConnection(arg *connectioning.CreateShopConnectionArgs, out *connectioning.ShopConnection) {
	out.ShopID = arg.ShopID                 // simple assign
	out.ConnectionID = arg.ConnectionID     // simple assign
	out.Token = arg.Token                   // simple assign
	out.TokenExpiresAt = arg.TokenExpiresAt // simple assign
	out.Status = 0                          // zero value
	out.CreatedAt = time.Time{}             // zero value
	out.UpdatedAt = time.Time{}             // zero value
	out.DeletedAt = time.Time{}             // zero value
	out.IsGlobal = false                    // zero value
	out.ExternalData = arg.ExternalData     // simple assign
}

func Apply_connectioning_UpdateShopConnectionExternalDataArgs_connectioning_ShopConnection(arg *connectioning.UpdateShopConnectionExternalDataArgs, out *connectioning.ShopConnection) *connectioning.ShopConnection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.ShopConnection{}
	}
	apply_connectioning_UpdateShopConnectionExternalDataArgs_connectioning_ShopConnection(arg, out)
	return out
}

func apply_connectioning_UpdateShopConnectionExternalDataArgs_connectioning_ShopConnection(arg *connectioning.UpdateShopConnectionExternalDataArgs, out *connectioning.ShopConnection) {
	out.ShopID = out.ShopID                 // identifier
	out.ConnectionID = out.ConnectionID     // identifier
	out.Token = arg.Token                   // simple assign
	out.TokenExpiresAt = arg.TokenExpiresAt // simple assign
	out.Status = out.Status                 // no change
	out.CreatedAt = out.CreatedAt           // no change
	out.UpdatedAt = out.UpdatedAt           // no change
	out.DeletedAt = out.DeletedAt           // no change
	out.IsGlobal = out.IsGlobal             // no change
	out.ExternalData = arg.ExternalData     // simple assign
}

//-- convert etop.vn/api/main/connectioning.ShopConnectionExternalData --//

func Convert_connectioningmodel_ShopConnectionExternalData_connectioning_ShopConnectionExternalData(arg *connectioningmodel.ShopConnectionExternalData, out *connectioning.ShopConnectionExternalData) *connectioning.ShopConnectionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioning.ShopConnectionExternalData{}
	}
	convert_connectioningmodel_ShopConnectionExternalData_connectioning_ShopConnectionExternalData(arg, out)
	return out
}

func convert_connectioningmodel_ShopConnectionExternalData_connectioning_ShopConnectionExternalData(arg *connectioningmodel.ShopConnectionExternalData, out *connectioning.ShopConnectionExternalData) {
	out.UserID = arg.UserID // simple assign
	out.Email = arg.Email   // simple assign
}

func Convert_connectioningmodel_ShopConnectionExternalDatas_connectioning_ShopConnectionExternalDatas(args []*connectioningmodel.ShopConnectionExternalData) (outs []*connectioning.ShopConnectionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioning.ShopConnectionExternalData, len(args))
	outs = make([]*connectioning.ShopConnectionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioningmodel_ShopConnectionExternalData_connectioning_ShopConnectionExternalData(args[i], &tmps[i])
	}
	return outs
}

func Convert_connectioning_ShopConnectionExternalData_connectioningmodel_ShopConnectionExternalData(arg *connectioning.ShopConnectionExternalData, out *connectioningmodel.ShopConnectionExternalData) *connectioningmodel.ShopConnectionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &connectioningmodel.ShopConnectionExternalData{}
	}
	convert_connectioning_ShopConnectionExternalData_connectioningmodel_ShopConnectionExternalData(arg, out)
	return out
}

func convert_connectioning_ShopConnectionExternalData_connectioningmodel_ShopConnectionExternalData(arg *connectioning.ShopConnectionExternalData, out *connectioningmodel.ShopConnectionExternalData) {
	out.UserID = arg.UserID // simple assign
	out.Email = arg.Email   // simple assign
}

func Convert_connectioning_ShopConnectionExternalDatas_connectioningmodel_ShopConnectionExternalDatas(args []*connectioning.ShopConnectionExternalData) (outs []*connectioningmodel.ShopConnectionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]connectioningmodel.ShopConnectionExternalData, len(args))
	outs = make([]*connectioningmodel.ShopConnectionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_connectioning_ShopConnectionExternalData_connectioningmodel_ShopConnectionExternalData(args[i], &tmps[i])
	}
	return outs
}
