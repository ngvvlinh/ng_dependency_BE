package whitelabel

import (
	"context"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/external/whitelabel"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

func (s *ImportService) ProductCollections(ctx context.Context, r *ProductCollectionsEndpoint) error {
	if len(r.ProductCollections) > MaximumItems {
		return cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var productIDs, collectionIDs []dot.ID
	for _, productCollection := range r.ProductCollections {
		product, err := shopProductStoreFactory(ctx).ExternalID(productCollection.ExternalProductID).GetShopProductDB()
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "external_product_id is invalid")
		}
		collection, err := shopCollectionStoreFactory(ctx).ExternalID(productCollection.ExternalCollectionID).GetShopCollectionDB()
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "external_collection_id is invalid")
		}
		productIDs = append(productIDs, product.ProductID)
		collectionIDs = append(collectionIDs, collection.ID)
		if _, err := shopProductCollectionStoreFactory(ctx).AddProductToCollection(&catalog.ShopProductCollection{
			PartnerID:            r.Context.AuthPartnerID,
			ExternalCollectionID: productCollection.ExternalCollectionID,
			ExternalProductID:    productCollection.ExternalProductID,
			ProductID:            product.ProductID,
			CollectionID:         collection.ID,
			ShopID:               r.Context.Shop.ID,
			CreatedAt:            productCollection.CreatedAt.ToTime(),
			UpdatedAt:            productCollection.UpdatedAt.ToTime(),
		}); err != nil {
			return err
		}
		if productCollection.CreatedAt.IsZero() {
			productCollection.CreatedAt = dot.Time(time.Now())
		}
		if productCollection.UpdatedAt.IsZero() {
			productCollection.UpdatedAt = dot.Time(time.Now())
		}
	}

	var productCollectionsResponse []*whitelabel.ProductCollection
	for i, productCollection := range r.ProductCollections {
		productCollectionsResponse = append(productCollectionsResponse, &whitelabel.ProductCollection{
			PartnerID:            r.Context.AuthPartnerID,
			ShopID:               r.Context.Shop.ID,
			ExternalProductID:    productCollection.ExternalProductID,
			ExternalCollectionID: productCollection.ExternalCollectionID,
			ProductID:            productIDs[i],
			CollectionID:         collectionIDs[i],
			CreatedAt:            productCollection.CreatedAt,
			UpdatedAt:            productCollection.UpdatedAt,
		})
	}
	r.Result = &whitelabel.ImportProductCollectionsResponse{ProductCollections: productCollectionsResponse}
	return nil
}
