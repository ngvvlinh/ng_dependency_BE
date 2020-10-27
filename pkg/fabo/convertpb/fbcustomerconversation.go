package convertpb

import (
	"o.o/api/fabo/fbmessagetemplate"
	"o.o/api/fabo/fbmessaging"
	"o.o/api/top/int/fabo"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
)

func PbFbCustomerConversation(m *fbmessaging.FbCustomerConversation) *fabo.FbCustomerConversation {
	if m == nil {
		return nil
	}
	result := &fabo.FbCustomerConversation{
		ID:                        m.ID,
		ExternalPageID:            m.ExternalPageID,
		ExternalID:                m.ExternalID,
		ExternalUserID:            m.ExternalUserID,
		ExternalUserName:          m.ExternalUserName,
		ExternalFrom:              PbFbObjectFrom(m.ExternalFrom),
		IsRead:                    m.IsRead,
		ExternalPostAttachments:   PbPostAttachments(m.ExternalPostAttachments),
		ExternalCommentAttachment: PbCommentAttachment(m.ExternalCommentAttachment),
		Type:                      m.Type.String(),
		ExternalUserPictureURL:    fbclientconvert.GenerateFacebookUserPicture(m.ExternalUserID),
		LastMessage:               m.LastMessage,
		LastMessageAt:             m.LastMessageAt,
		LastCustomerMessageAt:     m.LastCustomerMessageAt,
		CreatedAt:                 m.CreatedAt,
		UpdatedAt:                 m.UpdatedAt,
	}
	if m.ExternalUserPictureURL == "" {
		result.ExternalUserPictureURL = fbclientconvert.GenerateFacebookUserPicture(m.ExternalUserID)
	} else {
		result.ExternalUserPictureURL = m.ExternalUserPictureURL
	}
	return result
}

func PbFbCustomerConversations(ms []*fbmessaging.FbCustomerConversation) []*fabo.FbCustomerConversation {
	res := make([]*fabo.FbCustomerConversation, len(ms))
	for i, m := range ms {
		res[i] = PbFbCustomerConversation(m)
	}
	return res
}

