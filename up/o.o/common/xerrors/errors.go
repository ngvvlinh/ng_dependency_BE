package xerrors

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"o.o/common/jsonx"
	"o.o/common/xerrors/logline"
)

type TraceLevel int

const (
	LevelNoError TraceLevel = iota
	LevelTrival
	LevelPartialError
	LevelError
	LevelTrace
	LevelInternal
	LevelPanic
)

func (lvl TraceLevel) String() string {
	switch lvl {
	case LevelNoError:
		return "NoError"
	case LevelTrival:
		return "Trival"
	case LevelPartialError:
		return "PartialError"
	case LevelError:
		return "Error"
	case LevelTrace:
		return "Trace"
	case LevelInternal:
		return "Internal"
	case LevelPanic:
		return "Panic"
	default:
		return "TraceLevel(" + strconv.Itoa(int(lvl)) + ")"
	}
}

// Code ...
type Code int

// Error constants from twirp
const (
	NoError            = Code(0)
	Canceled           = Code(1)
	Unknown            = Code(2)
	InvalidArgument    = Code(3)
	DeadlineExceeded   = Code(4)
	NotFound           = Code(5)
	AlreadyExists      = Code(6)
	PermissionDenied   = Code(7)
	ResourceExhausted  = Code(8)
	FailedPrecondition = Code(9)
	Aborted            = Code(10)
	OutOfRange         = Code(11)
	Unimplemented      = Code(12)
	Internal           = Code(13)
	Unavailable        = Code(14)
	DataLoss           = Code(15)
	Unauthenticated    = Code(16)

	RuntimePanic = Code(100)
	BadRoute     = Code(101)
	Malformed    = Code(102)

	WrongPassword        = Code(1005)
	ValidationFailed     = Code(1501)
	STokenRequired       = Code(1607)
	CaptchaRequired      = Code(1610)
	CaptchaInvalid       = Code(1611)
	RegisterRequired     = Code(1702)
	ExternalServiceError = Code(1909)
	AccountClosed        = Code(2001)

	SkipSync = Code(2101)

	FacebookPermissionDenied = Code(3001)
	Facebook                 = Code(3005)
)

var (
	ErrTODO    = Error(Unimplemented, "TODO", nil)
	ErrREMOVED = Error(Unavailable, "The function is no longer available", nil)

	ErrUnauthenticated  = Error(Unauthenticated, "", nil).MarkTrivial()
	ErrPermissionDenied = Error(PermissionDenied, "", nil).MarkTrivial()
)

type CustomCode struct {
	StdCode        Code
	String         string
	DefaultMessage string
}

var (
	mapCodes       [Unauthenticated + 1]string
	mapCustomCodes map[Code]*CustomCode
)

func init() {
	mapCodes[Canceled] = "canceled"
	mapCodes[Unknown] = "unknown"
	mapCodes[InvalidArgument] = "invalid_argument"
	mapCodes[DeadlineExceeded] = "deadline_exceeded"
	mapCodes[NotFound] = "not_found"
	mapCodes[AlreadyExists] = "already_exists"
	mapCodes[PermissionDenied] = "permission_denied"
	mapCodes[Unauthenticated] = "unauthenticated"
	mapCodes[ResourceExhausted] = "resource_exhausted"
	mapCodes[FailedPrecondition] = "failed_precondition"
	mapCodes[Aborted] = "aborted"
	mapCodes[OutOfRange] = "out_of_range"
	mapCodes[Unimplemented] = "unimplemented"
	mapCodes[Internal] = "internal"
	mapCodes[Unavailable] = "unavailable"
	mapCodes[DataLoss] = "data_loss"
	mapCodes[NoError] = "ok"

	mapCustomCodes = make(map[Code]*CustomCode)
	mapCustomCodes[RuntimePanic] = &CustomCode{Internal, "runtime", ""}
	mapCustomCodes[BadRoute] = &CustomCode{NotFound, "bad_route", ""}
	mapCustomCodes[Malformed] = &CustomCode{InvalidArgument, "malformed", ""}
	mapCustomCodes[ValidationFailed] = &CustomCode{InvalidArgument, "validation_failed", "Dữ liệu không hợp lệ."}
	mapCustomCodes[ExternalServiceError] = &CustomCode{Unknown, "external_service", "Đã xảy ra lỗi khi kết nối với hệ thống bên ngoài"}
	mapCustomCodes[WrongPassword] = &CustomCode{Unauthenticated, "wrong_password", "Mật khẩu không hợp lệ"}
	mapCustomCodes[SkipSync] = &CustomCode{FailedPrecondition, "skip_sync", "skip_sync"}
	mapCustomCodes[STokenRequired] = &CustomCode{PermissionDenied, "stoken_required", "stoken_required"}
	mapCustomCodes[CaptchaRequired] = &CustomCode{Unauthenticated, "captcha_required", "Captcha is required"}
	mapCustomCodes[CaptchaInvalid] = &CustomCode{Unauthenticated, "invalid_captcha", "Mã xác thực không hợp lệ"}
	mapCustomCodes[RegisterRequired] = &CustomCode{FailedPrecondition, "register_required", "register_required"}
	mapCustomCodes[AccountClosed] = &CustomCode{Unavailable, "account_closed", "Tài khoản không còn được sử dụng hoặc đã bị xoá"}
	mapCustomCodes[FacebookPermissionDenied] = &CustomCode{PermissionDenied, "facebook_permission_denied", ""}
	mapCustomCodes[Facebook] = &CustomCode{FailedPrecondition, "facebook", ""}
}

