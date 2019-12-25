package sqlstore

import (
	"context"

	"etop.vn/backend/com/etc/logging/smslog/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
)

type SmsLogStoreFactory func(context.Context) *SmsLogStore

func NewSmsLogStore(db *cmsql.Database) SmsLogStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *SmsLogStore {
		return &SmsLogStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type SmsLogStore struct {
	query cmsql.QueryFactory
}

func (s *SmsLogStore) CreateSmsLog(sms *model.SmsLog) error {
	_, err := s.query().Insert(sms)
	return err
}
