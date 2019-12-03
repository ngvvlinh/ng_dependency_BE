package convert

import (
	"fmt"
	"time"

	"etop.vn/api/main/stocktaking"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/main/stocktaking/model->etop.vn/api/main/stocktaking
// +gen:convert: etop.vn/api/main/stocktaking

const (
	MaxCodeNorm = 999999
	codeRegex   = "^PKK([0-9]{6})$"
	codePrefix  = "PKK"
)

func ConvertCreateStocktake(arg *stocktaking.CreateStocktakeRequest, out *stocktaking.ShopStocktake) *stocktaking.ShopStocktake {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &stocktaking.ShopStocktake{}
	}
	apply_stocktaking_CreateStocktakeRequest_stocktaking_ShopStocktake(arg, out)
	out.ID = cm.NewID()
	out.UpdatedBy = out.CreatedBy
	out.CreatedAt = time.Now()
	out.UpdatedAt = time.Now()
	return out
}

func ConvertUpdateStocktake(arg *stocktaking.UpdateStocktakeRequest, out *stocktaking.ShopStocktake) *stocktaking.ShopStocktake {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &stocktaking.ShopStocktake{}
	}
	apply_stocktaking_UpdateStocktakeRequest_stocktaking_ShopStocktake(arg, out)
	out.UpdatedAt = time.Now()
	return out
}

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}