func (c Code) String() string {
	if IsValidStandardErrorCode(c) {
		return mapCodes[c]
	}
	if s := mapCustomCodes[c]; s != nil {
		return s.String
	}
	return "Code(" + strconv.Itoa(int(c)) + ")"
}

func DefaultErrorMessage(code Code) string {
	switch code {
	case NoError:
		return ""
	case NotFound:
		return "Không tìm thấy."
	case InvalidArgument:
		return "Có lỗi xảy ra.."
	case Internal:
		return "Lỗi không xác định."
	case Unauthenticated:
		return "Vui lòng đăng nhập (hoặc đăng ký nếu chưa có tài khoản)."
	case PermissionDenied:
		return "Không tìm thấy hoặc cần quyền truy cập.."
	case Unimplemented:
		return "TODO"
	}
	if s := mapCustomCodes[code]; s != nil && s.DefaultMessage != "" {
		return s.DefaultMessage
	}
	return "Lỗi không xác định."
}

func IsValidStandardErrorCode(c Code) bool {
	return c >= 0 && int(c) < len(mapCodes)
}

func GetCustomCode(c Code) *CustomCode {
	return mapCustomCodes[c]
}

func IsValidErrorCode(c Code) bool {
	return IsValidStandardErrorCode(c) || mapCustomCodes[c] != nil
}

// IError defines error interface returned by errors package
type IError interface {
	error
	IStack

	Format(st fmt.State, verb rune)
}

type IStack interface {
	StackTrace() errors.StackTrace
}

// APIError ...
type APIError struct {
	Code     Code
	XCode    Code
	Err      error
	Message  string
	Original string
	OrigFile string
	OrigLine int
	Stack    errors.StackTrace
	Trace    bool
	Trivial  bool
	Logs     []*logline.LogLine
	Meta     map[string]string
}

func ToError(err error) *APIError {
	if err, ok := err.(*APIError); ok {
		return err
	}
	return newError(true, true, Internal, "", err)
}

func Error(code Code, message string, err error) *APIError {
	return newError(false, true, code, message, err)
}

func ErrorTrace(code Code, message string, err error) *APIError {
	return newError(true, true, code, message, err)
}

// NoStackError
func NSErrorf(code Code, err error, message string, args ...interface{}) *APIError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return newError(false, false, code, message, err)
}

func Errorf(code Code, err error, message string, args ...interface{}) *APIError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return newError(false, true, code, message, err)
}

func ErrorTracef(code Code, err error, message string, args ...interface{}) *APIError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return newError(true, true, code, message, err)
}

func Trace(err error) *APIError {
	if xerr, ok := err.(*APIError); ok {
		xerr.Trace = true
		return xerr
	}
	if err != nil {
		return newError(true, true, Internal, err.Error(), err)
	}
	return newError(true, true, Internal, "Expected error!", nil)
}

