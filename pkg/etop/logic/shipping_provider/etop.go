package shipping_provider

/*
func (ctrl *ProviderManager) GetShippingService(ffm *model.Fulfillment, order *model.Order, weight int, valueInsurance int) (providerService *model.AvailableShippingService, etopService *model.AvailableShippingService, err error) {
	shopShipping := order.ShopShipping
	ctx := context.Background()

	_, fromDistrict, _, err := VerifyAddress(ffm.AddressFrom, false)
	if err != nil {
		return nil, nil, err
	}
	_, toDistrict, _, err := VerifyAddress(ffm.AddressTo, false)
	if err != nil {
		return nil, nil, cm.Errorf(cm.Internal, err, "ToAddress: %v", err)
	}

	args := GetShippingServicesArgs{
		ArbitraryID:      order.ShopID,
		FromDistrict:     fromDistrict,
		ToDistrict:       toDistrict,
		ChargeableWeight: weight,
		IncludeInsurance: true, // TODO: fix it
		BasketValue:      valueInsurance,
	}
	allServices, err := ctrl.GetShippingServices(ctx, args)
	if err != nil {
		return nil, nil, err
	}
	isEtopService, sType := CheckEtopServiceID(shopShipping.ProviderServiceID)
	if isEtopService {
		// ETOP serivce
		// => Get cheapest provider service
		etopService, err = GetEtopServiceFromShopShipping(order.ShopShipping, allServices)
		if err != nil {
			return nil, nil, err
		}
		providerService = GetCheapestService(allServices, sType)
		if providerService == nil {
			return nil, nil, cm.Error(cm.InvalidArgument, "Không có gói vận chuyển phù hợp.", nil)
		}
		return providerService, etopService, nil
	}
	// Provider service
	// => Check price
	// => Get this service
	providerService, err = checkShippingService(order, allServices)
	return providerService, nil, err
}
*/
