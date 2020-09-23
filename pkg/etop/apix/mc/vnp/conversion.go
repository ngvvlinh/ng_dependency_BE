package vnp

import (
	"o.o/api/top/external/mc/vnp"
	vnpentitytype "o.o/api/top/external/mc/vnp/etc/entity_type"
	"o.o/api/top/external/types"
)

func Convert_apix_Webhook_To_vnp_Webhook(in *types.Webhook) *vnp.Webhook {
	if in == nil {
		return nil
	}
	vnpEntities, _ := vnpentitytype.Convert_type_Entities_To_type_VnpEntities(in.Entities)
	res := &vnp.Webhook{
		ID:        in.Id,
		Entities:  vnpEntities,
		URL:       in.Url,
		CreatedAt: in.CreatedAt,
		States:    in.States,
	}
	return res
}

func Convert_apix_Webhooks_To_vnp_Webhooks(ins []*types.Webhook) (res []*vnp.Webhook) {
	if ins == nil {
		return nil
	}
	for _, in := range ins {
		res = append(res, Convert_apix_Webhook_To_vnp_Webhook(in))
	}
	return
}
