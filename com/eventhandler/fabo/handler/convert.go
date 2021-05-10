package handler

import (
	"o.o/api/fabo/fbmessaging"
	exttypes "o.o/api/top/external/types"
	"o.o/backend/com/eventhandler/fabo/types"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	"o.o/capi/dot"
)

func PbFbExternalComment(fbcomment *fbmessaging.FbExternalComment, fbParentComment *fbmessaging.FbExternalComment) *exttypes.FbExternalComment {
	if fbcomment == nil {
		return nil
	}
	return &exttypes.FbExternalComment{
		ID:                   fbcomment.ID,
		ExternalPageID:       dot.String(fbcomment.ExternalPageID),
		ExternalUserID:       dot.String(fbcomment.ExternalUserID),
		ExternalParentID:     dot.String(fbcomment.ExternalParentID),
		ExternalParent:       PbFbExternalComment(fbParentComment, nil),
		ExternalParentUserID: dot.String(fbcomment.ExternalParentUserID),
		ExternalMessage:      dot.String(fbcomment.ExternalMessage),
		ExternalCommentCount: dot.Int(fbcomment.ExternalCommentCount),
		ExternalFrom:         PbFbExternalFrom(fbcomment.ExternalFrom),
		ExternalAttachment:   PbFbCommentAttachment(fbcomment.ExternalAttachment),
		ExternalCreatedTime:  fbcomment.ExternalCreatedTime,
		ExternalID:           dot.String(fbcomment.ExternalID),
		ExternalPostID:       dot.String(fbcomment.ExternalPostID),
		IsLiked:              dot.Bool(fbcomment.IsLiked),
		IsHidden:             dot.Bool(fbcomment.IsHidden),
		IsPrivateReplied:     dot.Bool(fbcomment.IsPrivateReplied),
		CreatedBy:            fbcomment.CreatedBy.Wrap(),
		CreatedAt:            dot.Time(fbcomment.CreatedAt),
		UpdatedAt:            dot.Time(fbcomment.UpdatedAt),
	}
}

func PbFbExternalCommentEvent(fbcomment *fbmessaging.FbExternalComment, fbParentComment *fbmessaging.FbExternalComment, op string) *types.PgEventComment {
	if fbcomment == nil {
		return nil
	}
	return &types.PgEventComment{
		Op:             op,
		FbEventComment: PbFbExternalComment(fbcomment, fbParentComment),
	}
}

func PbFbCommentAttachment(fbcommentattachment *fbmessaging.CommentAttachment) *exttypes.CommentAttachment {
	if fbcommentattachment == nil {
		return nil
	}

	commentAttachment := &exttypes.CommentAttachment{
		Target: PbFbTargetDataSubAttachment(fbcommentattachment.Target),
		Title:  dot.String(fbcommentattachment.Title),
		Type:   dot.String(fbcommentattachment.Type),
		URL:    dot.String(fbcommentattachment.URL),
	}

	if fbcommentattachment.Media != nil {
		commentAttachment.Media = &exttypes.ImageMediaDataSubAttachment{
			Image: PbFbMediaDataSubAttachment(fbcommentattachment.Media.Image),
		}
	}

	return commentAttachment
}

func PbFbExternalParent(fbExternalParent *fbmessaging.FbObjectParent) *exttypes.FbObjectParent {
	if fbExternalParent == nil {
		return nil
	}
	return &exttypes.FbObjectParent{
		CreatedTime: fbExternalParent.CreatedTime,
		From:        PbFbExternalFrom(fbExternalParent.From),
		Message:     fbExternalParent.Message,
		ID:          fbExternalParent.ID,
	}
}

func PbFbExternalFrom(fbfrom *fbmessaging.FbObjectFrom) *exttypes.FbObjectFrom {
	if fbfrom == nil {
		return nil
	}
	result := &exttypes.FbObjectFrom{
		ID:    dot.String(fbfrom.ID),
		Name:  dot.String(fbfrom.Name),
		Email: dot.String(fbfrom.Email),
	}
	if fbfrom.ImageURL == "" {
		result.ExternalUserPictureURL = dot.String(fbclientconvert.GenerateFacebookUserPicture(fbfrom.ID))
	} else {
		result.ExternalUserPictureURL = dot.String(fbfrom.ImageURL)
	}
	return result
}

