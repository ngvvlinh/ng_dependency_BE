package main

import (
	"flag"
	"fmt"

	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	postgres := cfg.Databases.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	mapFbExternalPost := make(map[string]*model.FbExternalPost)
	{
		var fromID dot.ID
		for {
			fbExternalPosts, err := scanFbExternalPosts(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(fbExternalPosts) == 0 {
				break
			}
			for _, fbExternalPost := range fbExternalPosts {
				mapFbExternalPost[fbExternalPost.ExternalID] = fbExternalPost
			}
			fromID = max(fromID, fbExternalPosts[len(fbExternalPosts)-1].ID)
		}
	}

	for _, fbExternalPost := range mapFbExternalPost {
		for _, externalAttachment := range fbExternalPost.ExternalAttachments {
			if externalAttachment.Type == "photo" || externalAttachment.Type == "cover_photo" {
				externalAttachment.Media = &model.MediaPostAttachment{
					Image: &model.ImageMediaPostAttachment{
						Src: fbExternalPost.ExternalPicture,
					},
				}
			}
		}
	}

	{
		count, errCount, updateCount := len(mapFbExternalPost), 0, 0
		maxGoroutines := 8
		ch := make(chan string, maxGoroutines)
		chInsert := make(chan error, maxGoroutines)
		for key, value := range mapFbExternalPost {
			ch <- key
			externalID := key
			fbExternalPost := value
			go func(_externalID string, _fbExternalPost *model.FbExternalPost) (_err error) {
				defer func() {
					<-ch
					chInsert <- _err
				}()
				_, _err = db.
					Table("fb_external_post").
					Where("external_id = ?", _externalID).
					Update(_fbExternalPost)
				return _err
			}(externalID, fbExternalPost)
		}
		for i, n := 0, len(mapFbExternalPost); i < n; i++ {
			err := <-chInsert
			if err != nil {
				errCount++
			} else {
				updateCount++
			}
		}
		ll.S.Infof("Update external_attachments for fbExternalPost: success %v/%v, error %v/%v", updateCount, count, errCount, count)
	}

	fmt.Println(len(mapFbExternalPost))
	{
		count, errCount, updateCount := len(mapFbExternalPost), 0, 0
		maxGoroutines := 8
		ch := make(chan string, maxGoroutines)
		chInsert := make(chan error, maxGoroutines)
		for key, value := range mapFbExternalPost {
			ch <- key
			externalID := key
			fbExternalPost := value
			go func(_externalID string, _fbExternalPost *model.FbExternalPost) (_err error) {
				defer func() {
					<-ch
					chInsert <- _err
				}()
				_err = db.
					Table("fb_customer_conversation").
					Where("external_id = ? and external_page_id = ?", _externalID, _fbExternalPost.ExternalPageID).
					ShouldUpdateMap(map[string]interface{}{
						"external_post_attachments": core.JSON{_fbExternalPost.ExternalAttachments},
					})
				if cm.ErrorCode(_err) == cm.NotFound {
					return nil
				}
				return _err
			}(externalID, fbExternalPost)
		}
		for i, n := 0, len(mapFbExternalPost); i < n; i++ {
			err := <-chInsert
			if err != nil {
				errCount++
			} else {
				updateCount++
			}
		}
		ll.S.Infof("Update external_post_attachments for fbCustomerConversation: success %v/%v, error %v/%v", updateCount, count, errCount, count)
	}

}

func scanFbExternalPosts(fromID dot.ID) (fbExternalPosts model.FbExternalPosts, err error) {
	err = db.
		Where(`
			id > ? and
			jsonb_array_length(external_attachments) > 0 and 
			(external_attachments @> '[{"type":"photo"}]' or
			external_attachments @> '[{"type":"cover_photo"}]') and
			external_attachments @> '[{"media":null}]' and
			external_picture is not null`, fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&fbExternalPosts)
	return
}

func max(a, b dot.ID) dot.ID {
	if a > b {
		return a
	}
	return b
}
