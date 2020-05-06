package convert

import (
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
		outs = append(outs, &fbmessaging.FbMessageAttachment{
			ID:        in.ID,
			ImageData: imageData,
			MimeType:  in.MimeType,
			Name:      in.Name,
			Size:      in.Size,
		})
	}
	return outs
}
