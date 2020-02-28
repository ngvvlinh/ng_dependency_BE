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

	"go.uber.org/zap/zapcore"

	cmP "etop.vn/api/top/types/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/headers"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/capi"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
	"etop.vn/common/xerrors/logline"
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

func Censor(m interface{}) {
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

func EncodeTwirpError(w io.Writer, err xerrors.ErrorInterface) {
	twerr := TwerrJSON{
		Code: err.Code().String(),
		Msg:  err.Msg(),
		Meta: err.MetaMap(),
	}
	_ = json.NewEncoder(w).Encode(twerr)
}

func SendErrorToBot(ctx context.Context, bot *telebot.Channel, rpcName string, session *middleware.Session, req interface{}, err xerrors.TwError, errs []*cmP.Error, d time.Duration, lvl xerrors.TraceLevel, stacktrace []byte) {
	if bot == nil {
		return
	}

	buf := &strings.Builder{}
	if lvl >= xerrors.LevelTrace {
		buf.WriteString("🔥 ")
	}
	buf.WriteString("[")
	buf.WriteString(cmenv.Env().String())
	buf.WriteString(",")
	buf.WriteString(headers.GetHeader(ctx).Get("X-Forwarded-Host"))
	buf.WriteString("] @thangtran268 ERROR: ")
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
		if affiliate := session.Affiliate; affiliate != nil {
			buf.WriteString("\n–– Sale: ")
			buf.WriteString(affiliate.Name)
			buf.WriteString(" (")
			buf.WriteString(strconv.Itoa(int(affiliate.ID)))
			buf.WriteString(")")
		}
		if session.Partner != nil || session.CtxPartner != nil {
			partner := session.Partner
			if partner == nil {
				partner = session.CtxPartner
			}
			buf.WriteString("\n–– Partner: ")
			buf.WriteString(partner.Name)
			buf.WriteString(" (")
			buf.WriteString(strconv.Itoa(int(partner.ID)))
			buf.WriteString(")")
		}
	}
	buf.WriteString("\n")
	sortedHeaders := headers.GetSortedHeaders(ctx)
	for _, item := range sortedHeaders {
		for _, v := range item.Values {
			buf.WriteString(item.Key)
			buf.WriteString(": ")
			buf.WriteString(v)
			buf.WriteString("\n")
		}
	}
	buf.WriteString("→")

	switch req := req.(type) {
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
	default: // MUSTDO: interface for API message
		_ = jsonx.MarshalTo(buf, req)
	}

	if err != nil {
		buf.WriteString("\n")
		if cause := err.Cause(); cause != nil {
			buf.WriteString("• ")
			buf.WriteString(cause.Error())
			buf.WriteString(" ")
		}
		buf.WriteString("• ")
		buf.WriteString(xerrors.TrimFilePath(err.OrigFile()))
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(err.OrigLine()))
		for _, line := range err.Logs() {
			buf.WriteString("\n• ")
			buf.WriteString(line.Message)
			for _, field := range line.Fields {
				buf.WriteByte(' ')
				buf.WriteString(field.Key)
				buf.WriteByte('=')
				_, _ = fmt.Fprint(buf, logline.ValueOf(field))
			}
			buf.WriteString(" • ")
			buf.WriteString(xerrors.TrimFilePath(line.File))
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
			_ = jsonx.MarshalTo(buf, e)
		}
		buf.WriteByte(']')
	}
	if stacktrace != nil {
		buf.WriteString("\n")
		buf.Write(stacktrace)
	}

	bot.SendMessage(buf.String())
}

func RecoverAndLog(ctx context.Context, rpcName string, session *middleware.Session, req, resp capi.Message, recovered interface{}, err error, errs []*cmP.Error, t0 time.Time) (twError xerrors.TwError) {
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
	twError = xerrors.TwirpError(err)
	metrics.APIRequest(rpcName, d, twError)

	if err == nil {
		if errs != nil {
			ll.Warn("->"+rpcName,
				l.Duration("d", d),
				l.Stringer("req", req),
				l.Stringer("resp", resp))
			go SendErrorToBot(ctx, bot, rpcName, session, req, nil, errs, d, xerrors.LevelPartialError, stacktrace)
			return nil
		}
		ll.Debug("->"+rpcName,
			l.Duration("d", d),
			l.Stringer("req", req))
		return nil
	}

	lvl := xerrors.GetTraceLevel(err)
	if lvl <= xerrors.LevelTrival {
		if cmenv.IsDev() {
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
	if cmenv.NotProd() || cm.ErrorCode(err) == cm.RuntimePanic || headers.CtxDebug(ctx) != "" {
		PrintErrorWithStack(ctx, err, stacktrace)
	}
	go SendErrorToBot(ctx, bot, rpcName, session, req, twError, nil, d, lvl, stacktrace)
	return twError
}

func LogErrorAndTrace(ctx context.Context, err error, msg string, fields ...zapcore.Field) {
	if err == nil {
		return
	}

	ll.Error(msg, append(fields, l.Error(err))...)
	if xerrors.GetTraceLevel(err) >= xerrors.LevelTrace {
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
	if ll.Verbosed(6) {
		bus.PrintAllStack(ctx, true)
	} else if ll.Verbosed(3) {
		bus.PrintAllStack(ctx, false)
	} else if ll.Verbosed(1) {
		bus.PrintStack(ctx)
	}
}