func GetMeta(err error) map[string]string {
	if xerr, ok := err.(*APIError); ok {
		return xerr.Meta
	}
	return nil
}

func MarkTrivial(err error, code Code) error {
	if xerr, ok := err.(*APIError); ok &&
		(xerr.Code == code || xerr.XCode == code) {
		xerr.Trivial = true
	}
	return err
}

func MarkTrivials(err error, codes ...Code) error {
	if xerr, ok := err.(*APIError); ok {
		for _, code := range codes {
			if xerr.Code == code || xerr.XCode == code {
				xerr.Trivial = true
				return err
			}
		}
	}
	return err
}

func IsTrivial(err error) bool {
	if xerr, ok := err.(*APIError); ok {
		return xerr.Trivial
	}
	return false
}

func newError(trace bool, stack bool, code Code, message string, err error) *APIError {
	if message == "" {
		message = DefaultErrorMessage(code)
	}

	var xcode Code
	if !IsValidStandardErrorCode(code) {
		xcode = code
		if s := GetCustomCode(code); s != nil {
			code = s.StdCode
		} else {
			code = Unknown
		}
	}
	if err != nil {
		// Overwrite *Error
		if xerr, ok := err.(*APIError); ok {
			// Keep original message
			if xerr.Original == "" {
				xerr.Original = xerr.Message
			}
			xerr.Code = code
			xerr.XCode = xcode
			xerr.Message = message
			xerr.Trace = xerr.Trace || trace
			return xerr
		}
	}

	// Always include the original location
	_, file, line, _ := runtime.Caller(2)
	xerr := &APIError{
		Err:      err,
		Code:     code,
		XCode:    xcode,
		Message:  message,
		Original: "",
		OrigFile: file,
		OrigLine: line,
		Trace:    trace,
	}

	// Wrap error with stacktrace
	if stack {
		xerr.Stack = errors.New("").(IStack).StackTrace()[2:]
	}
	return xerr
}

func (e *APIError) Log(msg string, fields ...zapcore.Field) *APIError {
	_, file, line, _ := runtime.Caller(1)
	e.Logs = append(e.Logs, &logline.LogLine{
		Message: msg,
		File:    file,
		Line:    line,
		Fields:  fields,
	})
	return e
}

func (e *APIError) Logf(format string, args ...interface{}) *APIError {
	_, file, line, _ := runtime.Caller(1)
	e.Logs = append(e.Logs, &logline.LogLine{
		Message: fmt.Sprintf(format, args...),
		File:    file,
		Line:    line,
	})
	return e
}

// Error ...
func (e *APIError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	if e.Err != nil {
		b.WriteString(" cause=")
		b.WriteString(e.Err.Error())
	}
	if e.Original != "" {
		b.WriteString(" original=")
		b.WriteString(e.Original)
	}
	for k, v := range e.Meta {
		b.WriteByte(' ')
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
	}
	return b.String()
}

// MarkTrivial ...
func (e *APIError) MarkTrivial() *APIError {
	e.Trivial = true
	return e
}

// Cause ...
func (e *APIError) Cause() error {
	return e.Err
}

// StackTrace ...
func (e *APIError) StackTrace() errors.StackTrace {
	return e.Stack
}

// Format ...
func (e *APIError) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('#') || st.Flag('+'):
			_, _ = fmt.Fprintf(st, "\ncode=%v message=%v", e.Code, e.Message)
			lvl := GetTraceLevel(e)
			_, _ = fmt.Fprint(st, " level=", lvl.String())
			for k, v := range e.Meta {
				_, _ = fmt.Fprint(st, " ", k, "=", v)
			}
			if e.Original != "" {
				_, _ = fmt.Fprint(st, " original=", e.Original)
			}
			if e.Err != nil {
				_, _ = fmt.Fprint(st, " cause=", e.Err)
			}

			_, _ = fmt.Fprint(st, " •", TrimFilePath(e.OrigFile), ":", e.OrigLine)
			for _, log := range e.Logs {
				_, _ = fmt.Fprint(st, "\n• ", log.Message, " •", TrimFilePath(log.File), ":", strconv.Itoa(log.Line))
				for _, v := range log.Fields {
					_, _ = fmt.Fprint(st, " ", v.Key, "=", logline.ValueOf(v))
				}
			}
			fallthrough
		case st.Flag('+'):
			_, _ = fmt.Fprintf(st, "%+v", e.StackTrace())
		default:
			_, _ = io.WriteString(st, e.Message)
		}
	case 's':
		_, _ = io.WriteString(st, e.Message)
	case 'q':
		_, _ = fmt.Fprintf(st, "%q", e.Error())
	}
}

