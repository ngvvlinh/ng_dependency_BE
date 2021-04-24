package webhook

import (
	"fmt"
	"strings"

	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
)

type WebhookMessageType string
type MessageAttachmentType string

type WebhookUserType string
type LiveVideoStatus string

const (
	WebhookFeed           WebhookMessageType = "feed"
	WebhookMessage        WebhookMessageType = "message"
	WebhookInvalidMessage WebhookMessageType = "invalid"

	MessageAttachmentImage    MessageAttachmentType = "image"
	MessageAttachmentFallback MessageAttachmentType = "fallback"
	MessageAttachmentVideo    MessageAttachmentType = "video"

	FeedComment  = "comment"
	FeedEvent    = "event"
	FeedReaction = "reaction"
	FeedPhoto    = "photo"
	FeedStatus   = "status"

	EventPermalinkPrefix = "https://www.facebook.com/events/"

	WebhookUserFeed       WebhookUserType = "feed"
	WebhookUserLiveVideos WebhookUserType = "live_videos"

	UnpublishedStatus LiveVideoStatus = "unpublished"
	LiveStatus        LiveVideoStatus = "live"
	LiveStoppedStatus LiveVideoStatus = "live_stopped"
	VodStatus         LiveVideoStatus = "vod"

	separator = "-"
)

// Model for Message
type WebhookMessages struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

func (msg *WebhookMessages) MessageType() WebhookMessageType {
	if len(msg.Entry) > 0 {
		for _, val := range msg.Entry {
			if val.Changes != nil && len(val.Changes) > 0 {
				return WebhookFeed
			}

			if val.Messaging != nil && len(val.Messaging) > 0 {
				return WebhookMessage
			}
		}
	}
	return WebhookInvalidMessage
}

func (msg *WebhookMessages) GetKey() string {

	switch msg.MessageType() {
	case WebhookMessage:
		key := msg.Object
		if len(msg.Entry) == 0 {
			return key
		}
		entry := msg.Entry[0]
		key += separator + entry.ID

		if len(entry.Messaging) == 0 {
			return key
		}

		messaging := entry.Messaging[0]
		if messaging.Recipient != nil && messaging.Recipient.ID != entry.ID {
			key += separator + messaging.Recipient.ID
		}

		if messaging.Sender != nil && messaging.Sender.ID != entry.ID {
			key += separator + messaging.Sender.ID
		}

		return key
	case WebhookFeed:
		key := msg.Object
		if len(msg.Entry) == 0 {
			return key
		}
		entry := msg.Entry[0]
		key += separator + entry.ID

		if len(entry.Changes) == 0 {
			return key
		}

		feedChange := entry.Changes[0]
		key += separator + feedChange.Field
		changeValue := feedChange.Value

		key += separator + changeValue.PostID

		if changeValue.CommentID == "" {
			return key
		}

		return separator + changeValue.From.ID
	default:
		return ""
	}
}

func (msg *WebhookMessages) IsCreateOrEditCommentFromPageOwner() bool {
	for _, entry := range msg.Entry {
		pageId := entry.ID
		for _, change := range entry.Changes {
			if change.IsRemove() {
				return false
			}

			// like/ unlike or hide/ unhide
			if change.IsPageUnLikeComment(pageId) || change.IsPageLikeComment(pageId) ||
				change.IsPageHideComment(pageId) || change.IsPageUnHideComment(pageId) {
				return false
			}

			if change.Value.From.ID == pageId && change.IsComment() {
				return true
			}
		}
	}
	return false
}

type Entry struct {
	ID        string                     `json:"id"`
	Time      fbclientmodel.FacebookTime `json:"time"`
	Messaging []*Messaging               `json:"messaging"`
	Changes   []FeedChange               `json:"changes"`
}

func (e *Entry) ExternalID() string {
	switch {
	case len(e.Changes) > 0:
		return e.Changes[0].Value.PostID
	case len(e.Messaging) > 0:
		return e.Messaging[0].Message.Mid
	default:
		return ""
	}
}

