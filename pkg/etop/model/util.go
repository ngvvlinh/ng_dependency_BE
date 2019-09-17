package model

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/xerrors"
)

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`

	Meta map[string]string `json:"meta,omitempty"`
}

func (err *Error) Error() string {
	return err.Msg
}

func ToError(err error) *Error {
	if err == nil {
		return &Error{
			Code: "ok",
		}
	}
	xerr := xerrors.TwirpError(err)
	return &Error{
		Msg:  xerr.Msg(),
		Code: string(xerr.Code()),
		Meta: xerr.MetaMap(),
	}
}

func NewID() int64 {
	return cm.NewID()
}

func NewShopID() int64 {
	return cm.NewIDWithTag(TagShop)
}

func NewAffiliateID() int64 {
	return cm.NewIDWithTag(TagAffiliate)
}

func IsPartnerID(id int64) bool {
	return cm.GetTag(id) == TagPartner
}

func IsShopID(id int64) bool {
	return cm.GetTag(id) == TagShop
}

func IsAffiliateID(id int64) bool {
	return cm.GetTag(id) == TagAffiliate
}

func IsEtopAccountID(id int64) bool {
	return id == EtopAccountID
}

type UpdateListRequest struct {
	Adds       []string
	Deletes    []string
	ReplaceAll []string
	DeleteAll  bool
}

func (req *UpdateListRequest) Verify() error {
	c := 0
	if req.DeleteAll {
		c++
	}
	if len(req.ReplaceAll) > 0 {
		c++
	}
	if len(req.Adds) > 0 || len(req.Deletes) > 0 {
		c++
	}
	if c != 1 {
		return cm.Error(cm.InvalidArgument, "Must provide one of delete_all, replace_all or adds/deletes", nil)
	}
	return nil
}

func (req *UpdateListRequest) Patch(list []string) []string {
	if req.DeleteAll {
		return []string{}
	}
	if len(req.ReplaceAll) > 0 {
		return req.ReplaceAll
	}

	newList := make([]string, 0, len(list)+len(req.Adds))
	for _, item := range list {
		if !ListContain(newList, item) && !ListContain(req.Deletes, item) {
			newList = append(newList, item)
		}
	}
	for _, item := range req.Adds {
		if !ListContain(newList, item) && !ListContain(req.Deletes, item) {
			newList = append(newList, item)
		}
	}
	return newList
}

func ListContain(A []string, s string) bool {
	for _, a := range A {
		if a == s {
			return true
		}
	}
	return false
}

func PatchTag(sourceImages []string, cmd UpdateListRequest) ([]string, error) {
	for i, tag := range cmd.Adds {
		tag, ok := validate.NormalizeTag(tag)
		if !ok {
			return nil, cm.Error(cm.InvalidArgument, "Invalid tag: "+tag, nil)
		}
		cmd.Adds[i] = tag
	}
	for i, tag := range cmd.ReplaceAll {
		tag, ok := validate.NormalizeTag(tag)
		if !ok {
			return nil, cm.Error(cm.InvalidArgument, "Invalid tag: "+tag, nil)
		}
		cmd.ReplaceAll[i] = tag
	}

	return cmd.Patch(sourceImages), nil
}

func CoalesceString2(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}

func CalcVolumetricWeight(length, width, height int) int {
	return length * width * height / 5
}

func CalcChargeableWeight(weight, length, width, height int) int {
	volumetricWeight := CalcVolumetricWeight(length, width, height)
	if volumetricWeight > weight {
		return volumetricWeight
	}
	return weight
}
