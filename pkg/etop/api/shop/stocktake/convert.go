package stocktake

import (
	"o.o/api/main/stocktaking"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbStocktakes(args []*stocktaking.ShopStocktake) []*shop.Stocktake {
	var stocktakesPb []*shop.Stocktake
	for _, value := range args {
		stocktakesPb = append(stocktakesPb, PbStocktake(value))
	}
	return stocktakesPb
}

func PbStocktake(args *stocktaking.ShopStocktake) *shop.Stocktake {
	if args == nil {
		return nil
	}
	return &shop.Stocktake{
		Id:            args.ID,
		ShopId:        args.ShopID,
		TotalQuantity: args.TotalQuantity,
		Note:          args.Note,
		CreatedBy:     args.CreatedBy,
		UpdatedBy:     args.UpdatedBy,
		CancelReason:  args.CancelReason,
		CreatedAt:     cmapi.PbTime(args.CreatedAt),
		UpdatedAt:     cmapi.PbTime(args.UpdatedAt),
		ConfirmedAt:   cmapi.PbTime(args.ConfirmedAt),
		CancelledAt:   cmapi.PbTime(args.CancelledAt),
		Status:        args.Status,
		Type:          args.Type.String(),
		Code:          args.Code,
		Lines:         PbstocktakeLines(args.Lines),
	}
}

func PbstocktakeLines(args []*stocktaking.StocktakeLine) []*shop.StocktakeLine {
	var lines []*shop.StocktakeLine
	for _, value := range args {
		lines = append(lines, &shop.StocktakeLine{
			VariantId:   value.VariantID,
			OldQuantity: value.OldQuantity,
			NewQuantity: value.NewQuantity,
			VariantName: value.VariantName,
			ProductName: value.ProductName,
			CostPrice:   value.CostPrice,
			ProductId:   value.ProductID,
			Code:        value.Code,
			ImageUrl:    value.ImageURL,
			Attributes:  value.Attributes,
		})
	}
	return lines
}
