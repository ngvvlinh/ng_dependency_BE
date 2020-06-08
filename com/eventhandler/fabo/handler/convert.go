package handler

import (
	"o.o/api/fabo/fbmessaging"
	exttypes "o.o/api/top/external/types"
	"o.o/backend/com/eventhandler/fabo/types"
	"o.o/capi/dot"
)

func PbFbExternalComment(fbcomment *fbmessaging.FbExternalComment) *exttypes.FbExternalComment {
	if fbcomment == nil {
		return nil
	}
	return &exttypes.FbExternalComment{
		ID:                   fbcomment.ID,
		ExternalUserID:       dot.String(fbcomment.ExternalUserID),
		ExternalParentID:     dot.String(fbcomment.ExternalParentID),
		ExternalParentUserID: dot.String(fbcomment.ExternalParentUserID),
		ExternalMessage:      dot.String(fbcomment.ExternalMessage),
		ExternalCommentCount: dot.Int(fbcomment.ExternalCommentCount),
		ExternalFrom:         PbFbExternalFrom(fbcomment.ExternalFrom),
		ExternalAttachment:   PbFbCommentAttachment(fbcomment.ExternalAttachment),
		ExternalCreatedTime:  fbcomment.ExternalCreatedTime,
		ExternalID:           dot.String(fbcomment.ExternalUserID),
		CreatedAt:            dot.Time(fbcomment.CreatedAt),
		UpdatedAt:            dot.Time(fbcomment.UpdatedAt),
	}
}

func PbFbExternalCommentEvent(fbcomment *fbmessaging.FbExternalComment, op string) *types.PgEventComment {
	if fbcomment == nil {
		return nil
	}
	return &types.PgEventComment{
		Op:             op,
		FbEventComment: PbFbExternalComment(fbcomment),
	}
}

func PbFbCommentAttachment(fbcommentattachment *fbmessaging.CommentAttachment) *exttypes.CommentAttachment {
	if fbcommentattachment == nil {
		return nil
	}
	return &exttypes.CommentAttachment{
		Media: &exttypes.ImageMediaDataSubAttachment{
			Image: PbFbMediaDataSubAttachment(fbcommentattachment.Media.Image),
		},
		Target: PbFbTargetDataSubAttachment(fbcommentattachment.Target),
		Title:  dot.String(fbcommentattachment.Title),
		Type:   dot.String(fbcommentattachment.Type),
		URL:    dot.String(fbcommentattachment.URL),
	}
}

func PbFbExternalFrom(fbfrom *fbmessaging.FbObjectFrom) *exttypes.FbObjectFrom {
	if fbfrom == nil {
		return nil
	}
	return &exttypes.FbObjectFrom{
		ID:    dot.String(fbfrom.ID),
		Name:  dot.String(fbfrom.Name),
		Email: dot.String(fbfrom.Email),
	}
}

func PbFbExternalTo(fbto *fbmessaging.FbObjectTo) *exttypes.FbObjectTo {
	if fbto == nil {
		return nil
	}
	return &exttypes.FbObjectTo{
		ID:    dot.String(fbto.ID),
		Name:  dot.String(fbto.Name),
		Email: dot.String(fbto.Email),
	}
}

func PbFbExternalTos(fbtos []*fbmessaging.FbObjectTo) []*exttypes.FbObjectTo {
	out := make([]*exttypes.FbObjectTo, len(fbtos))
	for i, fbto := range fbtos {
		out[i] = PbFbExternalTo(fbto)
	}
	return out
}

func PbFbPostAttachments(attachments []*fbmessaging.PostAttachment) []*exttypes.PostAttachment {
	out := make([]*exttypes.PostAttachment, len(attachments))
	for i, attachment := range attachments {
		out[i] = PbFbPostAttachment(attachment)
	}
	return out
}

func PbFbPostAttachment(attachment *fbmessaging.PostAttachment) *exttypes.PostAttachment {
	if attachment == nil {
		return nil
	}
	return &exttypes.PostAttachment{
		MediaType:      dot.String(attachment.MediaType),
		Type:           dot.String(attachment.Type),
		SubAttachments: PbFbSubAttachments(attachment.SubAttachments),
	}
}

func PbFbSubAttachments(subattachments []*fbmessaging.SubAttachment) []*exttypes.SubAttachment {
	out := make([]*exttypes.SubAttachment, len(subattachments))
	for i, sub := range subattachments {
		out[i] = PbFbSubAttachment(sub)
	}
	return out
}

