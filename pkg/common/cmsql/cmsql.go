package cmsql

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"sync"
	"time"

	"github.com/lib/pq"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

// map[connection-string]Database
var dbPool = map[string]*Database{}
var mu sync.RWMutex
var ll = l.New()

const sqlScope = "https://www.googleapis.com/auth/sqlservice.admin"

// Config ...
type ConfigPostgres struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
	Timeout  int    `yaml:"timeout"`

	MaxOpenConns    int `yaml:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns"`
	MaxConnLifetime int `yaml:"max_conn_lifetime"`

	GoogleAuthFile string `yaml:"google_auth_file"`
}

// RegisterCloudSQL ...
func (c *ConfigPostgres) RegisterCloudSQL() error {
	switch c.Protocol {
	case "":
		return nil
	case "cloudsql":
	default:
		ll.Panic("Invalid postgres protocol")
	}

	conf, err := ioutil.ReadFile(c.GoogleAuthFile)
	if err != nil {
		ll.Error("Unable to read auth file", l.String("file", c.GoogleAuthFile), l.Error(err))
		return err
	}
	config, err := google.JWTConfigFromJSON(conf, sqlScope)
	if err != nil {
		ll.Error("Failed to decode auth file", l.String("GoogleAuthFile", c.GoogleAuthFile), l.Error(err))
		return err
	}

	ctx := context.Background()
	proxy.Init(config.Client(ctx), nil, nil)
	return nil
}

// ConnectionString ...
func (c *ConfigPostgres) ConnectionString() (driver string, connStr string) {
	sslmode := c.SSLMode
	if c.SSLMode == "" {
		sslmode = "disable"
	}
	if c.Timeout == 0 {
		c.Timeout = 15
	}

	switch c.Protocol {
	case "":
		connStr = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v connect_timeout=%v", c.Host, c.Port, c.Username, c.Password, c.Database, sslmode, c.Timeout)
	case "cloudsql":
		connStr = fmt.Sprintf("host=%v user=%v password=%v dbname=%v sslmode=disable",
			c.Host, c.Username, c.Password, c.Database)
	default:
		ll.Panic("postgres: Invalid protocol", l.Object("config", c))
	}
	return c.Driver(), connStr
}

func (c *ConfigPostgres) Driver() string {
	switch c.Protocol {
	case "":
		return "postgres"
	case "cloudsql":
		return "cloudsqlpostgres"
	default:
		return "unknown"
	}
}

func (c *ConfigPostgres) connectionStringIdentifier() string {
	return fmt.Sprintf("driver=%v host=%v port=%v user=%v dbname=%v",
		c.Driver(), c.Host, c.Port, c.Username, c.Database)
}

// Dialer ...
type Dialer struct {
}

// NewDialer ...
func NewDialer() *Dialer {
	return &Dialer{}
}

// [extopvn:asia-southeast1:etoppg1]:5432
var reAddress = regexp.MustCompile(`\[([^]]+)]`)

// Dial ...
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	parts := reAddress.FindStringSubmatch(address)
	if len(parts) != 2 {
		ll.Fatal("Unexpected address", l.String("address", address))
	}

	conn, err := proxy.Dial(parts[1])
	if err != nil {
		ll.Fatal("Unable to listen dial postgres host", l.String("host", parts[1]), l.Error(err))
	}
	return conn, nil
}

// DialTimeout ...
func (d *Dialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	parts := reAddress.FindStringSubmatch(address)
	if len(parts) != 2 {
		ll.Fatal("Unexpected address", l.String("address", address))
	}

	conn, err := proxy.Dial(parts[1])
	if err != nil {
		ll.Fatal("Unable to listen dial postgres host", l.String("host", parts[1]), l.Error(err))
	}
	return conn, nil
}

// Database ...
type Database struct {
	id   int64
	db   sq.Database
	dlog sq.DynamicLogger
	errs []error
}

