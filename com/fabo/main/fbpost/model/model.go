package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type FbPost struct {
	ID                  dot.ID
	ExternalID          string
	FbPageID            dot.ID
	ParentID            dot.ID
	ExternalParentID    string
	ExternalFrom        *FbPostFrom
	ExternalPicture     string
	ExternalIcon        string
	ExternalMessage     string
	ExternalAttachments []*Attachment
	ExternalCreatedTime time.Time
	ExternalUpdatedTime time.Time
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
}

type FbPostFrom struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Attachment struct {
	MediaType      string           `json:"media_type"`
	Type           string           `json:"type"`
	SubAttachments []*SubAttachment `json:"sub_attachments"`
}

type SubAttachment struct {
	Media  *MediaDataSubAttachment  `json:"media"`
	Target *TargetDataSubAttachment `json:"target"`
	Type   string                   `json:"type"`
	URL    string                   `json:"url"`
}

type MediaDataSubAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type TargetDataSubAttachment struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