func PbFbSubAttachment(subattachment *fbmessaging.SubAttachment) *exttypes.SubAttachment {
	if subattachment == nil {
		return nil
	}
	return &exttypes.SubAttachment{
		Media:  PbFbMediaDataSubAttachment(subattachment.Media),
		Target: PbFbTargetDataSubAttachment(subattachment.Target),
		Type:   dot.String(subattachment.Type),
		URL:    dot.String(subattachment.URL),
	}
}

func PbFbMediaDataSubAttachment(mediadata *fbmessaging.MediaDataSubAttachment) *exttypes.MediaDataSubAttachment {
	if mediadata == nil {
		return nil
	}
	return &exttypes.MediaDataSubAttachment{
		Height: dot.Int(mediadata.Height),
		Width:  dot.Int(mediadata.Width),
		Src:    dot.String(mediadata.Src),
	}
}

func PbFbTargetDataSubAttachment(targetdata *fbmessaging.TargetDataSubAttachment) *exttypes.TargetDataSubAttachment {
	if targetdata == nil {
		return nil
	}
	return &exttypes.TargetDataSubAttachment{
		ID:  dot.String(targetdata.ID),
		URL: dot.String(targetdata.URL),
	}
}

func PbFbMessageAttachment(fbattachment *fbmessaging.FbMessageAttachment) *exttypes.FbMessageAttachment {
	if fbattachment == nil {
		return nil
	}
	return &exttypes.FbMessageAttachment{
		ID:        dot.String(fbattachment.ID),
		ImageData: PbFbMessageAttachmentImage(fbattachment.ImageData),
		MimeType:  dot.String(fbattachment.MimeType),
		Name:      dot.String(fbattachment.Name),
		Size:      dot.Int(fbattachment.Size),
	}
}

func PbFbMessageAttachments(messageattachments []*fbmessaging.FbMessageAttachment) []*exttypes.FbMessageAttachment {
	out := make([]*exttypes.FbMessageAttachment, len(messageattachments))
	for i, messageattachment := range messageattachments {
		out[i] = PbFbMessageAttachment(messageattachment)
	}
	return out
}

func PbFbMessageAttachmentImage(fbmessageomage *fbmessaging.FbMessageAttachmentImageData) *exttypes.FbMessageAttachmentImageData {
	if fbmessageomage == nil {
		return nil
	}
	return &exttypes.FbMessageAttachmentImageData{
		Width:           dot.Int(fbmessageomage.Width),
		Height:          dot.Int(fbmessageomage.Height),
		MaxWidth:        dot.Int(fbmessageomage.MaxWidth),
		MaxHeight:       dot.Int(fbmessageomage.MaxHeight),
		URL:             dot.String(fbmessageomage.URL),
		PreviewURL:      dot.String(fbmessageomage.PreviewURL),
		ImageType:       dot.Int(fbmessageomage.Width),
		RenderAsSticker: dot.Bool(fbmessageomage.RenderAsSticker),
	}
}

func PbFbExternalPost(fbpost *fbmessaging.FbExternalPost) *exttypes.FbExternalPost {
	if fbpost == nil {
		return nil
	}
	return &exttypes.FbExternalPost{
		Id:                  fbpost.ID,
		ExternalID:          dot.String(fbpost.ExternalID),
		ExternalParentID:    dot.String(fbpost.ExternalParentID),
		ExternalFrom:        PbFbExternalFrom(fbpost.ExternalFrom),
		ExternalPicture:     dot.String(fbpost.ExternalPicture),
		ExternalPageID:      dot.String(fbpost.ExternalPageID),
		ExternalIcon:        dot.String(fbpost.ExternalIcon),
		ExternalMessage:     dot.String(fbpost.ExternalMessage),
		ExternalAttachments: PbFbPostAttachments(fbpost.ExternalAttachments),
		ExternalCreatedTime: fbpost.ExternalCreatedTime,
		CreatedAt:           dot.Time(fbpost.CreatedAt),
		UpdatedAt:           dot.Time(fbpost.UpdatedAt),
	}
}

func PbFbExternalConversationEvent(fbConversation *fbmessaging.FbExternalConversation, op string) *types.PgEventConversation {
	if fbConversation == nil {
		return nil
	}
	return &types.PgEventConversation{
		FbEventConversation: PbFbExternalConversation(fbConversation),
		Op:                  op,
	}
}

