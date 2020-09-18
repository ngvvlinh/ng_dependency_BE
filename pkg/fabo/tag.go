package fabo

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
	cm "o.o/api/top/types/common"
	"o.o/backend/pkg/fabo/convertpb"
)

var (
	defaultTagTemplate = []*fbusering.FbShopTag{
		{
			Name:  "Chốt Đơn",
			Color: "3498db",
		},
		{
			Name:  "Đã Ship",
			Color: "2ecc71",
		},
		{
			Name:  "Hỏi Giá",
			Color: "95a5a6",
		},
		{
			Name:  "Tư Vấn",
			Color: "e74c3c",
		},
		{
			Name:  "Bank",
			Color: "9b59b6",
		},
		{
			Name:  "COD",
			Color: "f39c12",
		},
	}
)

func (s *PageService) CreateTag(ctx context.Context, req *fabo.CreateFbShopTagRequest) (*fabo.FbShopTagResponse, error) {
	cmd := &fbusering.CreateShopTagCommand{
		Args: &fbusering.CreateShopTagArgs{
			Name:   req.Name,
			Color:  req.Color,
			ShopID: s.SS.Shop().ID,
		},
	}
	if err := s.FBExternalUserAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	tag := cmd.Result
	resp := convertpb.ConvertFbUserringTagToResponseTag(tag)
	return resp, nil
}

func (s *PageService) UpdateTag(ctx context.Context, req *fabo.UpdateFbShopTagRequest) (*fabo.FbShopTagResponse, error) {
	query := &fbusering.GetShopTagQuery{
		Args: &fbusering.GetShopTagArgs{
			ID:     req.ID,
			ShopID: s.SS.Shop().ID,
		},
	}
	if err := s.FBExternalUserQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	cmdUpdate := &fbusering.UpdateShopTagCommand{
		Args: &fbusering.UpdateShopTagArgs{
			Name:  req.Name,
			Color: req.Color,
			ID:    query.Result.ID,
		},
		Result: nil,
	}
	if err := s.FBExternalUserAggr.Dispatch(ctx, cmdUpdate); err != nil {
		return nil, err
	}

	tag := cmdUpdate.Result
	resp := convertpb.ConvertFbUserringTagToResponseTag(tag)
	return resp, nil
}

func (s *PageService) DeleteTag(ctx context.Context, req *fabo.DeleteFbShopTagRequest) (*cm.Empty, error) {
	cmdDelete := &fbusering.DeleteShopTagCommand{
		Args: &fbusering.DeleteShopTagArgs{
			ID:     req.ID,
			ShopID: s.SS.Shop().ID,
		},
	}
	err := s.FBExternalUserAggr.Dispatch(ctx, cmdDelete)
	return &cm.Empty{}, err
}

func (s *PageService) ListTag(ctx context.Context, req *fabo.ListFbShopTagRequest) (*fabo.ListFbShopTagResponse, error) {
	shopID := s.SS.Shop().ID
	queryListTag := &fbusering.ListShopTagQuery{
		Args: &fbusering.ListShopTagArgs{ShopID: shopID},
	}
	if err := s.FBExternalUserQuery.Dispatch(ctx, queryListTag); err != nil {
		return nil, err
	}

	tags := queryListTag.Result
	if len(tags) == 0 {
		for _, _tag := range defaultTagTemplate {
			cmd := &fbusering.CreateShopTagCommand{
				Args: &fbusering.CreateShopTagArgs{
					Name:   _tag.Name,
					Color:  _tag.Color,
					ShopID: shopID,
				},
			}
			_ = s.FBExternalUserAggr.Dispatch(ctx, cmd)
			tags = append(tags, cmd.Result)
		}
	}
	resp := &fabo.ListFbShopTagResponse{ShopTags: convertpb.ConvertFbUserringTagsToResponseTags(tags)}
	return resp, nil
}