func (db *Database) GetSchemaErrors() error {
	for _, v := range db.errs {
		l.Error(v)
	}
	if len(db.errs) > 0 {
		return cm.Error(cm.Internal, "Found error in database migration", nil)
	}
	return nil
}

func (db *Database) RecordError(err error) {
	db.errs = append(db.errs, err)
}

func MustConnect(c ConfigPostgres) *Database {
	db, err := Connect(c)
	if err != nil {
		ll.Fatal("Error while connecting to database", l.Error(err))
	}
	return db
}

// Connect ...
func Connect(c ConfigPostgres) (*Database, error) {
	if err := c.RegisterCloudSQL(); err != nil {
		return &Database{}, err
	}

	identifier := c.connectionStringIdentifier()
	mu.RLock()
	if db, ok := dbPool[identifier]; ok {
		mu.RUnlock()
		return db, nil
	}
	mu.RUnlock()

	dlog := sq.NewDynamicLogger(DefaultLogger)
	driver, conn := c.ConnectionString()
	db, err := sq.Connect(driver, conn, dlog,
		sq.SetErrorMapper(DefaultErrorMapper))
	if err != nil {
		return &Database{}, err
	}
	if _, err := db.Exec("SELECT 1"); err != nil {
		return &Database{}, err
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 10
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 3
	}
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetMaxIdleConns(c.MaxIdleConns)

	mu.Lock()
	database := &Database{id: cm.NewID(), db: *db, dlog: *dlog}
	dbPool[identifier] = database
	defer mu.Unlock()
	return database, nil
}

func (db Database) TxKey() TxKey {
	return TxKey{db.id}
}

// DB ...
func (db Database) DB() *sq.Database {
	return &db.db
}

func (db Database) Opts() core.Opts {
	return db.db.Opts()
}