func PbFbExternalMessage(m *fbmessaging.FbExternalMessage) *fabo.FbExternalMessage {
	if m == nil {
		return nil
	}
	return &fabo.FbExternalMessage{
		ID:                     m.ID,
		ExternalConversationID: m.ExternalConversationID,
		ExternalPageID:         m.ExternalPageID,
		ExternalID:             m.ExternalID,
		ExternalMessage:        m.ExternalMessage,
		ExternalSticker:        m.ExternalSticker,
		ExternalTo:             PbFbObjectsTo(m.ExternalTo),
		ExternalFrom:           PbFbObjectFrom(m.ExternalFrom),
		ExternalAttachments:    PbFbMessageAttachments(m.ExternalAttachments),
		ExternalCreatedTime:    m.ExternalCreatedTime,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}

func PbFbExternalMessages(ms []*fbmessaging.FbExternalMessage) []*fabo.FbExternalMessage {
	res := make([]*fabo.FbExternalMessage, len(ms))
	for i, m := range ms {
		res[i] = PbFbExternalMessage(m)
	}
	return res
}

func PbFbObjectTo(m *fbmessaging.FbObjectTo) *fabo.FbObjectTo {
	if m == nil {
		return nil
	}
	result := &fabo.FbObjectTo{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
	if m.ImageURL == "" {
		result.ExternalUserPictureURL = fbclientconvert.GenerateFacebookUserPicture(m.ID)
	} else {
		result.ExternalUserPictureURL = m.ImageURL
	}
	return result
}

func PbFbObjectsTo(ms []*fbmessaging.FbObjectTo) []*fabo.FbObjectTo {
	res := make([]*fabo.FbObjectTo, len(ms))
	for i, m := range ms {
		res[i] = PbFbObjectTo(m)
	}
	return res
}

func PbFbObjectFrom(m *fbmessaging.FbObjectFrom) *fabo.FbObjectFrom {
	if m == nil {
		return nil
	}
	result := &fabo.FbObjectFrom{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
	if m.ImageURL == "" {
		result.ExternalUserPictureURL = fbclientconvert.GenerateFacebookUserPicture(m.ID)
	} else {
		result.ExternalUserPictureURL = m.ImageURL
	}
	return result
}

func PbFbMessageAttachment(m *fbmessaging.FbMessageAttachment) *fabo.FbMessageAttachment {
	if m == nil {
		return nil
	}
	var imageData *fabo.FbMessageAttachmentImageData
	if m.ImageData != nil {
		imageData = &fabo.FbMessageAttachmentImageData{
			Width:           m.ImageData.Width,
			Height:          m.ImageData.Height,
			MaxWidth:        m.ImageData.MaxWidth,
			MaxHeight:       m.ImageData.MaxHeight,
			URL:             m.ImageData.URL,
			PreviewURL:      m.ImageData.PreviewURL,
			ImageType:       m.ImageData.ImageType,
			RenderAsSticker: m.ImageData.RenderAsSticker,
		}
	}

	var videoData *fabo.FbMessageDataAttachmentVideoData
	if m.VideoData != nil {
		videoData = &fabo.FbMessageDataAttachmentVideoData{
			Width:      m.VideoData.Width,
			Height:     m.VideoData.Height,
			Length:     m.VideoData.Length,
			VideoType:  m.VideoData.VideoType,
			URL:        m.VideoData.URL,
			PreviewURL: m.VideoData.PreviewURL,
			Rotation:   m.VideoData.Rotation,
		}
	}

	return &fabo.FbMessageAttachment{
		ID:        m.ID,
		ImageData: imageData,
		MimeType:  m.MimeType,
		Name:      m.Name,
		Size:      m.Size,
		VideoData: videoData,
		FileURL:   m.FileURL,
	}
}

func PbFbMessageAttachments(ms []*fbmessaging.FbMessageAttachment) []*fabo.FbMessageAttachment {
	res := make([]*fabo.FbMessageAttachment, len(ms))
	for i, m := range ms {
		res[i] = PbFbMessageAttachment(m)
	}
	return res
}

func PbFbExternalPost(m *fbmessaging.FbExternalPost) *fabo.FbExternalPost {
	if m == nil {
		return nil
	}
	return &fabo.FbExternalPost{
		ID:                  m.ID,
		ExternalPageID:      m.ExternalPageID,
		ExternalID:          m.ExternalID,
		ExternalParentID:    m.ExternalParentID,
		ExternalFrom:        PbFbObjectFrom(m.ExternalFrom),
		ExternalPicture:     m.ExternalPicture,
		ExternalIcon:        m.ExternalIcon,
		ExternalMessage:     m.ExternalMessage,
		ExternalParent:      PbFbExternalPost(m.ExternalParent),
		ExternalAttachments: PbPostAttachments(m.ExternalAttachments),
		ExternalCreatedTime: m.ExternalCreatedTime,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func PbPostAttachments(ms []*fbmessaging.PostAttachment) []*fabo.PostAttachment {
	res := make([]*fabo.PostAttachment, len(ms))
	for i, m := range ms {
		res[i] = PbPostAttachment(m)
	}
	return res
}

func PbPostAttachment(m *fbmessaging.PostAttachment) *fabo.PostAttachment {
	if m == nil {
		return nil
	}
	var subAttachments []*fabo.SubAttachment
	for _, subAttachment := range m.SubAttachments {
		var media *fabo.MediaDataSubAttachment
		var target *fabo.TargetDataSubAttachment
		if subAttachment.Media != nil {
			media = &fabo.MediaDataSubAttachment{
				Width:  subAttachment.Media.Width,
				Height: subAttachment.Media.Height,
				Src:    subAttachment.Media.Src,
			}
		}
		if subAttachment.Target != nil {
			target = &fabo.TargetDataSubAttachment{
				ID:  subAttachment.Target.ID,
				URL: subAttachment.Target.URL,
			}
		}
		subAttachments = append(subAttachments, &fabo.SubAttachment{
			Media:  media,
			Target: target,
			Type:   subAttachment.Type,
			URL:    subAttachment.URL,
		})
	}

	return &fabo.PostAttachment{
		MediaType:      m.MediaType,
		Type:           m.Type,
		Media:          PbMediaForCustomerConversation(m.Media),
		SubAttachments: subAttachments,
	}
}

func PbMediaForCustomerConversation(m *fbmessaging.MediaPostAttachment) *fabo.MediaPostAttachment {
	if m == nil {
		return nil
	}
	return &fabo.MediaPostAttachment{
		Image: PbImageForCustomerConversation(m.Image),
	}
}

func PbImageForCustomerConversation(m *fbmessaging.ImageMediaPostAttachment) *fabo.ImageMediaPostAttachment {
	if m == nil {
		return nil
	}
	return &fabo.ImageMediaPostAttachment{
		Height: m.Height,
		Width:  m.Width,
		Src:    m.Src,
	}
}

func PbFbExternalComments(ms []*fbmessaging.FbExternalComment) []*fabo.FbExternalComment {
	res := make([]*fabo.FbExternalComment, len(ms))
	for i, m := range ms {
		res[i] = PbFbExternalComment(m)
	}
	return res
}

func PbFbExternalComment(m *fbmessaging.FbExternalComment) *fabo.FbExternalComment {
	if m == nil {
		return nil
	}
	return &fabo.FbExternalComment{
		ID:                   m.ID,
		ExternalPostID:       m.ExternalPostID,
		ExternalPageID:       m.ExternalPageID,
		ExternalID:           m.ExternalID,
		ExternalUserID:       m.ExternalUserID,
		ExternalParentID:     m.ExternalParentID,
		ExternalParentUserID: m.ExternalParentUserID,
		ExternalMessage:      m.ExternalMessage,
		ExternalCommentCount: m.ExternalCommentCount,
		ExternalFrom:         PbFbObjectFrom(m.ExternalFrom),
		ExternalAttachment:   PbCommentAttachment(m.ExternalAttachment),
		ExternalCreatedTime:  m.ExternalCreatedTime,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func PbCommentAttachment(m *fbmessaging.CommentAttachment) *fabo.CommentAttachment {
	if m == nil {
		return nil
	}
	var media *fabo.ImageMediaDataSubAttachment
	var target *fabo.TargetDataSubAttachment
	if m.Media != nil && m.Media.Image != nil {
		media = &fabo.ImageMediaDataSubAttachment{
			Image: &fabo.MediaDataSubAttachment{
				Width:  m.Media.Image.Width,
				Height: m.Media.Image.Height,
				Src:    m.Media.Image.Src,
			},
		}
	}
	if m.Target != nil {
		target = &fabo.TargetDataSubAttachment{
			ID:  m.Target.ID,
			URL: m.Target.URL,
		}
	}
	return &fabo.CommentAttachment{
		Media:  media,
		Target: target,
		Title:  m.Title,
		Type:   m.Type,
		URL:    m.URL,
	}
}

func PbFbObjectParent(m *fbmessaging.FbObjectParent) *fabo.FbObjectParent {
	if m == nil {
		return nil
	}
	return &fabo.FbObjectParent{
		CreatedTime: m.CreatedTime,
		From:        PbFbObjectFrom(m.From),
		Message:     m.Message,
		ID:          m.ID,
	}
}

func FbMessageTemplate(m *fbmessagetemplate.FbMessageTemplate) *fabo.MessageTemplate {
	if m == nil {
		return nil
	}

	return &fabo.MessageTemplate{
		ID:        m.ID,
		Template:  m.Template,
		ShortCode: m.ShortCode,
		ShopID:    m.ShopID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FbMessageTemplates(ms []*fbmessagetemplate.FbMessageTemplate) []*fabo.MessageTemplate {
	if ms == nil {
		return nil
	}

	var result []*fabo.MessageTemplate
	for _, template := range ms {
		result = append(result, FbMessageTemplate(template))
	}
	return result
}
