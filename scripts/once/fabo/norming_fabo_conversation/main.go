package main

import (
	"flag"
	"fmt"

	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/com/fabo/main/fbcustomerconversationsearch/model"
	fbCustomerConvModel "o.o/backend/com/fabo/main/fbmessaging/model"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

var (
	ll    = l.New()
	cfg   config.Config
	db    *cmsql.Database
	table int
)

const (
	CustomerConversationTableName = "fb_customer_conversation"
	FbExternalCommentTableName    = "fb_external_comment"
	FbExternalMessageTableName    = "fb_external_message"

	All                  = 0
	CustomerConversation = 1
	ExternalMessage      = 2
	ExternalComment      = 3
)

func main() {
	cc.InitFlags()
	flag.IntVar(
		&table,
		"table",
		-1,
		"run migrate for table: \n - 0 (all)\n - 1 (fb_customer_conversation)\n - 2 (fb_external_message)\n - 3 (fb_external_comment)",
	)
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	postgres := cfg.Databases.Postgres
	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	switch table {
	case All:
		normCustomerConversation()
		normFbComment()
		normFbMessage()
	case CustomerConversation:
		normCustomerConversation()
	case ExternalMessage:
		normFbMessage()
	case ExternalComment:
		normFbComment()
	default:
		ll.Error("invalid table flag")
	}
}

func normCustomerConversation() {
	totalConversations, err := db.Table(CustomerConversationTableName).Count(&fbCustomerConvModel.FbCustomerConversation{})
	if err != nil {
		ll.Fatal(err.Error())
	}

	var fromID int64 = 0
	counter := 0
	for {
		conversations, err := scanCustomerConversations(fromID)
		if err != nil {
			ll.Error(fmt.Sprintf("fetch conversations got error: %s", err.Error()))
			continue
		}

		lenConversations := len(conversations)
		if lenConversations == 0 {
			break
		}

		for _, conv := range conversations {
			counter++
			fromID = conv.ID.Int64()
			err := normingConversation(conv)
			if err != nil {
				ll.Error(fmt.Sprintf("norming conversation got error: %s", err.Error()))
				continue
			}
			ll.Info(fmt.Sprintf("successsfuly norming conversation: %d/%d", counter, totalConversations))
		}
	}

	ll.Info("successful norming for fb_customer_conversation")
}

func normFbComment() {
	totalComments, err := db.Table(FbExternalCommentTableName).Count(&fbCustomerConvModel.FbExternalComment{})
	if err != nil {
		ll.Fatal(err.Error())
	}

	var fromID int64 = 0
	counter := 0

	for {
		comments, err := scanFbExternalComment(fromID)
		if err != nil {
			ll.Error(fmt.Sprintf("fetch comments got error: %s", err.Error()))
			continue
		}

		lenComments := len(comments)
		if lenComments == 0 {
			break
		}

		for _, cmt := range comments {
			fromID = cmt.ID.Int64()
			counter++
			err := normingFbExternalComment(cmt)
			if err != nil {
				ll.Error(fmt.Sprintf("norming comment got error: %s", err.Error()))
				continue
			}
			ll.Info(fmt.Sprintf("successsfuly norming comments: %d/%d", counter, totalComments))
		}
	}
	ll.Info("successful norming for fb_external_comment")
}

func normFbMessage() {
	totalMessages, err := db.Table(FbExternalMessageTableName).Count(&fbCustomerConvModel.FbExternalMessage{})
	if err != nil {
		ll.Fatal(err.Error())
	}

	var fromID int64 = 0
	counter := 0

	for {
		messages, err := scanFbExternalMessage(fromID)
		if err != nil {
			ll.Error(fmt.Sprintf("fetch comments got error: %s", err.Error()))
			continue
		}

		lenMessages := len(messages)
		if lenMessages == 0 {
			break
		}

		for _, msg := range messages {
			fromID = msg.ID.Int64()
			counter++
			err := normingFbExternalMessage(msg)
			if err != nil {
				ll.Error(fmt.Sprintf("norming message got error: %s", err.Error()))
				continue
			}
			ll.Info(fmt.Sprintf("successsfuly norming messages: %d/%d", counter, totalMessages))
		}
	}
	ll.Info("successful norming for fb_external_message")
}

func normalizeText(s string) string {
	return validate.NormalizedSearchToTsVector(validate.NormalizeSearch(s))
}

func scanCustomerConversations(from int64) (fbCustomerConvModel.FbCustomerConversations, error) {
	var conversations fbCustomerConvModel.FbCustomerConversations
	err := db.
		Where("id > ?", from).
		OrderBy("id").
		Limit(1000).
		Find(&conversations)
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func normingConversation(conv *fbCustomerConvModel.FbCustomerConversation) error {
	normModel := &model.FbCustomerConversationSearch{
		ID:                   conv.ID,
		ExternalUserNameNorm: normalizeText(conv.ExternalUserName),
		CreatedAt:            conv.CreatedAt,
		ExternalPageID:       conv.ExternalPageID,
	}
	return db.ShouldInsert(normModel)
}

func scanFbExternalComment(from int64) (fbCustomerConvModel.FbExternalComments, error) {
	var externalComments fbCustomerConvModel.FbExternalComments
	err := db.
		Where("id > ?", from).
		OrderBy("id").
		Limit(1000).
		Find(&externalComments)
	if err != nil {
		return nil, err
	}
	return externalComments, nil
}

func normingFbExternalComment(cmt *fbCustomerConvModel.FbExternalComment) error {
	externalUserID := cmt.ExternalUserID
	if externalUserID == cmt.ExternalPageID {
		externalUserID = cmt.ExternalParentUserID
	}
	normModel := &model.FbExternalCommentSearch{
		ID:                  cmt.ID,
		ExternalMessageNorm: normalizeText(cmt.ExternalMessage),
		ExternalPostID:      cmt.ExternalPostID,
		ExternalUserID:      externalUserID,
		ExternalPageID:      cmt.ExternalPageID,
		CreatedAt:           cmt.CreatedAt,
	}
	return db.ShouldInsert(normModel)
}

func scanFbExternalMessage(from int64) (fbCustomerConvModel.FbExternalMessages, error) {
	var externalMessages fbCustomerConvModel.FbExternalMessages
	err := db.
		Where("id > ?", from).
		OrderBy("id").
		Limit(1000).
		Find(&externalMessages)
	if err != nil {
		return nil, err
	}
	return externalMessages, nil
}

func normingFbExternalMessage(msg *fbCustomerConvModel.FbExternalMessage) error {
	normModel := &model.FbExternalMessageSearch{
		ID:                     msg.ID,
		ExternalMessageNorm:    normalizeText(msg.ExternalMessage),
		ExternalPageID:         msg.ExternalPageID,
		CreatedAt:              msg.CreatedAt,
		ExternalConversationID: msg.ExternalConversationID,
	}
	return db.ShouldInsert(normModel)
}
