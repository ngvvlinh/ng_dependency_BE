package convert

import (
	"time"

	"o.o/api/fabo/fbmessaging"
)

// +gen:convert: o.o/backend/com/fabo/main/fbmessaging/model -> o.o/api/fabo/fbmessaging
// +gen:convert: o.o/api/fabo/fbmessaging

func ConvertCreateFbExternalCommentArgsToFbExternalComment(in *fbmessaging.CreateFbExternalCommentArgs, out *fbmessaging.FbExternalComment) *fbmessaging.FbExternalComment {
	if in == nil {
		return nil
	}
	if out == nil {
		return &fbmessaging.FbExternalComment{}
	}
	apply_fbmessaging_CreateFbExternalCommentArgs_fbmessaging_FbExternalComment(in, out)
	out.CreatedAt = time.Now()
	out.UpdatedAt = time.Now()
	return out
}
