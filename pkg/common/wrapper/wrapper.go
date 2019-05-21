package cmWrapper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap/zapcore"

	cmP "etop.vn/backend/pb/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/logline"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/etop/authorize/middleware"
)

var (
	ll  = l.New()
	bot *telebot.Channel
)

//go:generate ffjson -nodecoder $GOFILE

func InitBot(b *telebot.Channel) {
	if b != nil {
		ll.Info("Enabled sending request errors to telegram")
	}
	bot = b
}

type CensorInterface interface {
	Censor()
}

func Censor(m proto.Message) {
	if m, ok := m.(CensorInterface); ok {
		m.Censor()
	}
}

type HasErrorsInterface interface {
	HasErrors() []*cmP.Error
}

func HasErrors(m interface{}) []*cmP.Error {
	me, ok := m.(HasErrorsInterface)
	if !ok {
		return nil
	}

	errs := me.HasErrors()
	for _, err := range errs {
		if err.Code != "ok" {
			return errs
		}
	}
	return nil
}

type TwerrJSON struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

var marshaler = jsonpb.Marshaler{OrigName: true, EmitDefaults: true}

func EncodeTwirpError(w io.Writer, err twirp.Error) {
	twerr := TwerrJSON{
		Code: string(err.Code()),
		Msg:  err.Msg(),
		Meta: err.MetaMap(),
	}
	json.NewEncoder(w).Encode(twerr)
}

func SendErrorToBot(bot *telebot.Channel, rpcName string, session *middleware.Session, req interface{}, err cm.TwError, errs []*cmP.Error, d time.Duration, lvl cm.TraceLevel, stacktrace []byte) {
	if bot == nil {
		return
	}

	buf := &strings.Builder{}
	if lvl >= cm.LevelTrace {
		buf.WriteString("🔥 ")
	}
	buf.WriteString("ERROR: ")
	buf.WriteString(rpcName)
	buf.WriteString(" (")
	buf.WriteString(strconv.Itoa(int(d / time.Millisecond)))
	buf.WriteString("ms)")
	if session != nil {
		if user := session.User; user != nil {
			buf.WriteString("\n–– User: ")
			buf.WriteString(user.FullName)
			buf.WriteString(" (")
			buf.WriteString(strconv.Itoa(int(user.ID)))
			buf.WriteString(")")
		}
		if shop := session.Shop; shop != nil {
			buf.WriteString("\n–– Shop: ")
			buf.WriteString(shop.Name)
			buf.WriteString(" (")
			buf.WriteString(strconv.Itoa(int(shop.ID)))
			buf.WriteString(")")
		}
	}
	buf.WriteString("\n→")

	switch req := req.(type) {
	case proto.Message:
		marshaler.Marshal(buf, req)
	case []byte:
		if len(req) == 0 {
			buf.WriteString("<empty>")
		} else {
			if req[len(req)-1] == '\n' {
				req = req[:len(req)-1]
			}
			buf.Write(req)
		}
	case string:
		if len(req) == 0 {
			buf.WriteString("<empty>")
		} else {
			if req[len(req)-1] == '\n' {
				req = req[:len(req)-1]
			}
			buf.WriteString(req)
		}
	default:
		fmt.Fprintf(buf, "<unknown type=%T>", req)
	}

	if err != nil {
		buf.WriteString("\n")
		if cause := err.Cause(); cause != nil {
			buf.WriteString("• ")
			buf.WriteString(cause.Error())
			buf.WriteString(" ")
		}
		buf.WriteString("• ")
		buf.WriteString(cm.TrimFilePath(err.OrigFile()))
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(err.OrigLine()))
		for _, line := range err.Logs() {
			buf.WriteString("\n• ")
			buf.WriteString(line.Message)
			for _, field := range line.Fields {
				buf.WriteByte(' ')
				buf.WriteString(field.Key)
				buf.WriteByte('=')
				fmt.Fprint(buf, logline.ValueOf(field))
			}
			buf.WriteString(" • ")
			buf.WriteString(cm.TrimFilePath(line.File))
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(line.Line))
		}
		buf.WriteString("\n⇐")
		EncodeTwirpError(buf, err)

	} else {
		buf.WriteString("\n⇐")
		buf.WriteByte('[')
		for i, e := range errs {
			if i > 0 {
				buf.WriteByte(',')
			}
			marshaler.Marshal(buf, e)
		}
		buf.WriteByte(']')
	}
	if stacktrace != nil {
		buf.WriteString("\n")
		buf.Write(stacktrace)
	}

	bot.SendMessage(buf.String())
}

func RecoverAndLog(ctx context.Context, rpcName string, session *middleware.Session, req, resp proto.Message, recovered interface{}, err error, errs []*cmP.Error, t0 time.Time) (twError cm.TwError) {
	var stacktrace []byte
	if recovered != nil {
		stacktrace = debug.Stack()
		if _err, ok := recovered.(error); ok {
			err = cm.Error(cm.RuntimePanic, "", _err)
		} else {
			err = cm.Error(cm.RuntimePanic, "", errors.New(fmt.Sprint(recovered)))
		}
	}
	t1 := time.Now()
	d := t1.Sub(t0)

	if err == nil {
		if errs != nil {
			ll.Warn("->"+rpcName,
				l.Duration("d", d),
				l.Stringer("req", req),
				l.Stringer("resp", resp))
			go SendErrorToBot(bot, rpcName, session, req, nil, errs, d, cm.LevelPartialError, stacktrace)
			return nil
		}
		ll.Debug("->"+rpcName,
			l.Duration("d", d),
			l.Stringer("req", req))
		return nil
	}

	twError = cm.TwirpError(err)
	lvl := cm.GetTraceLevel(err)
	if lvl <= cm.LevelTrival {
		if cm.IsDev() {
			ll.Warn("->"+rpcName,
				l.Duration("d", d),
				l.Stringer("req", req),
				l.Error(err))
		}
		return twError
	}

	ll.Error("->"+rpcName,
		l.Duration("d", d),
		l.Stringer("req", req),
		l.Error(err))
	if cm.NotProd() || cm.ErrorCode(err) == cm.RuntimePanic || middleware.CtxDebug(ctx) != "" {
		PrintErrorWithStack(ctx, err, stacktrace)
	}
	go SendErrorToBot(bot, rpcName, session, req, twError, nil, d, lvl, stacktrace)
	return twError
}

func LogErrorAndTrace(ctx context.Context, err error, msg string, fields ...zapcore.Field) {
	if err == nil {
		return
	}

	ll.Error(msg, append(fields, l.Error(err))...)
	if cm.GetTraceLevel(err) >= cm.LevelTrace {
		fmt.Printf("%+v\n", err)
	}
}

func PrintErrorWithStack(ctx context.Context, err error, stacktrace []byte) {
	if stacktrace == nil {
		fmt.Printf("%+v\n", err) // Print err with stacktrace
	} else {
		fmt.Printf("%#v\n", err) // Print err without stacktrace
		fmt.Printf("%s", stacktrace)
	}
	bus.PrintAllStack(ctx, true)
}
