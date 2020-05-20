	alter table fulfillment add column delivery_route text;
	alter table history.fulfillment add column delivery_route text;

	alter table shipnow_fulfillment add column  address_to_province_code text;
	alter table shipnow_fulfillment add column  address_to_district_code text;
