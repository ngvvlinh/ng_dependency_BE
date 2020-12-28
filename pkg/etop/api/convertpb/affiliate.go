package convertpb

import (
	"o.o/api/main/identity"
	"o.o/api/top/int/etop"
)

func Convert_core_Affiliate_To_api_Affiliate(in *identity.Affiliate) *etop.Affiliate {
	if in == nil {
		return nil
	}
	return &etop.Affiliate{
		Id:          in.ID,
		Name:        in.Name,
		Status:      in.Status,
		IsTest:      in.IsTest != 0,
		Phone:       in.Phone,
		Email:       in.Email,
		BankAccount: Convert_core_BankAccount_To_api_BankAccount(in.BankAccount),
	}
}
