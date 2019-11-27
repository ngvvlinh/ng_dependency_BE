package aggregate

import (
	"context"

	"etop.vn/api/etc/logging/smslog"
	"etop.vn/backend/com/etc/logging/smslog/convert"
	"etop.vn/backend/com/etc/logging/smslog/model"
	"etop.vn/backend/com/etc/logging/smslog/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/capi"
)

var _ smslog.Aggregate = &SmsLogAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SmsLogAggregate struct {
	db       *cmsql.Database
	store    sqlstore.SmsLogStoreFactory
	eventBus capi.EventBus
}

func NewSmsLogAggregate(eventBus capi.EventBus, db *cmsql.Database) *SmsLogAggregate {
	return &SmsLogAggregate{
		db:       db,
		store:    sqlstore.NewSmsLogStore(db),
		eventBus: eventBus,
	}
}

func (a *SmsLogAggregate) MessageBus() smslog.CommandBus {
	b := bus.New()
	return smslog.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SmsLogAggregate) CreateSmsLog(
	ctx context.Context, args *smslog.CreateSmsArgs) (err error) {
	sms := &smslog.SmsLog{}
	if err := scheme.Convert(args, sms); err != nil {
		return err
	}
	out := &model.SmsLog{}
	if err := scheme.Convert(sms, out); err != nil {
		return err
	}

	return a.store(ctx).CreateSmsLog(out)
}