// Exec ...
func (db Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

// ExecContext ...
func (db Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// MustExec ...
func (db Database) MustExec(query string, args ...interface{}) sql.Result {
	res, err := db.db.Exec(query, args...)
	if err != nil {
		ll.Fatal("Unable to execute query", l.Error(err))
	}
	return res
}

// Query ...
func (db Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

// QueryContext ...
func (db Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

// QueryRow ...
func (db Database) QueryRow(query string, args ...interface{}) sq.Row {
	return db.db.QueryRow(query, args...)
}

// QueryRowContext ...
func (db Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) sq.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

// NewQuery ...
func (db Database) NewQuery() Query {
	return Query{db.db.NewQuery()}
}

func (db Database) WithContext(ctx context.Context) Query {
	return Query{db.db.NewQuery().WithContext(ctx)}
}

// Begin ...
func (db Database) Begin() (Tx, error) {
	t, err := db.db.Begin()
	return tx{tx: t, db: db}, err
}

func (db Database) BeginContext(ctx context.Context) (Tx, error) {
	t, err := db.db.BeginContext(ctx)
	return tx{tx: t, db: db}, err
}

// Query ...
type Query struct {
	query sq.Query
}

// WithContext ...
func (q Query) WithContext(ctx context.Context) Query {
	return Query{q.query.WithContext(ctx)}
}

// Clone ...
func (q Query) Clone() Query {
	return Query{q.query.Clone()}
}

// Exec ...
func (q Query) Exec() (sql.Result, error) {
	return q.query.Exec()
}

// Query ...
func (q Query) Query() (*sql.Rows, error) {
	return q.query.Query()
}

// QueryRow ...
func (q Query) QueryRow() (sq.Row, error) {
	return q.query.QueryRow()
}

// Scan ...
func (q Query) Scan(dest ...interface{}) error {
	return q.query.Scan(dest...)
}

type tx struct {
	tx sq.Tx
	db Database
}

func (tx tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.db.Exec(query, args...)
}

func (tx tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.db.ExecContext(ctx, query, args...)
}

func (tx tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.db.Query(query, args...)
}

func (tx tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.db.Query(query, args...)
}

func (tx tx) QueryRow(query string, args ...interface{}) sq.Row {
	return tx.db.QueryRow(query, args...)
}

func (tx tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) sq.Row {
	return tx.db.QueryRow(query, args...)
}

func (tx tx) Commit() error {
	return tx.tx.Commit()
}

func (tx tx) Rollback() error {
	return tx.tx.Rollback()
}

// DefaultLogger ...
func DefaultLogger(entry *sq.LogEntry) {
	if ctx, ok := entry.Ctx.(*bus.NodeContext); ok {
		ctx.WithMessage(entry)
	}
	if entry.IsTx() {
		if entry.Type() == sq.TypeCommit || entry.Type() == sq.TypeRollback {
			for _, item := range entry.TxQueries {
				logQuery(item)
			}
		}
		return
	}
	logQuery(entry)
}

func logQuery(entry *sq.LogEntry) {
	d, query, err := entry.Duration, entry.Query, entry.Error
	args, _ := entry.Args.ToSQLValues()
	if err != nil && xerrors.ErrorCode(err) != xerrors.NotFound {
		ll.Error(query, l.Any("args", args), l.Error(err), l.Duration("t", d), l.Bool("tx", entry.IsTx()))
	} else if d >= 200*time.Millisecond {
		ll.Warn(query, l.Any("args", args), l.Duration("t", d), l.Bool("tx", entry.IsTx()))
	} else {
		ll.Debug(query, l.Any("args", args), l.Duration("t", d), l.Bool("tx", entry.IsTx()))
	}
}

// DefaultErrorMapper ...
func DefaultErrorMapper(err error, entry *sq.LogEntry) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return xerrors.Error(xerrors.NotFound, "", err)
	default:
		if err, ok := err.(core.InvalidArgumentError); ok {
			return xerrors.Error(xerrors.InvalidArgument, string(err), nil)
		}
		return xerrors.Error(xerrors.Internal, "", err)
	}
}

// InTransaction ...
func (db Database) InTransaction(ctx context.Context, callback func(QueryInterface) error) (err error) {
	txKey := db.TxKey()
	{
		tx := ctx.Value(txKey)
		if tx != nil {
			return callback(tx.(Tx))
		}
	}

	tx, err := db.BeginContext(ctx)
	if err != nil {
		return err
	}
	ctx2, ok := ctx.(bus.WithValuer)
	if !ok {
		panic("cmsql: InTransaction only accepts bus.NodeContext")
	}
	ctx2.WithValue(txKey, tx)

	defer func() {
		e := recover()
		if e != nil {
			ll.Error("common/sql: panic (recover)", l.Any("err", e), l.Stack())
			err = xerrors.Error(xerrors.Internal, fmt.Sprint(e), nil)
		}
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
		ctx2.ResetValue(txKey)
	}()
	return callback(tx)
}

func DefaultListenerProblemReport(ev pq.ListenerEventType, err error) {
	if err != nil {
		ll.Error("Error while listening to event", l.Int("type", int(ev)), l.Error(err))
	}
}

func NewListener(
	cfg ConfigPostgres,
	minReconnectInterval time.Duration,
	maxReconnectInterval time.Duration,
	eventCallback pq.EventCallbackType,
) *pq.Listener {

	switch cfg.Protocol {
	case "":
		_, str := cfg.ConnectionString()
		listener := pq.NewListener(str, minReconnectInterval, maxReconnectInterval, eventCallback)
		return listener

	case "cloudsql":
		dialer := NewDialer()
		_, str := cfg.ConnectionString()
		listener := pq.NewDialListener(dialer, str, minReconnectInterval, maxReconnectInterval, eventCallback)
		return listener

	default:
		ll.Panic("Invalid protocol", l.Object("config", cfg))
		return nil
	}
}

func ListenTo(ctx context.Context, listener *pq.Listener, channels ...string) error {
	var g errgroup.Group
	for _, c := range channels {
		g.Go(func() error {
			return listener.Listen(c)
		})
	}
	return g.Wait()
}
