package affiliate

import (
	"etop.vn/api/main/identity"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	"etop.vn/backend/pkg/etop/model"
)

func Convert_core_Affiliate_To_api_Affiliate(in *identity.Affiliate) *Affiliate {
	if in == nil {
		return nil
	}
	return &Affiliate{
		Id:     in.ID,
		Name:   in.Name,
		Status: pbs3.Pb(model.Status3(in.Status)),
		IsTest: in.IsTest != 0,
		Phone:  in.Phone,
		Email:  in.Email,
	}
}
