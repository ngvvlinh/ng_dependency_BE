package sqlstore

import (
	"context"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/etop/model"
)

func createCode(ctx context.Context, x Qx, cmd *model.CreateCodeCommand) (int, error) {
	code := cmd.Code
	if err := cmd.Code.Validate(); err != nil {
		return 0, err
	}
	return x.Table("code").
		Suffix("ON CONFLICT (code,type) DO NOTHING").
		Insert(code)
}

func GenerateCodeWithoutTransaction(ctx context.Context, typeCode model.CodeType, subCode string) (string, error) {
	return GenerateCode(ctx, x, typeCode, subCode)
}

func GenerateCode(ctx context.Context, x Qx, typeCode model.CodeType, subCode string) (string, error) {
	var fn func() string
	switch typeCode {
	case model.CodeTypeShop:
		fn = gencode.GenerateShopCode
	case model.CodeTypeOrder:
		fn = func() string {
			return gencode.GenerateOrderCode(subCode, time.Now())
		}
	case model.CodeTypeMoneyTransaction:
		fn = func() string {
			return gencode.GenerateCodeWithType("M", subCode, time.Now())
		}
	case model.CodeTypeMoneyTransactionEtop:
		fn = func() string {
			return gencode.GenerateCodeWithType("Mt", subCode, time.Now())
		}
	case model.CodeTypeMoneyTransactionExternal:
		fn = func() string {
			return gencode.GenerateCodeWithType("Mx", subCode, time.Now())
		}
	default:
		return "", cm.Errorf(cm.Internal, nil, "Invalid code type: %v", typeCode)
	}

	return generateCode(ctx, x, typeCode, fn)
}

func generateCode(ctx context.Context, x Qx, typeCode model.CodeType, fn func() string) (string, error) {
	const maxRetry = 5
	for retry := 0; retry < maxRetry; retry++ {
		code := fn()
		cmd := &model.CreateCodeCommand{
			Code: &model.Code{
				Code: code,
				Type: typeCode,
			},
		}
		n, err := createCode(ctx, x, cmd)
		if err != nil {
			return "", cm.Errorf(cm.Internal, err, "Can not generate code for type: %v", typeCode)
		}
		if n != 0 {
			return code, nil
		}
	}
	return "", cm.Errorf(cm.Internal, nil, "Can not generate code for type: %v", typeCode).WithMeta("reason", "retried too many times")
}
