package _all

import (
	"o.o/api/supporting/ticket"
	"o.o/api/top/int/shop/types"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_core_TicketComment_to_api_TicketComment(in *ticket.TicketComment) *types.TicketComment {
	if in == nil {
		return nil
	}
	res := &types.TicketComment{
		ID:        in.ID,
		TicketID:  in.TicketID,
		CreatedBy: in.CreatedBy,
		AccountID: in.AccountID,
		ParentID:  in.ParentID,
		Message:   in.Message,
		ImageUrls: in.ImageUrls,
		DeletedBy: in.DeletedBy,
		CreatedAt: cmapi.PbTime(in.CreatedAt),
		UpdatedAt: cmapi.PbTime(in.UpdatedAt),
		From: &types.TicketFrom{
			ID:     in.CreatedBy,
			Name:   in.CreatedName,
			Source: in.CreatedSource,
		},
	}
	// backward-compatible
	if len(in.ImageUrls) > 0 {
		res.ImageUrl = in.ImageUrls[0]
	}
	return res
}

func Convert_core_TicketComments_to_api_TicketComments(items []*ticket.TicketComment) []*types.TicketComment {
	result := make([]*types.TicketComment, len(items))
	for i, item := range items {
		result[i] = Convert_core_TicketComment_to_api_TicketComment(item)
	}
	return result
}

func Convert_core_TicketLabel_to_api_TicketLabel(in *ticket.TicketLabel) *types.TicketLabel {
	if in == nil {
		return nil
	}
	return &types.TicketLabel{
		ID:       in.ID,
		Name:     in.Name,
		Code:     in.Code,
		Color:    in.Color,
		ParentID: in.ParentID,
		Children: Convert_core_TicketLabels_to_api_TicketLabels(in.Children),
	}
}

func Convert_core_TicketLabels_to_api_TicketLabels(items []*ticket.TicketLabel) []*types.TicketLabel {
	result := make([]*types.TicketLabel, len(items))
	for i, item := range items {
		result[i] = Convert_core_TicketLabel_to_api_TicketLabel(item)
	}
	return result
}

func Convert_core_Ticket_to_api_Ticket(in *ticket.Ticket) *types.Ticket {
	if in == nil {
		return nil
	}
	return &types.Ticket{
		ID:              in.ID,
		Code:            in.Code,
		AssignedUserIDs: in.AssignedUserIDs,
		AccountID:       in.AccountID,
		LabelIDs:        in.LabelIDs,
		Title:           in.Title,
		Description:     in.Description,
		Note:            in.Note,
		ExternalID:      in.ExternalID,
		AdminNote:       in.AdminNote,
		RefID:           in.RefID,
		RefType:         in.RefType,
		RefTicketID:     in.RefTicketID.ID,
		Source:          in.Source,
		RefCode:         in.RefCode,
		State:           in.State,
		Status:          in.Status,
		CreatedBy:       in.CreatedBy,
		UpdatedBy:       in.UpdatedBy,
		ConfirmedBy:     in.ConfirmedBy,
		ClosedBy:        in.ClosedBy,
		CreatedAt:       cmapi.PbTime(in.CreatedAt),
		UpdatedAt:       cmapi.PbTime(in.UpdatedAt),
		ConfirmedAt:     cmapi.PbTime(in.ConfirmedAt),
		ClosedAt:        cmapi.PbTime(in.ClosedAt),
		From: &types.TicketFrom{
			ID:     in.CreatedBy,
			Name:   in.CreatedName,
			Source: in.CreatedSource,
		},
	}
}

func Convert_core_Tickets_to_api_Tickets(items []*ticket.Ticket) []*types.Ticket {
	result := make([]*types.Ticket, len(items))
	for i, item := range items {
		result[i] = Convert_core_Ticket_to_api_Ticket(item)
	}
	return result
}
