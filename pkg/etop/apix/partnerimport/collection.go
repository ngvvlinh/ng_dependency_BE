package partnerimport

import (
	"context"

	"o.o/api/main/catalog"
	api "o.o/api/top/external/whitelabel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *ImportService) Collections(ctx context.Context, r *api.ImportCollectionsRequest) (*api.ImportCollectionsResponse, error) {
	if len(r.Collections) > MaximumItems {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, collection := range r.Collections {

		if collection.ExternalID == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "id should not be null")
		}
		shopCollection := &catalog.ShopCollection{
			ShopID:      s.SS.Shop().ID,
			PartnerID:   s.SS.Claim().AuthPartnerID,
			ExternalID:  collection.ExternalID,
			Name:        collection.Name,
			Description: collection.Description,
			DescHTML:    collection.DescHTML,
			ShortDesc:   collection.ShortDesc,
			CreatedAt:   collection.CreatedAt.ToTime(),
			UpdatedAt:   collection.UpdatedAt.ToTime(),
		}

		oldShopCollection, err := s.shopCollectionStoreFactory(ctx).ExternalID(collection.ExternalID).GetShopCollectionDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			shopCollection.ID = id
			ids = append(ids, id)
			if _err := s.shopCollectionStoreFactory(ctx).CreateShopCollection(shopCollection); _err != nil {
				return nil, _err
			}
		case cm.NoError:
			shopCollection.ID = oldShopCollection.ID
			ids = append(ids, oldShopCollection.ID)
			if _err := s.shopCollectionStoreFactory(ctx).UpdateShopCollection(shopCollection); _err != nil {
				return nil, _err
			}
		default:
			return nil, err
		}
	}

	modelCollections, err := s.shopCollectionStoreFactory(ctx).IDs(ids).ListShopCollectionsDB()
	if err != nil {
		return nil, err
	}

	var collectionsResponse []*api.Collection
	for _, collection := range modelCollections {
		collectionsResponse = append(collectionsResponse, &api.Collection{
			ID:          collection.ID,
			ShopID:      collection.ShopID,
			PartnerID:   collection.PartnerID,
			ExternalID:  collection.ExternalID,
			Name:        collection.Name,
			Description: collection.Description,
			DescHTML:    collection.DescHTML,
			ShortDesc:   collection.ShortDesc,
			CreatedAt:   cmapi.PbTime(collection.CreatedAt),
			UpdatedAt:   cmapi.PbTime(collection.UpdatedAt),
		})
	}
	result := &api.ImportCollectionsResponse{Collections: collectionsResponse}
	return result, nil
}
