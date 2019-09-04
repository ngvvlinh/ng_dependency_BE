package sync

import (
	"context"
	"net/url"
	"time"

	"etop.vn/api/supporting/crm/vtiger"
	"etop.vn/backend/com/supporting/crm/vtiger/convert"
	"etop.vn/backend/com/supporting/crm/vtiger/mapping"
	"etop.vn/backend/com/supporting/crm/vtiger/model"
	"etop.vn/backend/com/supporting/crm/vtiger/sqlstore"
	"etop.vn/backend/com/supporting/crm/vtiger/vtigerstore"
	cm "etop.vn/backend/pkg/common"
	sqlbuilder "etop.vn/backend/pkg/common/simple-sql-builder"
)

type SyncVtiger struct {
	vs            *vtigerstore.VtigerStore
	vtigerContact sqlstore.VtigerContactStoreFactory
}

func NewSyncVtiger(vtigerStore *vtigerstore.VtigerStore, fatctory sqlstore.VtigerContactStoreFactory) *SyncVtiger {
	return &SyncVtiger{
		vs:            vtigerStore,
		vtigerContact: fatctory,
	}
}

// sync data vtiger
func (v *SyncVtiger) SyncContact(syncTime time.Time) error {
	if _, err := v.vs.Client.GetSessionKey(); err != nil {
		return err
	}
	fileMapData := v.vs.FieldMap
	ctx := context.Background()
	page := 0
	perPage := 50
	modifiedTime := syncTime.Format(mapping.TimeLayout)
	for true {
		var b sqlbuilder.SimpleSQLBuilder

		b.Printf(`SELECT * FROM Contacts WHERE modifiedtime > ? LIMIT ?, ? ;`, modifiedTime, page*perPage, perPage)
		sqlQuery, err := b.String()
		if err != nil {
			return err
		}

		queryValues := make(url.Values)
		queryValues.Set("operation", "query")
		queryValues.Set("sessionName", v.vs.Client.SessionInfo.VtigerSession.SessionName)
		queryValues.Set("query", sqlQuery)

		result, err := v.vs.Client.RequestGet(queryValues)
		if err != nil {
			return err
		}

		if len(result.Result) == 0 {
			break
		}
		for _, value := range result.Result {
			var contact *vtiger.Contact

			contact, err = fileMapData.MapingContactVtiger2Etop(value)
			if err != nil {
				return err
			}
			modelContact := convert.ConvertModelContact(contact, v.vs.Client.SessionInfo.VtigerSession.UserID)
			err = v.createOrUpdateContactToDB(ctx, modelContact)
			if err != nil {
				return err
			}
		}
		page = page + 1
	}

	return nil
}

func (v *SyncVtiger) createOrUpdateContactToDB(ctx context.Context, contact *model.VtigerContact) error {
	_, err := v.vtigerContact(ctx).ByEtopUserID(contact.EtopUserID).GetContact()

	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		return v.vtigerContact(ctx).ByEtopUserID(contact.EtopUserID).CreateVtigerContact(contact)
	} else if err != nil {
		return err
	} else {
		return v.vtigerContact(ctx).ByEtopUserID(contact.EtopUserID).UpdateVtigerContact(contact)
	}
}
