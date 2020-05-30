package aggregate

import (
	"context"

	"o.o/api/etc/logging/smslog"
	"o.o/backend/com/etc/logging/smslog/convert"
	"o.o/backend/com/etc/logging/smslog/model"
	"o.o/backend/com/etc/logging/smslog/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
)

var _ smslog.Aggregate = &SmsLogAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SmsLogAggregate struct {
	db       *cmsql.Database
	store    sqlstore.SmsLogStoreFactory
	eventBus capi.EventBus
}

func NewSmsLogAggregate(eventBus capi.EventBus, db com.LogDB) *SmsLogAggregate {
	return &SmsLogAggregate{
		db:       db,
		store:    sqlstore.NewSmsLogStore(db),
		eventBus: eventBus,
	}
}

func SmsLogAggregateMessageBus(a *SmsLogAggregate) smslog.CommandBus {
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
