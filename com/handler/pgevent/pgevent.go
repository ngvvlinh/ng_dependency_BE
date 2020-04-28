package pgevent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

const PgChannel = "pgrid"

type TGOP string

const (
	OpInsert = "INSERT"
	OpUpdate = "UPDATE"
	OpDelete = "DELETE"
)

func (op TGOP) String() string { return string(op) }

func (op TGOP) ToEventName(name string) (string, error) {
	switch op {
	case OpInsert:
		return name + "_created", nil
	case OpUpdate:
		return name + "_updated", nil
	case OpDelete:
		return name + "_deleted", nil
	default:
		return "", fmt.Errorf("pgevent: Invalid op (%v)", string(op))
	}
}

type PgEvent struct {
	EventKey string `json:"-"`

	ID    int64  `json:"id"`
	RID   int64  `json:"rid"`
	Table string `json:"table"`
	Op    TGOP   `json:"op"`

	Keys map[string]int64 `json:"keys"`

	Timestamp int64 `json:"t"`
}

type Service struct {
	pglistener *pq.Listener
	producer   *mq.KafkaProducer
	prefix     string
	dbName     model.DBName
}

func NewService(ctx context.Context, dbName model.DBName, pgcfg cmsql.ConfigPostgres, producer *mq.KafkaProducer, prefix string) (Service, error) {
	listener := cmsql.NewListener(pgcfg, 10*time.Millisecond, 120*time.Second, cmsql.DefaultListenerProblemReport)
	if err := cmsql.ListenTo(ctx, listener, PgChannel); err != nil {
		return Service{}, err
	}

	return Service{
		pglistener: listener,
		producer:   producer,
		prefix:     prefix + "_pgrid_",
		dbName:     dbName,
	}, nil
}

func (s Service) StartForwarding(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case noti := <-s.pglistener.Notify:
			if noti != nil {
				s.HandleNotification(noti)
			}
		}
	}
}

func StartForwardings(ctx context.Context, ss []Service) {
	for _, s := range ss {
		go s.StartForwarding(ctx)
	}
}

func (s Service) HandleNotification(noti *pq.Notification) {
	ll.Info("HandleNotification :: ", l.Object("noti", noti))
	err := s.HandleNotificationWithError(noti)
	if err != nil {
		ll.Error("HandleNotification", l.Error(err), l.Object("notification", noti))
	}
}

func (s Service) HandleNotificationWithError(noti *pq.Notification) error {
	if noti.Channel != PgChannel {
		return fmt.Errorf("unknown channel: %v", noti.Channel)
	}

	event, err := ParseEventPayload(noti.Extra)
	if err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	data, err := jsonx.Marshal(event)
	if err != nil {
		panic(err)
	}

	if s.prefix == "" {
		panic("empty prefix")
	}
	topic := s.prefix + event.Table
	d, ok := TopicMap[event.Table]
	if !ok {
		return fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	if d.DBName != s.dbName {
		return fmt.Errorf("The topic is in a different database, table: %v, database: %v, expected database: %v", d.Name, d.DBName, s.dbName)
	}

	partition := int(event.ID % int64(d.Partitions)) // TODO: composition primary key

	ll.Info("HandleNotificationWithError :: ", l.String("topic", topic), l.Object("topic", d), l.Int("partition", partition))
	s.producer.Send(topic, partition, event.EventKey, data)
	return nil
}

// table:rid:op:id
func ParseEventPayload(data string) (*PgEvent, error) {
	parts := strings.Split(data, ":")
	if len(parts) < 4 {
		return nil, errors.New("Must be format table:rid:op:id or table:rid:op:k123:...")
	}

	rid, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	table := parts[0]
	op := TGOP(parts[2])
	event := &PgEvent{
		EventKey:  parts[0] + ":" + parts[3],
		Op:        op,
		RID:       rid,
		Table:     table,
		Timestamp: time.Now().Unix(),
	}

	if parts[3] == "" {
		return nil, errors.New("id is empty")
	}
	if ch := parts[3][0]; ch >= '0' && ch <= '9' {
		if len(parts) != 4 {
			return nil, errors.New("Must be format table:rid:op:id (len=4)")
		}

		id, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return nil, err
		}
		event.ID = id

	} else {
		// Parse the remaining as keyed ids
		keys := make(map[string]int64)
		for i, p := range parts[3:] {
			key, id, err := parseKeyedID(p)
			if err != nil {
				return nil, err
			}
			keys[key] = id

			// Store the first part as id
			if i == 0 {
				event.ID = id
			}
		}
		event.Keys = keys
	}
	return event, nil
}

func parseKeyedID(s string) (key string, id int64, err error) {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch >= '0' && ch <= '9' {
			key = s[:i]
			break
		}
	}
	if key == "" {
		return "", 0, fmt.Errorf("Empty key (%v)", s)
	}

	idStr := s[len(key):]
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("%v (%v)", err, s)
	}
	return key, id, nil
}
