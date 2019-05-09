package main

import (
	"net/http"
	"testing"

	"etop.vn/backend/pkg/common/bus"
	etopW "etop.vn/backend/wrapper/etop"
	adminW "etop.vn/backend/wrapper/etop/admin"
	sadminW "etop.vn/backend/wrapper/etop/sadmin"
	shopW "etop.vn/backend/wrapper/etop/shop"
	supplierW "etop.vn/backend/wrapper/etop/supplier"
)

func TestWrapper(t *testing.T) {
	mux := http.NewServeMux()
	etopW.NewEtopServer(mux, nil)
	sadminW.NewSadminServer(mux, nil)
	adminW.NewAdminServer(mux, nil)
	supplierW.NewSupplierServer(mux, nil)
	shopW.NewShopServer(mux, nil)

	if !bus.Validate() {
		t.Error("Bus validation failed")
	}
}
