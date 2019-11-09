package shop

import (
	"context"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/export"
)

func init() {
	bus.AddHandlers("export",
		exportService.RequestExport,
		exportService.GetExports,
	)
}

func (s *ExportService) RequestExport(ctx context.Context, r *RequestExportEndpoint) (_err error) {
	resp, err := export.ServiceImpl.RequestExport(ctx, r.Context, r.Context.Shop, r.Context.UserID, r.RequestExportRequest)
	r.Result = resp
	return err
}

func (s *ExportService) GetExports(ctx context.Context, r *GetExportsEndpoint) error {
	resp, err := export.ServiceImpl.GetExports(ctx, r.Context.Shop.ID, r.GetExportsRequest)
	r.Result = resp
	return err
}
