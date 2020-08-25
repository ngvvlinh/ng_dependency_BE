package partnerimport

import (
	"context"
	"time"

	"o.o/api/main/catalog"
	api "o.o/api/top/external/whitelabel"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func (s *ImportService) ProductCollections(ctx context.Context, r *api.ImportProductCollectionsRequest) (*api.ImportProductCollectionsResponse, error) {
	if len(r.ProductCollections) > MaximumItems {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var productIDs, collectionIDs []dot.ID
	for _, productCollection := range r.ProductCollections {
		product, err := s.shopProductStoreFactory(ctx).ExternalID(productCollection.ExternalProductID).GetShopProductDB()
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_product_id is invalid")
		}
		collection, err := s.shopCollectionStoreFactory(ctx).ExternalID(productCollection.ExternalCollectionID).GetShopCollectionDB()
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_collection_id is invalid")
		}
		productIDs = append(productIDs, product.ProductID)
		collectionIDs = append(collectionIDs, collection.ID)
		if _, err = s.shopProductCollectionStoreFactory(ctx).AddProductToCollection(&catalog.ShopProductCollection{
			PartnerID:            s.SS.Claim().AuthPartnerID,
			ExternalCollectionID: productCollection.ExternalCollectionID,
			ExternalProductID:    productCollection.ExternalProductID,
			ProductID:            product.ProductID,
			CollectionID:         collection.ID,
			ShopID:               s.SS.Shop().ID,
			CreatedAt:            productCollection.CreatedAt.ToTime(),
			UpdatedAt:            productCollection.UpdatedAt.ToTime(),
		}); err != nil {
			return nil, err
		}
		if productCollection.CreatedAt.IsZero() {
			productCollection.CreatedAt = dot.Time(time.Now())
		}
		if productCollection.UpdatedAt.IsZero() {
			productCollection.UpdatedAt = dot.Time(time.Now())
		}
	}

	var productCollectionsResponse []*api.ProductCollection
	for i, productCollection := range r.ProductCollections {
		productCollectionsResponse = append(productCollectionsResponse, &api.ProductCollection{
			PartnerID:            s.SS.Claim().AuthPartnerID,
			ShopID:               s.SS.Shop().ID,
			ExternalProductID:    productCollection.ExternalProductID,
			ExternalCollectionID: productCollection.ExternalCollectionID,
			ProductID:            productIDs[i],
			CollectionID:         collectionIDs[i],
			CreatedAt:            productCollection.CreatedAt,
			UpdatedAt:            productCollection.UpdatedAt,
		})
	}
	result := &api.ImportProductCollectionsResponse{ProductCollections: productCollectionsResponse}
	return result, nil
}