func PbFbExternalTo(fbto *fbmessaging.FbObjectTo) *exttypes.FbObjectTo {
	if fbto == nil {
		return nil
	}
	result := &exttypes.FbObjectTo{
		ID:    dot.String(fbto.ID),
		Name:  dot.String(fbto.Name),
		Email: dot.String(fbto.Email),
	}
	if fbto.ImageURL == "" {
		result.ExternalUserPictureURL = dot.String(fbclientconvert.GenerateFacebookUserPicture(fbto.ID))
	} else {
		result.ExternalUserPictureURL = dot.String(fbto.ImageURL)
	}
	return result
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
		Media:          PbFbPostAttachmentMedia(attachment.Media),
	}
}

func PbFbPostAttachmentMedia(media *fbmessaging.MediaPostAttachment) *exttypes.MediaPostAttachment {
	if media == nil || media.Image == nil {
		return nil
	}
	return &exttypes.MediaPostAttachment{
		Image: &exttypes.ImageMediaPostAttachment{
			Height: media.Image.Height,
			Width:  media.Image.Width,
			Src:    media.Image.Src,
		},
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
		ExternalUserPictureURL:     dot.String(fbCustomerConversation.ExternalUserPictureURL),
		ExternalFrom:               PbFbExternalFrom(fbCustomerConversation.ExternalFrom),
		IsRead:                     dot.Bool(fbCustomerConversation.IsRead),
		Type:                       fbCustomerConversation.Type,
		ExternalPostAttachments:    PbFbPostAttachments(fbCustomerConversation.ExternalPostAttachments),
		ExternalCommentAttachment:  PbFbCommentAttachment(fbCustomerConversation.ExternalCommentAttachment),
		ExternalMessageAttachments: PbFbMessageAttachments(fbCustomerConversation.ExternalMessageAttachments),
		LastMessage:                dot.String(fbCustomerConversation.LastMessage),
		LastMessageAt:              fbCustomerConversation.LastMessageAt,
		LastCustomerMessageAt:      fbCustomerConversation.LastCustomerMessageAt,
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
		CreatedBy:              fbmessage.CreatedBy.Wrap(),
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

func PbFbExternalPostEvent(fbPost *fbmessaging.FbExternalPost, op string) *types.PgEventPost {
	if fbPost == nil {
		return nil
	}
	return &types.PgEventPost{
		Op:          op,
		FbEventPost: PbFbExternalPost(fbPost),
	}
}

func PbFbExternalPost(fbPost *fbmessaging.FbExternalPost) *exttypes.FbExternalPost {
	if fbPost == nil {
		return nil
	}
	return &exttypes.FbExternalPost{
		Id:                      fbPost.ID,
		ExternalPageID:          dot.String(fbPost.ExternalPageID),
		ExternalID:              dot.String(fbPost.ExternalID),
		ExternalUserID:          dot.String(fbPost.ExternalUserID),
		ExternalParentID:        dot.String(fbPost.ExternalParentID),
		ExternalFrom:            PbFbExternalFrom(fbPost.ExternalFrom),
		ExternalPicture:         dot.String(fbPost.ExternalPicture),
		ExternalIcon:            dot.String(fbPost.ExternalIcon),
		ExternalMessage:         dot.String(fbPost.ExternalMessage),
		ExternalAttachments:     PbFbPostAttachments(fbPost.ExternalAttachments),
		ExternalCreatedTime:     fbPost.ExternalCreatedTime,
		CreatedAt:               dot.Time(fbPost.CreatedAt),
		UpdatedAt:               dot.Time(fbPost.UpdatedAt),
		FeedType:                fbPost.FeedType,
		StatusType:              fbPost.StatusType,
		Type:                    fbPost.Type,
		TotalComments:           dot.Int(fbPost.TotalComments),
		TotalReactions:          dot.Int(fbPost.TotalReactions),
		IsLiveVideo:             dot.Bool(fbPost.IsLiveVideo),
		ExternalLiveVideoStatus: dot.String(fbPost.ExternalLiveVideoStatus),
		LiveVideoStatus:         fbPost.LiveVideoStatus,
	}
}
