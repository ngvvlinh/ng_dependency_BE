package shopping

import (
	"context"

	"o.o/api/main/catalog"
	externaltypes "o.o/api/top/external/types"
	"o.o/api/top/types/common"
	cm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
)

func (s *Shopping) GetCollection(ctx context.Context, shopID dot.ID, request *externaltypes.GetCollectionRequest) (*externaltypes.ProductCollection, error) {
	query := &catalog.GetShopCollectionQuery{
		ID:     request.ID,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopProductCollection(query.Result), nil
}

func (s *Shopping) ListCollections(ctx context.Context, shopID dot.ID, request *externaltypes.ListCollectionsRequest) (*externaltypes.ProductCollectionsResponse, error) {
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
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.ProductCollectionsResponse{
		Collections: convertpb.PbShopProductCollections(query.Result.Collections),
		Paging:      convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func (s *Shopping) CreateCollection(ctx context.Context, shopID, partnerID dot.ID, request *externaltypes.CreateCollectionRequest) (*externaltypes.ProductCollection, error) {
	cmd := &catalog.CreateShopCollectionCommand{
		ShopID:      shopID,
		PartnerID:   partnerID,
		Name:        request.Name,
		Description: request.Description,
		ShortDesc:   request.ShortDesc,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopProductCollection(cmd.Result), nil
}

func (s *Shopping) UpdateCollection(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateCollectionRequest) (*externaltypes.ProductCollection, error) {
	cmd := &catalog.UpdateShopCollectionCommand{
		ID:          request.ID,
		ShopID:      shopID,
		Name:        request.Name,
		Description: request.Description,
		ShortDesc:   request.ShortDesc,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopProductCollection(cmd.Result), nil
}

func (s *Shopping) DeleteCollection(ctx context.Context, shopID dot.ID, request *externaltypes.GetCollectionRequest) (*cm.Empty, error) {
	cmd := &catalog.DeleteShopCollectionCommand{
		Id:     request.ID,
		ShopId: shopID,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (s *Shopping) ListRelationshipsProductCollection(ctx context.Context, shopID dot.ID, request *externaltypes.ListProductCollectionRelationshipsRequest) (*externaltypes.ProductCollectionRelationshipsResponse, error) {
	// TODO: add cursor paging
	query := &catalog.ListShopProductsCollectionsQuery{
		ProductIds:     request.Filter.ProductID,
		CollectionIDs:  request.Filter.CollectionID,
		ShopID:         shopID,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.ProductCollectionRelationshipsResponse{
		Relationships: convertpb.PbShopProductCollectionRelationships(query.Result.ProductsCollections),
	}, nil
}

func (s *Shopping) CreateRelationshipProductCollection(ctx context.Context, shopID dot.ID, request *externaltypes.CreateProductCollectionRelationshipRequest) (*cm.Empty, error) {
	cmd := &catalog.AddShopProductCollectionCommand{
		ProductID:     request.ProductId,
		ShopID:        shopID,
		CollectionIDs: []dot.ID{request.CollectionId},
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (s *Shopping) DeleteRelationshipProductCollection(ctx context.Context, shopID dot.ID, request *externaltypes.RemoveProductCollectionRequest) (*cm.Empty, error) {

	cmd := &catalog.RemoveShopProductCollectionCommand{
		ProductID:     request.ProductId,
		ShopID:        shopID,
		CollectionIDs: []dot.ID{request.CollectionId},
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