func PbFbExternalConversation(fbConversation *fbmessaging.FbExternalConversation) *exttypes.FbExternalConversation {
	if fbConversation == nil {
		return nil
	}
	return &exttypes.FbExternalConversation{
		ID:                   fbConversation.ID,
		ExternalID:           dot.String(fbConversation.ExternalID),
		ExternalPageID:       dot.String(fbConversation.ExternalPageID),
		PSID:                 dot.String(fbConversation.PSID),
		ExternalUserID:       dot.String(fbConversation.ExternalUserID),
		ExternalUserName:     dot.String(fbConversation.ExternalUserName),
		ExternalLink:         dot.String(fbConversation.ExternalLink),
		ExternalUpdatedTime:  fbConversation.ExternalUpdatedTime,
		ExternalMessageCount: dot.Int(fbConversation.ExternalMessageCount),
		CreatedAt:            fbConversation.CreatedAt,
		UpdatedAt:            fbConversation.UpdatedAt,
	}
}

func PbFbCustomerConversation(fbCustomerConversation *fbmessaging.FbCustomerConversation) *exttypes.FbCustomerConversation {
	if fbCustomerConversation == nil {
		return nil
	}
	return &exttypes.FbCustomerConversation{
		ID:                         fbCustomerConversation.ID,
		ExternalPageID:             dot.String(fbCustomerConversation.ExternalPageID),
		ExternalID:                 dot.String(fbCustomerConversation.ExternalID),
		ExternalUserID:             dot.String(fbCustomerConversation.ExternalUserID),
		ExternalUserName:           dot.String(fbCustomerConversation.ExternalUserName),
		ExternalFrom:               PbFbExternalFrom(fbCustomerConversation.ExternalFrom),
		IsRead:                     dot.Bool(fbCustomerConversation.IsRead),
		Type:                       fbCustomerConversation.Type,
		ExternalPostAttachments:    PbFbPostAttachments(fbCustomerConversation.ExternalPostAttachments),
		ExternalCommentAttachment:  PbFbCommentAttachment(fbCustomerConversation.ExternalCommentAttachment),
		ExternalMessageAttachments: PbFbMessageAttachments(fbCustomerConversation.ExternalMessageAttachments),
		LastMessage:                dot.String(fbCustomerConversation.LastMessage),
		LastMessageAt:              fbCustomerConversation.LastMessageAt,
		CreatedAt:                  fbCustomerConversation.CreatedAt,
		UpdatedAt:                  fbCustomerConversation.UpdatedAt,
	}
}

func PbFbCustomerConversationEvent(fbConversation *fbmessaging.FbCustomerConversation, op string) *types.PgEventCustomerConversation {
	if fbConversation == nil {
		return nil
	}
	return &types.PgEventCustomerConversation{
		FbEventCustomerConversation: PbFbCustomerConversation(fbConversation),
		Op:                          op,
	}
}

func PbFbExternalMessage(fbmessage *fbmessaging.FbExternalMessage) *exttypes.FbExternalMessage {
	if fbmessage == nil {
		return nil
	}
	return &exttypes.FbExternalMessage{
		ID:                     fbmessage.ID,
		ExternalConversationID: dot.String(fbmessage.ExternalConversationID),
		ExternalID:             dot.String(fbmessage.ExternalID),
		ExternalMessage:        dot.String(fbmessage.ExternalMessage),
		ExternalPageID:         dot.String(fbmessage.ExternalPageID),
		ExternalSticker:        dot.String(fbmessage.ExternalSticker),
		ExternalTo:             PbFbExternalTos(fbmessage.ExternalTo),
		ExternalFrom:           PbFbExternalFrom(fbmessage.ExternalFrom),
		ExternalAttachments:    PbFbMessageAttachments(fbmessage.ExternalAttachments),
		ExternalCreatedTime:    fbmessage.ExternalCreatedTime,
		CreatedAt:              fbmessage.CreatedAt,
		UpdatedAt:              fbmessage.UpdatedAt,
	}
}

func PbFbExternalMessageEvent(fbmessage *fbmessaging.FbExternalMessage, op string) *types.PgEventMessage {
	if fbmessage == nil {
		return nil
	}
	return &types.PgEventMessage{
		Op:             op,
		FbEventMessage: PbFbExternalMessage(fbmessage),
	}
}