// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	bankstatement "o.o/api/main/bankstatement"
	bankstatementmodel "o.o/backend/com/main/bankstatement/model"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions:
    BankStatementToModel    // in use

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*bankstatementmodel.BankStatement)(nil), (*bankstatement.BankStatement)(nil), func(arg, out interface{}) error {
		Convert_bankstatementmodel_BankStatement_bankstatement_BankStatement(arg.(*bankstatementmodel.BankStatement), out.(*bankstatement.BankStatement))
		return nil
	})
	s.Register(([]*bankstatementmodel.BankStatement)(nil), (*[]*bankstatement.BankStatement)(nil), func(arg, out interface{}) error {
		out0 := Convert_bankstatementmodel_BankStatements_bankstatement_BankStatements(arg.([]*bankstatementmodel.BankStatement))
		*out.(*[]*bankstatement.BankStatement) = out0
		return nil
	})
	s.Register((*bankstatement.BankStatement)(nil), (*bankstatementmodel.BankStatement)(nil), func(arg, out interface{}) error {
		Convert_bankstatement_BankStatement_bankstatementmodel_BankStatement(arg.(*bankstatement.BankStatement), out.(*bankstatementmodel.BankStatement))
		return nil
	})
	s.Register(([]*bankstatement.BankStatement)(nil), (*[]*bankstatementmodel.BankStatement)(nil), func(arg, out interface{}) error {
		out0 := Convert_bankstatement_BankStatements_bankstatementmodel_BankStatements(arg.([]*bankstatement.BankStatement))
		*out.(*[]*bankstatementmodel.BankStatement) = out0
		return nil
	})
	s.Register((*bankstatement.CreateBankStatementArgs)(nil), (*bankstatement.BankStatement)(nil), func(arg, out interface{}) error {
		Apply_bankstatement_CreateBankStatementArgs_bankstatement_BankStatement(arg.(*bankstatement.CreateBankStatementArgs), out.(*bankstatement.BankStatement))
		return nil
	})
}

//-- convert o.o/api/main/bankstatement.BankStatement --//

func Convert_bankstatementmodel_BankStatement_bankstatement_BankStatement(arg *bankstatementmodel.BankStatement, out *bankstatement.BankStatement) *bankstatement.BankStatement {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &bankstatement.BankStatement{}
	}
	convert_bankstatementmodel_BankStatement_bankstatement_BankStatement(arg, out)
	return out
}

func convert_bankstatementmodel_BankStatement_bankstatement_BankStatement(arg *bankstatementmodel.BankStatement, out *bankstatement.BankStatement) {
	out.ID = arg.ID                                       // simple assign
	out.Amount = arg.Amount                               // simple assign
	out.Description = arg.Description                     // simple assign
	out.AccountID = arg.AccountID                         // simple assign
	out.TransferedAt = arg.TransferedAt                   // simple assign
	out.ExternalTransactionID = arg.ExternalTransactionID // simple assign
	out.SenderName = arg.SenderName                       // simple assign
	out.SenderBankAccount = arg.SenderBankAccount         // simple assign
	out.OtherInfo = nil                                   // types do not match
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
}

func Convert_bankstatementmodel_BankStatements_bankstatement_BankStatements(args []*bankstatementmodel.BankStatement) (outs []*bankstatement.BankStatement) {
	if args == nil {
		return nil
	}
	tmps := make([]bankstatement.BankStatement, len(args))
	outs = make([]*bankstatement.BankStatement, len(args))
	for i := range tmps {
		outs[i] = Convert_bankstatementmodel_BankStatement_bankstatement_BankStatement(args[i], &tmps[i])
	}
	return outs
}

func Convert_bankstatement_BankStatement_bankstatementmodel_BankStatement(arg *bankstatement.BankStatement, out *bankstatementmodel.BankStatement) *bankstatementmodel.BankStatement {
	return BankStatementToModel(arg)
}

func convert_bankstatement_BankStatement_bankstatementmodel_BankStatement(arg *bankstatement.BankStatement, out *bankstatementmodel.BankStatement) {
	out.ID = arg.ID                                       // simple assign
	out.Amount = arg.Amount                               // simple assign
	out.Description = arg.Description                     // simple assign
	out.AccountID = arg.AccountID                         // simple assign
	out.TransferedAt = arg.TransferedAt                   // simple assign
	out.ExternalTransactionID = arg.ExternalTransactionID // simple assign
	out.SenderName = arg.SenderName                       // simple assign
	out.SenderBankAccount = arg.SenderBankAccount         // simple assign
	out.OtherInfo = nil                                   // types do not match
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
}

func Convert_bankstatement_BankStatements_bankstatementmodel_BankStatements(args []*bankstatement.BankStatement) (outs []*bankstatementmodel.BankStatement) {
	if args == nil {
		return nil
	}
	tmps := make([]bankstatementmodel.BankStatement, len(args))
	outs = make([]*bankstatementmodel.BankStatement, len(args))
	for i := range tmps {
		outs[i] = Convert_bankstatement_BankStatement_bankstatementmodel_BankStatement(args[i], &tmps[i])
	}
	return outs
}

func Apply_bankstatement_CreateBankStatementArgs_bankstatement_BankStatement(arg *bankstatement.CreateBankStatementArgs, out *bankstatement.BankStatement) *bankstatement.BankStatement {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &bankstatement.BankStatement{}
	}
	apply_bankstatement_CreateBankStatementArgs_bankstatement_BankStatement(arg, out)
	return out
}

func apply_bankstatement_CreateBankStatementArgs_bankstatement_BankStatement(arg *bankstatement.CreateBankStatementArgs, out *bankstatement.BankStatement) {
	out.ID = arg.ID                                       // simple assign
	out.Amount = arg.Amount                               // simple assign
	out.Description = arg.Description                     // simple assign
	out.AccountID = arg.AccountID                         // simple assign
	out.TransferedAt = arg.TransferedAt                   // simple assign
	out.ExternalTransactionID = arg.ExternalTransactionID // simple assign
	out.SenderName = arg.SenderName                       // simple assign
	out.SenderBankAccount = arg.SenderBankAccount         // simple assign
	out.OtherInfo = nil                                   // types do not match
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
}
