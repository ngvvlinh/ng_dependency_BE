package convert

import (
	"etop.vn/api/main/location"
	"etop.vn/backend/pkg/services/location/types"
)

func Provinces(ins []*types.Province, outs *[]*location.Province) error {
	if *outs == nil {
		*outs = make([]*location.Province, len(ins))
	}
	for i, in := range ins {
		out := &location.Province{}
		if err := Province(in, out); err != nil {
			return err
		}
		(*outs)[i] = out
	}
	return nil
}

func Districts(ins []*types.District, outs *[]*location.District) error {
	if *outs == nil {
		*outs = make([]*location.District, len(ins))
	}
	for i, in := range ins {
		out := &location.District{}
		if err := District(in, out); err != nil {
			return err
		}
		(*outs)[i] = out
	}
	return nil
}

func Wards(ins []*types.Ward, outs *[]*location.Ward) error {
	if *outs == nil {
		*outs = make([]*location.Ward, len(ins))
	}
	for i, in := range ins {
		out := &location.Ward{}
		if err := Ward(in, out); err != nil {
			return err
		}
		(*outs)[i] = out
	}
	return nil
}

func PtrProvince(in *types.Province, out **location.Province) error {
	if in == nil {
		return nil
	}
	if *out == nil {
		*out = &location.Province{}
	}
	return Province(in, *out)
}

func Province(in *types.Province, out *location.Province) error {
	*out = location.Province{
		Name:   in.Name,
		Code:   in.Code,
		Region: in.Region,
		Extra: location.LocationExtra{
			Special:  in.Special,
			GhnId:    0,
			VtpostId: in.VTPostID,
		},
	}
	return nil
}

func PtrDistrict(in *types.District, out **location.District) error {
	if in == nil {
		return nil
	}
	if *out == nil {
		*out = &location.District{}
	}
	return District(in, *out)
}

func District(in *types.District, out *location.District) error {
	*out = location.District{
		Name:         in.Name,
		Code:         in.Code,
		ProvinceCode: in.ProvinceCode,
		UrbanType:    location.UrbanType(in.UrbanType),
		Extra: location.LocationExtra{
			GhnId:    in.GhnID,
			VtpostId: in.VTPostID,
		},
	}
	return nil
}

func PtrWard(in *types.Ward, out **location.Ward) error {
	if in == nil {
		return nil
	}
	if *out == nil {
		*out = &location.Ward{}
	}
	return Ward(in, *out)
}

func Ward(in *types.Ward, out *location.Ward) error {
	*out = location.Ward{
		Name:         in.Name,
		Code:         in.Code,
		DistrictCode: in.DistrictCode,
		Extra: location.LocationExtra{
			VtpostId: in.VTPostID,
		},
	}
	return nil
}