func (e *APIError) WithMeta(key string, value string) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = value
	return e
}

func (e *APIError) WithMetab(key string, value []byte) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = string(value)
	return e
}

func (e *APIError) WithMetaJson(key string, value interface{}) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	data, err := jsonx.Marshal(value)
	if err != nil {
		e.Meta[key] = fmt.Sprint(err)
	} else {
		e.Meta[key] = string(data)
	}
	return e
}

func (e *APIError) WithMetaID(key string, value interface{ Int64() int64 }) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = strconv.FormatInt(value.Int64(), 10)
	return e
}

func (e *APIError) WithMetap(key string, value interface{}) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = fmt.Sprint(value)
	return e
}

func (e *APIError) WithMetaf(key string, format string, args ...interface{}) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = fmt.Sprintf(format, args...)
	return e
}

func (e *APIError) WithMetaM(m map[string]string) *APIError {
	if e.Meta == nil {
		e.Meta = m
	} else {
		for k, v := range m {
			e.Meta[k] = v
		}
	}
	return e
}

// MarshalJSON ...
func (e *APIError) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 2048)

	b = append(b, '{')
	b = append(b, `"code":`...)
	b = append(b, marshal(e.Code.String())...)

	if e.XCode != 0 {
		b = append(b, ',')
		b = append(b, `"xcode":`...)
		b = append(b, marshal(e.XCode.String())...)
	}

	if e.Err != nil {
		b = append(b, ',')
		b = append(b, `"err":`...)
		b = append(b, marshal(e.Err.Error())...)
	}

	b = append(b, ',')
	b = append(b, `"msg":`...)
	b = append(b, marshal(e.Message)...)

	if e.Original != "" {
		b = append(b, ',')
		b = append(b, `"orig":`...)
		b = append(b, marshal(e.Original)...)
	}

	b = append(b, ',')
	b = append(b, `"logs":`...)
	b = append(b, '[')
	for i, line := range e.Logs {
		if i > 0 {
			b = append(b, ',')
		}
		b = line.MarshalTo(b)
	}
	b = append(b, ']')

	if e.Trace {
		b = append(b, ',')
		b = append(b, `"stack":`...)
		b = append(b, marshal(fmt.Sprintf("%+v", e.Stack))...)
	}

	b = append(b, '}')
	return b, nil
}

func TrimFilePath(file string) string {
	const commonPath = "o.o/backend/"
	if idx := strings.Index(file, commonPath); idx >= 0 {
		file = file[idx+len(commonPath):]
	}
	return file
}

func PrintStack(err error) {
	if s, ok := err.(IStack); ok {
		fmt.Printf("%+v", s.StackTrace())
	} else if err != nil {
		fmt.Printf("%+v", err)
	}
}

func marshal(v interface{}) []byte {
	data, err := jsonx.Marshal(v)
	if err != nil {
		data, _ = jsonx.Marshal(err)
	}
	return data
}

func ErrorCode(err error) Code {
	if err == nil {
		return NoError
	}
	if err, ok := err.(*APIError); ok {
		return err.Code
	}
	if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
		return NoError
	}
	return Unknown
}

func ErrorStack(err error) errors.StackTrace {
	if s, ok := err.(IStack); ok {
		return s.StackTrace()
	}
	return nil
}

// IsTrace ...
func IsTrace(err error) bool {
	if err == nil {
		return false
	}
	if err, ok := err.(*APIError); ok {
		return err.Trace
	}
	return true
}