type Messaging struct {
	Sender    *Sender                    `json:"sender"`
	Recipient *Recipient                 `json:"recipient"`
	Timestamp fbclientmodel.FacebookTime `json:"timestamp"`
	Message   *Message                   `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Mid         string               `json:"mid"`
	Text        string               `json:"text"`
	StickerID   int                  `json:"sticker_id"`
	Attachments []*MessageAttachment `json:"attachments"`
}

func (m *Message) ConvertToMessageData(fromID, toID string, createdTime fbclientmodel.FacebookTime) *fbclientmodel.MessageData {
	if m == nil {
		return nil
	}

	var (
		shares      []*fbclientmodel.MessageDataShare
		attachments []*fbclientmodel.MessageDataAttachment
		sticker     string
	)

	for _, attachment := range m.Attachments {
		attachmentPayload := attachment.Payload
		if attachmentPayload == nil {
			continue
		}

		// type == 'image'
		if attachment.Type == string(MessageAttachmentImage) {
			// message is a sticker
			if attachmentPayload.StickerID != 0 {
				shares = append(shares, &fbclientmodel.MessageDataShare{
					ID:   fmt.Sprintf("%d", attachmentPayload.StickerID),
					Link: attachmentPayload.Url,
				})
				sticker = attachmentPayload.Url
				continue
			}

			attachments = append(attachments, &fbclientmodel.MessageDataAttachment{
				ImageData: &fbclientmodel.MessageDataAttachmentImage{
					URL:        attachmentPayload.Url,
					PreviewURL: attachmentPayload.Url,
				},
			})

		}

		// type == 'fallback'
		// message is fallback (a post is shared or link is sent into messenger)
		if attachment.Type == string(MessageAttachmentFallback) {
			shares = append(shares, &fbclientmodel.MessageDataShare{
				ID:   m.Mid,
				Link: attachmentPayload.Url,
				Name: attachmentPayload.Title,
			})
		}

		// type == 'video'
		if attachment.Type == string(MessageAttachmentVideo) {
			attachments = append(attachments, &fbclientmodel.MessageDataAttachment{
				VideoData: &fbclientmodel.MessageDataAttachmentVideoData{
					URL:        attachmentPayload.Url,
					PreviewURL: attachmentPayload.Url,
				},
			})
		}
	}

	// handle case when user send link or share a post into messenger
	{
		var strs []string
		if m.Text != "" {
			strs = append(strs, m.Text)
		}
		// Get first share when sticker is empty
		if len(shares) > 0 && sticker == "" {
			if shares[0].Description != "" {
				strs = append(strs, shares[0].Description)
			}
			if shares[0].Link != "" {
				strs = append(strs, shares[0].Link)
			} else {
				strs = append(strs, shares[0].Name)
			}
		}
		m.Text = strings.Join(strs, "\n")
	}

	messageData := &fbclientmodel.MessageData{
		ID:          m.Mid,
		CreatedTime: &createdTime,
		Message:     m.Text,
		To: &fbclientmodel.ObjectsTo{
			Data: []*fbclientmodel.ObjectTo{
				{
					ID: toID,
				},
			},
		},
		From: &fbclientmodel.ObjectFrom{
			ID: fromID,
		},
	}
	if len(shares) != 0 {
		messageData.Shares = &fbclientmodel.MessageDataShares{
			Data: shares,
		}
	}
	if len(attachments) != 0 {
		messageData.Attachments = &fbclientmodel.MessageDataAttachments{
			Data: attachments,
		}
	}
	messageData.Sticker = sticker

	return messageData
}

type MessageAttachment struct {
	Type    string                    `json:"type"`
	Payload *MessageAttachmentPayload `json:"payload"`
}

type MessageAttachmentPayload struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	StickerID int64  `json:"sticker_id"`
}

type FeedEntry struct {
	ID      string                     `json:"id"`
	Time    fbclientmodel.FacebookTime `json:"time"`
	Changes []FeedChange               `json:"changes"`
}

type FeedChange struct {
	Field string      `json:"field"`
	Value ChangeValue `json:"value"`
}

func (v FeedChange) IsEventComment() bool {
	return strings.HasPrefix(v.Value.Post.PermalinkUrl, EventPermalinkPrefix)
}

func (v FeedChange) IsComment() bool {
	return v.Value.Item == FeedComment || v.Value.CommentID != ""
}

func (v FeedChange) IsEvent() bool {
	return v.Value.Item == FeedEvent
}

func (v FeedChange) IsAdminPost(externalPageID string) bool {
	return !v.IsComment() && v.Value.Item != FeedReaction && v.Value.From.ID == externalPageID
}

func (v FeedChange) IsCreated() bool {
	return v.Value.Verb == FeedAdd
}

func (v FeedChange) IsEdited() bool {
	return v.Value.Verb == FeedEdited
}

func (v FeedChange) IsOnChildPost() bool {
	return v.Value.Item == FeedPhoto
}

func (v FeedChange) IsOnParentPost() bool {
	return v.Value.Item == FeedStatus
}

func (v FeedChange) IsRemove() bool {
	return v.Value.Verb == FeedRemove && v.Value.ReactionType == ""
}

func (v FeedChange) IsDelete() bool {
	return v.Value.Verb == FeedDelete
}

// page like user's comment
func (v FeedChange) IsPageLikeComment(extPageID string) bool {
	if !v.IsComment() {
		return false
	}

	return v.Value.From.ID == extPageID && v.Value.ReactionType == string(ReactionLike) && v.Value.Verb == FeedAdd
}

// page unlike user's comment
func (v FeedChange) IsPageUnLikeComment(extPageID string) bool {
	if !v.IsComment() {
		return false
	}

	return v.Value.From.ID == extPageID && v.Value.ReactionType == string(ReactionLike) && v.Value.Verb == FeedRemove
}

// page hide user's comment
func (v FeedChange) IsPageHideComment(extPageID string) bool {
	if !v.IsComment() {
		return false
	}

	return v.Value.From.ID == extPageID && v.Value.Verb == FeedHide
}

// page unhide user's comment
func (v FeedChange) IsPageUnHideComment(extPageID string) bool {
	if !v.IsComment() {
		return false
	}

	return v.Value.From.ID == extPageID && v.Value.Verb == FeedUnhide
}

type FeedVerb string
type CommentVerb string
type CommentReaction string

// add, block, edit, edited, delete, follow, hide, mute, remove, unblock, unhide, update
const (
	FeedAdd     FeedVerb = "add"
	FeedBlock   FeedVerb = "block"
	FeedEdit    FeedVerb = "edit"
	FeedEdited  FeedVerb = "edited"
	FeedDelete  FeedVerb = "delete"
	FeedFollow  FeedVerb = "follow"
	FeedHide    FeedVerb = "hide"
	FeedMute    FeedVerb = "mute"
	FeedRemove  FeedVerb = "remove"
	FeedUnblock FeedVerb = "unblock"
	FeedUnhide  FeedVerb = "unhide"
	FeedUpdate  FeedVerb = "update"

	ReactionLike CommentReaction = "like"
)

type ChangeValue struct {
	EditedTime fbclientmodel.FacebookTime `json:"edited_time"`

	From FeedFrom `json:"from"`
	Post FeedPost `json:"post"`

	IsHidden bool   `json:"is_hidden"`
	Link     string `json:"link"`
	Message  string `json:"message"`

	PostID string `json:"post_id"`

	PhotoID string   `json:"photo_id"`
	Photo   string   `json:"photo"`
	Photos  []string `json:"photos"`

	CommentID string `json:"comment_id"`

	CreatedTime  fbclientmodel.FacebookTime `json:"created_time"`
	Item         string                     `json:"item"`
	ParentID     string                     `json:"parent_id"`
	ReactionType string                     `json:"reaction_type"`
	Verb         FeedVerb                   `json:"verb"`
}

type FeedFrom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FeedPost struct {
	ID              string `json:"id"`
	StatusType      string `json:"status_type"`
	IsPublished     bool   `json:"is_published"`
	UpdatedTime     string `json:"updated_time"`
	PermalinkUrl    string `json:"permalink_url"`
	PromotionStatus string `json:"promotion_status"`
}

// FbUser:
// {
//		"object": "user",
//		"entry": [
//			{
//				"id": "1616082448555860",
//				"uid": "1616082448555860",
//				"time": 1619265393,
//				"changes": [
//					{
//						"field": "feed"
//					}
//				]
//			}
//		]
// }

type WebhookUser struct {
	Object string       `json:"object"`
	Entry  []*UserEntry `json:"entry"`
}

func (m *WebhookUser) Type() WebhookUserType {
	if m == nil || len(m.Entry) == 0 || len(m.Entry[0].Changes) == 0 {
		return ""
	}
	return m.Entry[0].Changes[0].Field
}

func (m *WebhookUser) GetKey() string {
	switch m.Type() {
	case WebhookUserLiveVideos:
		return string(m.Type()) + separator + m.Entry[0].ID
	default:
		return string(m.Type())
	}
}

type UserEntry struct {
	ID      string                     `json:"id"`
	UID     string                     `json:"uid"`
	Time    fbclientmodel.FacebookTime `json:"time"`
	Changes []*UserEntryChange         `json:"changes"`
}

type UserEntryChange struct {
	Field WebhookUserType       `json:"field"`
	Value *UserEntryChangeValue `json:"value"`
}

type UserEntryChangeValue struct {
	ID     string          `json:"id"`
	Status LiveVideoStatus `json:"status"`
}

type LiveVideoComment struct {
	ID      string                    `json:"id"`
	Message string                    `json:"message"`
	From    *fbclientmodel.ObjectFrom `json:"from"`
}
