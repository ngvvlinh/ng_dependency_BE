package convertpb

import (
	identitytypes "o.o/api/main/identity/types"
	etop "o.o/api/top/int/etop"
)

func Convert_api_BankAccount_To_core_BankAccount(in *etop.BankAccount) *identitytypes.BankAccount {
	if in == nil {
		return nil
	}
	return &identitytypes.BankAccount{
		Name:          in.Name,
		Province:      in.Province,
		Branch:        in.Branch,
		AccountNumber: in.AccountNumber,
		AccountName:   in.AccountName,
	}
}