// GetTraceLevel ...
func GetTraceLevel(err error) TraceLevel {
	if err == nil {
		return LevelNoError
	}
	if xerr, ok := err.(*APIError); ok {
		switch {
		case xerr.XCode == RuntimePanic:
			return LevelPanic
		case xerr.Trace:
			return LevelTrace
		case xerr.Code == Internal:
			return LevelInternal
		case xerr.Trivial:
			return LevelTrival
		default:
			return LevelError
		}
	}
	return LevelInternal
}

type MapErrors struct {
	Error  error
	Code   Code
	Result *APIError
}

type MapErrorItem struct {
	Code    Code
	Message string
}

func MapError(err error) *MapErrors {
	return &MapErrors{
		Error: err,
		Code:  ErrorCode(err),
	}
}

func (m *MapErrors) Wrap(code Code, message string) *MapErrors {
	return m.mapError(code, code, message, false)
}

func (m *MapErrors) Wrapf(code Code, message string, args ...interface{}) *MapErrors {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}
	return m.mapError(code, code, message, false)
}

func (m *MapErrors) Map(code Code, newCode Code, message string) *MapErrors {
	return m.mapError(code, newCode, message, false)
}

func (m *MapErrors) Mapf(code Code, newCode Code, message string, args ...interface{}) *MapErrors {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}
	return m.mapError(code, newCode, message, false)
}

func (m *MapErrors) MapTrace(code Code, newCode Code, message string) *MapErrors {
	return m.mapError(code, newCode, message, true)
}

func (m *MapErrors) MapTracef(code Code, newCode Code, message string, args ...interface{}) *MapErrors {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}
	return m.mapError(code, newCode, message, true)
}

func (m *MapErrors) mapError(code Code, newCode Code, message string, trace bool) *MapErrors {
	if m.Code == code {
		m.Result = newError(trace, true, newCode, message, m.Error)
	}
	return m
}

func (m *MapErrors) DefaultInternal() *APIError {
	return m.defaultError(false, Internal, "")
}

func (m *MapErrors) DefaultInternalTrace() *APIError {
	return m.defaultError(true, Internal, "")
}

func (m *MapErrors) Default(code Code, message string) *APIError {
	return m.defaultError(false, code, message)
}

func (m *MapErrors) DefaultTrace(code Code, message string) *APIError {
	return m.defaultError(true, code, message)
}

func (m *MapErrors) Throw() *APIError {
	if m.Result != nil {
		return m.Result
	}
	if err, ok := m.Error.(*APIError); ok {
		return err
	}
	return newError(true, true, Internal, "", m.Error)
}

func (m *MapErrors) defaultError(trace bool, code Code, message string) *APIError {
	if m.Result != nil {
		return m.Result
	}
	if code == 0 {
		code = m.Code
	}
	if message == "" {
		if xerr, ok := m.Error.(*APIError); ok {
			message = xerr.Message
		} else {
			message = m.Error.Error()
		}
	}
	return newError(trace, true, code, message, m.Error)
}

func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func FirstErrorWithMsg(msg string, errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return Errorf(ErrorCode(err), err, "%v: %v", msg, err)
		}
	}
	return nil
}

type Errors []error

func (errs Errors) ToError() error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errs
	}
}

func (errs Errors) Error() string {
	switch len(errs) {
	case 0:
		return ""
	case 1:
		e := errs[0]
		if e == nil {
			return "ok"
		}
		return e.Error()
	}

	var b strings.Builder
	for i, e := range errs {
		if i > 0 {
			b.WriteString("; ")
		}
		if e == nil {
			b.WriteString("ok")
		} else {
			b.WriteString(e.Error())
		}
	}
	return b.String()
}

func (errs Errors) NErrors() int {
	c := 0
	for _, err := range errs {
		if err != nil {
			c++
		}
	}
	return c
}

func (errs Errors) IsAll() bool {
	if len(errs) == 0 {
		return false
	}
	for _, err := range errs {
		if err == nil {
			return false
		}
	}
	return true
}

func (errs Errors) HasAny() bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

func (errs Errors) All() error {
	if errs.IsAll() {
		return errs
	}
	return nil
}

func (errs Errors) Any() error {
	if errs.HasAny() {
		return errs
	}
	return nil
}

func (errs Errors) Last() error {
	if len(errs) == 0 {
		return nil
	}
	return errs[len(errs)-1]
}

