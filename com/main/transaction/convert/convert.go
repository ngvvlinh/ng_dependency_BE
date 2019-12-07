package convert

import (
	"etop.vn/api/main/transaction"
	transactionmodel "etop.vn/backend/com/main/transaction/model"
)

func TransactionMetadata(in *transactionmodel.TransactionMetadata) *transaction.TransactionMetadata {
	if in == nil {
		return nil
	}
	return &transaction.TransactionMetadata{
		ReferralType: transaction.ReferralType(in.ReferralType),
		ReferralIDs:  in.ReferralIDs,
	}
}

func TransactionMetadataDB(in *transaction.TransactionMetadata) *transactionmodel.TransactionMetadata {
	if in == nil {
		return nil
	}
	return &transactionmodel.TransactionMetadata{
		ReferralType: string(in.ReferralType),
		ReferralIDs:  in.ReferralIDs,
	}
}

func Transaction(in *transactionmodel.Transaction) *transaction.Transaction {
	if in == nil {
		return nil
	}
	return &transaction.Transaction{
		ID:        in.ID,
		Amount:    in.Amount,
		AccountID: in.AccountID,
		Status:    in.Status,
		Type:      transaction.TransactionType(in.Type),
		Note:      in.Note,
		Metadata:  TransactionMetadata(in.Metadata),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func Transactions(ins []*transactionmodel.Transaction) (outs []*transaction.Transaction) {
	for _, in := range ins {
		outs = append(outs, Transaction(in))
	}
	return
}

func TransactionDB(in *transaction.Transaction) *transactionmodel.Transaction {
	if in == nil {
		return nil
	}
	return &transactionmodel.Transaction{
		ID:        in.ID,
		Amount:    in.Amount,
		AccountID: in.AccountID,
		Status:    in.Status,
		Type:      string(in.Type),
		Note:      in.Note,
		Metadata:  TransactionMetadataDB(in.Metadata),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}
