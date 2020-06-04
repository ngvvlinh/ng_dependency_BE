package convert

import (
	"fmt"

	"o.o/api/fabo/fbmessaging"
	"o.o/backend/com/fabo/pkg/fbclient/model"
)

func ConvertObjectsTo(ins *model.ObjectsTo) []*fbmessaging.FbObjectTo {
	if ins == nil {
		return nil
	}
	var outs []*fbmessaging.FbObjectTo

	for _, in := range ins.Data {
		outs = append(outs, &fbmessaging.FbObjectTo{
			ID:    in.ID,
			Name:  in.Name,
			Email: in.Email,
		})
	}

	return outs
}

func ConvertObjectFrom(in *model.ObjectFrom) *fbmessaging.FbObjectFrom {
	if in == nil {
		return nil
	}
	return &fbmessaging.FbObjectFrom{
		ID:    in.ID,
		Name:  in.Name,
		Email: in.Email,
	}
}

func ConvertObjectsFrom(ins *model.ObjectsFrom) []*fbmessaging.FbObjectFrom {
	if ins == nil {
		return nil
	}
	var outs []*fbmessaging.FbObjectFrom

	for _, in := range ins.Data {
		outs = append(outs, &fbmessaging.FbObjectFrom{
			ID:    in.ID,
			Name:  in.Name,
			Email: in.Email,
		})
	}

	return outs
}

func ConvertMessageDataAttachments(ins []*model.MessageDataAttachment) []*fbmessaging.FbMessageAttachment {
	outs := make([]*fbmessaging.FbMessageAttachment, 0, len(ins))
	for _, in := range ins {
		var imageData *fbmessaging.FbMessageAttachmentImageData
		if in.ImageData != nil {
			imageData = &fbmessaging.FbMessageAttachmentImageData{
				Width:           in.ImageData.Width,
				Height:          in.ImageData.Height,
				MaxWidth:        in.ImageData.MaxWidth,
				MaxHeight:       in.ImageData.MaxHeight,
				URL:             in.ImageData.URL,
				PreviewURL:      in.ImageData.PreviewURL,
				ImageType:       in.ImageData.ImageType,
				RenderAsSticker: in.ImageData.RenderAsSticker,
			}
		}
		var videoData *fbmessaging.FbMessageDataAttachmentVideoData
		if in.VideoData != nil {
			videoData = &fbmessaging.FbMessageDataAttachmentVideoData{
				Width:      in.VideoData.Width,
				Height:     in.VideoData.Height,
				Length:     in.VideoData.Length,
				VideoType:  in.VideoData.VideoType,
				URL:        in.VideoData.URL,
				PreviewURL: in.VideoData.PreviewURL,
				Rotation:   in.VideoData.Rotation,
			}
		}
		outs = append(outs, &fbmessaging.FbMessageAttachment{
			ID:        in.ID,
			ImageData: imageData,
			MimeType:  in.MimeType,
			Name:      in.Name,
			Size:      in.Size,
			VideoData: videoData,
			FileURL:   in.FileURL,
		})
	}
	return outs
}

func ConvertAttachments(in *model.Attachments) []*fbmessaging.PostAttachment {
	if in == nil {
		return nil
	}
	return ConvertDataAttachments(in.Data)
}

func ConvertDataAttachments(ins []*model.DataAttachment) []*fbmessaging.PostAttachment {
	outs := make([]*fbmessaging.PostAttachment, 0, len(ins))
	for _, in := range ins {
		// TODO: Ngoc check type before convert subAttachments
		var subAttachments []*fbmessaging.SubAttachment
		var media *fbmessaging.MediaPostAttachment
		if in.SubAttachments != nil {
			for _, subAttachment := range in.SubAttachments.Data {
				var media *fbmessaging.MediaDataSubAttachment
				if subAttachment.Media != nil && subAttachment.Media.Image != nil {
					media = &fbmessaging.MediaDataSubAttachment{
						Height: subAttachment.Media.Image.Height,
						Width:  subAttachment.Media.Image.Width,
						Src:    subAttachment.Media.Image.Src,
					}
				}
				var target *fbmessaging.TargetDataSubAttachment
				if subAttachment.Target != nil {
					target = &fbmessaging.TargetDataSubAttachment{
						ID:  subAttachment.Target.ID,
						URL: subAttachment.Target.URL,
					}
				}
				subAttachments = append(subAttachments, &fbmessaging.SubAttachment{
					Media:  media,
					Target: target,
					Type:   subAttachment.Type,
					URL:    subAttachment.URL,
				})
			}
		}
		if in.Media != nil {
			if in.Media.Image != nil {
				media = &fbmessaging.MediaPostAttachment{
					Image: &fbmessaging.ImageMediaPostAttachment{
						Height: in.Media.Image.Height,
						Width:  in.Media.Image.Width,
						Src:    in.Media.Image.Src,
					},
				}
			}
		}
		outs = append(outs, &fbmessaging.PostAttachment{
			Media:          media,
			MediaType:      in.MediaType,
			Type:           in.Type,
			SubAttachments: subAttachments,
		})
	}
	return outs
}

func ConvertFbObjectParent(in *model.ObjectParent) *fbmessaging.FbObjectParent {
	if in == nil {
		return nil
	}
	return &fbmessaging.FbObjectParent{
		CreatedTime: in.CreatedTime.ToTime(),
		From:        ConvertObjectFrom(in.From),
		Message:     in.Message,
		ID:          in.ID,
	}
}

func ConvertFbCommentAttachment(in *model.CommentAttachment) *fbmessaging.CommentAttachment {
	if in == nil {
		return nil
	}
	var media *fbmessaging.ImageMediaDataSubAttachment
	if in.Media != nil && in.Media.Image != nil {
		media = &fbmessaging.ImageMediaDataSubAttachment{
			Image: &fbmessaging.MediaDataSubAttachment{
				Height: in.Media.Image.Height,
				Width:  in.Media.Image.Width,
				Src:    in.Media.Image.Src,
			},
		}
	}

	var target *fbmessaging.TargetDataSubAttachment
	if in.Target != nil {
		target = &fbmessaging.TargetDataSubAttachment{
			ID:  in.Target.ID,
			URL: in.Target.URL,
		}
	}
	return &fbmessaging.CommentAttachment{
		Media:  media,
		Target: target,
		Title:  in.Title,
		Type:   in.Type,
		URL:    in.URL,
	}
}

func GenerateFacebookUserPicture(userID string) string {
	return fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", userID)
}