// ErrorCollector ...
type ErrorCollector struct {
	Count  int
	Errors Errors
}

func (e *ErrorCollector) CollectOne(err error) *ErrorCollector {
	e.Count++
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
	return e
}

// Collect collects errors
func (e *ErrorCollector) Collect(errs ...error) *ErrorCollector {
	for _, err := range errs {
		e.Count++
		if err != nil {
			e.Errors = append(e.Errors, err)
		}
	}
	return e
}

// All returns error if all results collected are error
func (e *ErrorCollector) All() error {
	if len(e.Errors) == 0 || e.Count > len(e.Errors) {
		return nil
	}
	return e.concat()
}

// Any returns error if any result collected is error
func (e *ErrorCollector) Any() error {
	return e.concat()
}

func (e *ErrorCollector) N() int {
	return e.Count
}

func (e *ErrorCollector) NErrors() int {
	return len(e.Errors)
}

func (e *ErrorCollector) concat() error {
	return e.Errors.ToError()
}

func (e *ErrorCollector) Last() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e.Errors[len(e.Errors)-1]
}

type ErrorInterface interface {
	// Code is of the valid error codes.
	Code() Code

	// Msg returns a human-readable, unstructured messages describing the error.
	Msg() string

	// WithMeta returns a copy of the Error with the given key-value pair attached
	// as metadata. If the key is already set, it is overwritten.
	WithMeta(key string, val string) ErrorInterface

	// Meta returns the stored value for the given key. If the key has no set
	// value, Meta returns an empty string. There is no way to distinguish between
	// an unset value and an explicit empty string.
	Meta(key string) string

	// MetaMap returns the complete key-value metadata map stored on the error.
	MetaMap() map[string]string

	// Error returns a string of the form "twirp error <Type>: <Msg>"
	Error() string
}

type TwError interface {
	ErrorInterface
	Logs() []*logline.LogLine
	Cause() error
	OrigFile() string
	OrigLine() int
}

type twError struct {
	err *APIError
}

func (t twError) Code() Code {
	return t.err.Code
}

func (t twError) Msg() string {
	return t.err.Message
}

func (t twError) Meta(key string) string {
	meta := t.err.Meta
	if meta != nil {
		return meta[key]
	}
	return ""
}

func (t twError) WithMeta(key string, val string) ErrorInterface {
	_ = t.err.WithMeta(key, val)
	return t
}

func (t twError) MetaMap() map[string]string {
	return t.err.Meta
}

func (t twError) Error() string {
	return t.err.Error()
}

func (t twError) Logs() []*logline.LogLine {
	return t.err.Logs
}

func (t twError) Cause() error {
	if t.err.Err != nil {
		return t.err.Err
	}
	if t.err.Original != "" {
		return errors.New(t.err.Original)
	}
	return nil
}

func (t twError) OrigFile() string {
	return t.err.OrigFile
}

func (t twError) OrigLine() int {
	return t.err.OrigLine
}

func TwirpError(err error) TwError {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(TwError); ok {
		return xerr
	}
	xerr, ok := err.(*APIError)
	if !ok {
		xerr = newError(true, true, Internal, "", err)
	}

	if xerr.Meta == nil {
		xerr.Meta = map[string]string{}
	}
	if xerr.Err != nil {
		xerr.Meta["cause"] = xerr.Err.Error()
	}
	if xerr.Original != "" {
		xerr.Meta["orig"] = xerr.Original
	}
	if xerr.XCode != 0 {
		xerr.Meta["xcode"] = xerr.XCode.String()
	}
	return twError{xerr}
}

type ErrorJSON struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

func ToErrorJSON(twerr ErrorInterface) *ErrorJSON {
	return &ErrorJSON{
		Code: twerr.Code().String(),
		Msg:  twerr.Msg(),
		Meta: twerr.MetaMap(),
	}
}

func (e *ErrorJSON) Error() (s string) {
	if len(e.Meta) == 0 {
		return e.Msg
	}
	b := strings.Builder{}
	b.WriteString(e.Msg)
	b.WriteString(" (")
	for _, v := range e.Meta {
		b.WriteString(v)
		break
	}
	b.WriteString(")")
	return b.String()
}
