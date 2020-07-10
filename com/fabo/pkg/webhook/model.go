package webhook

import fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"

type WebhookMessageType string

const (
	WebhookFeed           WebhookMessageType = "feed"
	WebhookMessage        WebhookMessageType = "message"
	WebhookInvalidMessage WebhookMessageType = "invalid"
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

func (msg *WebhookMessages) IsOwnerPageComment() bool {
	if msg.MessageType() == WebhookFeed {
		return false
	}

	for _, entry := range msg.Entry {
		pageId := entry.ID
		for _, change := range entry.Changes {
			if change.Value.From.ID == pageId {
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
	StickerID   string               `json:"sticher_id"`
	Attachments []*MessageAttachment `json:"attachments"`
}

type MessageAttachment struct {
	Title   string                 `json:"title"`
	Url     string                 `json:"url"`
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
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

func (v FeedChange) IsComment() bool {
	return v.Value.Item == "comment"
}

func (v FeedChange) IsAdminPost(externalPageID string) bool {
	return !v.IsComment() && v.Value.Item != "reaction" && v.Value.From.ID == externalPageID
}

func (v FeedChange) IsCreated() bool {
	return v.Value.Verb == FeedAdd
}

func (v FeedChange) IsEdited() bool {
	return v.Value.Verb == FeedEdited
}

func (v FeedChange) IsOnChildPost() bool {
	return v.Value.Item == "photo"
}

func (v FeedChange) IsOnParentPost() bool {
	return v.Value.Item == "status"
}

type FeedVerb string

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
