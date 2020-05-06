package convertpb

import (
	"fmt"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/top/int/fabo"
)

func PbFbCustomerConversation(m *fbmessaging.FbCustomerConversation) *fabo.FbCustomerConversation {
	if m == nil {
		return nil
	}
	return &fabo.FbCustomerConversation{
		ID:                     m.ID,
		FbPageID:               m.FbPageID,
		ExternalID:             m.ExternalID,
		ExternalUserID:         m.ExternalUserID,
		ExternalUserName:       m.ExternalUserName,
		IsRead:                 m.IsRead,
		PostAttachments:        PbFbCustomerConversationPostAttachments(m.PostAttachments),
		Type:                   m.Type,
		ExternalUserPictureURL: GenerateFacebookUserPicture(m.ExternalUserID),
		LastMessage:            m.LastMessage,
		LastMessageAt:          m.LastMessageAt,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}

func PbFbCustomerConversations(ms []*fbmessaging.FbCustomerConversation) []*fabo.FbCustomerConversation {
	res := make([]*fabo.FbCustomerConversation, len(ms))
	for i, m := range ms {
		res[i] = PbFbCustomerConversation(m)
	}
	return res
}

func PbFbCustomerConversationPostAttachment(m *fbmessaging.PostAttachment) *fabo.PostAttachment {
	if m == nil {
		return nil
	}
	var media *fabo.PostAttachmentMedia
	if m.Media != nil {
		media = &fabo.PostAttachmentMedia{
			Height: m.Media.Height,
			Width:  m.Media.Width,
			Src:    m.Media.Src,
		}
	}
	return &fabo.PostAttachment{
		Media: media,
		Type:  m.Type,
	}
}

func PbFbCustomerConversationPostAttachments(ms []*fbmessaging.PostAttachment) []*fabo.PostAttachment {
	res := make([]*fabo.PostAttachment, len(ms))
	for i, m := range ms {
		res[i] = PbFbCustomerConversationPostAttachment(m)
	}
	return res
}

func GenerateFacebookUserPicture(userID string) string {
	return fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", userID)
}

func PbFbExternalMessage(m *fbmessaging.FbExternalMessage) *fabo.FbExternalMessage {
	if m == nil {
		return nil
	}
	return &fabo.FbExternalMessage{
		ID:                     m.ID,
		FbConversationID:       m.FbConversationID,
		ExternalConversationID: m.ExternalConversationID,
		FbPageID:               m.FbPageID,
		ExternalID:             m.ExternalID,
		ExternalMessage:        m.ExternalMessage,
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
	return &fabo.FbObjectTo{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
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
	return &fabo.FbObjectFrom{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
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
	return &fabo.FbMessageAttachment{
		ID:        m.ID,
		ImageData: imageData,
		MimeType:  m.MimeType,
		Name:      m.Name,
		Size:      m.Size,
	}
}

func PbFbMessageAttachments(ms []*fbmessaging.FbMessageAttachment) []*fabo.FbMessageAttachment {
	res := make([]*fabo.FbMessageAttachment, len(ms))
	for i, m := range ms {
		res[i] = PbFbMessageAttachment(m)
	}
	return res
}
