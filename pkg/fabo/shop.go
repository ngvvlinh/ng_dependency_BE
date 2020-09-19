package fabo

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
	cm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
)

type ShopService struct {
	session.Session

	FBExternalUserQuery fbusering.QueryBus
	FBExternalUserAggr  fbusering.CommandBus
}

func (s *ShopService) Clone() fabo.ShopService {
	res := *s
	return &res
}

func (s *ShopService) CreateTag(ctx context.Context, req *fabo.CreateFbShopTagRequest) (*fabo.FbShopUserTag, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	cmd := &fbusering.CreateShopUserTagCommand{
		Name:   req.Name,
		Color:  req.Color,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.FBExternalUserAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	tag := cmd.Result
	resp := convertpb.Convert_core_FbShopUserTag_To_api_FbShopUserTag(tag)
	return resp, nil
}

func (s *ShopService) UpdateTag(ctx context.Context, req *fabo.UpdateFbShopTagRequest) (*fabo.FbShopUserTag, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	query := &fbusering.GetShopUserTagQuery{
		ID:     req.ID,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.FBExternalUserQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	cmdUpdate := &fbusering.UpdateShopUserTagCommand{
		Name:  req.Name,
		Color: req.Color,
		ID:    query.Result.ID,
	}
	if err := s.FBExternalUserAggr.Dispatch(ctx, cmdUpdate); err != nil {
		return nil, err
	}

	tag := cmdUpdate.Result
	resp := convertpb.Convert_core_FbShopUserTag_To_api_FbShopUserTag(tag)
	return resp, nil
}

func (s *ShopService) DeleteTag(ctx context.Context, req *fabo.DeleteFbShopTagRequest) (*cm.Empty, error) {
	cmdDelete := &fbusering.DeleteShopUserTagCommand{
		ID:     req.ID,
		ShopID: s.SS.Shop().ID,
	}
	err := s.FBExternalUserAggr.Dispatch(ctx, cmdDelete)
	return &cm.Empty{}, err
}

func (s *ShopService) GetTags(ctx context.Context, req *cm.Empty) (*fabo.ListFbShopTagResponse, error) {
	shopID := s.SS.Shop().ID
	queryListTag := &fbusering.ListShopUserTagsQuery{
		ShopID: shopID,
	}
	if err := s.FBExternalUserQuery.Dispatch(ctx, queryListTag); err != nil {
		return nil, err
	}

	tags := queryListTag.Result
	resp := &fabo.ListFbShopTagResponse{ShopTags: convertpb.Convert_core_FbShopUserTags_To_api_FbShopUserTags(tags)}
	return resp, nil
}
