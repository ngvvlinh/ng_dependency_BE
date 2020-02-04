package shopping

import (
	"context"

	"etop.vn/api/main/catalog"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/api/top/types/common"
	cm "etop.vn/api/top/types/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
)

func GetCollection(ctx context.Context, shopID dot.ID, request *externaltypes.GetCollectionRequest) (*externaltypes.ProductCollection, error) {
	query := &catalog.GetShopCollectionQuery{
		ID:     request.ID,
		ShopID: shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopProductCollection(query.Result), nil
}

func ListCollections(ctx context.Context, shopID dot.ID, request *externaltypes.ListCollectionsRequest) (*externaltypes.ProductCollectionsResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	query := &catalog.ListShopCollectionsByIDsQuery{
		IDs:            request.Filter.ID,
		ShopID:         shopID,
		Paging:         *paging,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.ProductCollectionsResponse{
		Collections: convertpb.PbShopProductCollections(query.Result.Collections),
		Paging:      convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func CreateCollection(ctx context.Context, shopID, partnerID dot.ID, request *externaltypes.CreateCollectionRequest) (*externaltypes.ProductCollection, error) {
	cmd := &catalog.CreateShopCollectionCommand{
		ShopID:      shopID,
		PartnerID:   partnerID,
		Name:        request.Name,
		Description: request.Description,
		ShortDesc:   request.ShortDesc,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopProductCollection(cmd.Result), nil
}

func UpdateCollection(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateCollectionRequest) (*externaltypes.ProductCollection, error) {
	cmd := &catalog.UpdateShopCollectionCommand{
		ID:          request.ID,
		ShopID:      shopID,
		Name:        request.Name,
		Description: request.Description,
		ShortDesc:   request.ShortDesc,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopProductCollection(cmd.Result), nil
}

func DeleteCollection(ctx context.Context, shopID dot.ID, request *externaltypes.GetCollectionRequest) (*cm.Empty, error) {
	cmd := &catalog.DeleteShopCollectionCommand{
		Id:     request.ID,
		ShopId: shopID,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func ListRelationshipsProductCollection(ctx context.Context, shopID dot.ID, request *externaltypes.ListProductCollectionRelationshipsRequest) (*externaltypes.ProductCollectionRelationshipsResponse, error) {
	// TODO: add cursor paging
	query := &catalog.ListShopProductsCollectionsQuery{
		ProductIds:     request.Filter.ProductID,
		CollectionIDs:  request.Filter.CollectionID,
		ShopID:         shopID,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.ProductCollectionRelationshipsResponse{
		Relationships: convertpb.PbShopProductCollectionRelationships(query.Result.ProductsCollections),
	}, nil
}

func CreateRelationshipProductCollection(ctx context.Context, shopID dot.ID, request *externaltypes.CreateProductCollectionRelationshipRequest) (*cm.Empty, error) {
	cmd := &catalog.AddShopProductCollectionCommand{
		ProductID:     request.ProductId,
		ShopID:        shopID,
		CollectionIDs: []dot.ID{request.CollectionId},
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func DeleteRelationshipProductCollection(ctx context.Context, shopID dot.ID, request *externaltypes.RemoveProductCollectionRequest) (*cm.Empty, error) {

	cmd := &catalog.RemoveShopProductCollectionCommand{
		ProductID:     request.ProductId,
		ShopID:        shopID,
		CollectionIDs: []dot.ID{request.CollectionId},
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
