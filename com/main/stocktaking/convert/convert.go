package convert

import (
	"fmt"
	"time"

	"o.o/api/main/stocktaking"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/stocktaking/model -> o.o/api/main/stocktaking
// +gen:convert:  o.o/api/main/stocktaking

const (
	MaxCodeNorm = 999999
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
