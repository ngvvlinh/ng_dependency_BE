package cc

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var ll = l.New()

var (
	flConfigFile = ""
	flConfigYaml = ""
	flExample    = false
	flNoEnv      = false
	flCommitMsgs = false
)

func InitFlags() {
	flag.StringVar(&flConfigFile, "config-file", "", "Path to config file")
	flag.StringVar(&flConfigYaml, "config-yaml", "", "Config as yaml string")
	flag.BoolVar(&flNoEnv, "no-env", false, "Don't read config from environment")
	flag.BoolVar(&flExample, "example", false, "Print example config then exit")
	flag.BoolVar(&flCommitMsgs, "commit-messages", false, "Print commit messages then exit")
}

func ParseFlags() {
	flag.Parse()
	if flCommitMsgs {
		fmt.Println(cm.CommitMessage())
		os.Exit(2)
	}
}

func FlagConfigFile() string {
	return flConfigFile
}

func FlagExample() bool {
	return flExample
}

func FlagNoEnv() bool {
	return flNoEnv
}

func LoadWithDefault(v, def interface{}) (err error) {
	defer func() {
		if flExample {
			if err != nil {
				ll.Fatal("Error while loading config", l.Error(err))
			}
			PrintExample(v)
			os.Exit(2)
		}
	}()

	if (flConfigFile != "") && (flConfigYaml != "") {
		return errors.New("must provide only -config-file or -config-yaml")
	}
	if flConfigFile != "" {
		err := LoadFromFile(flConfigFile, v)
		if err != nil {
			ll.S.Errorf("can not load config from file: %v", flConfigFile)
		}
		return err
	}
	if flConfigYaml != "" {
		return LoadFromYaml([]byte(flConfigYaml), v)
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(def))
	return nil
}

// LoadFromFile loads config from file
func LoadFromFile(configPath string, v interface{}) (err error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return LoadFromYaml(data, v)
}

func LoadFromYaml(input []byte, v interface{}) (err error) {
	return yaml.Unmarshal(input, v)
}

func EnvPrefix(prefix []string, def string) string {
	if len(prefix) > 0 {
		return prefix[1]
	}
	return def
}

type EnvMap map[string]interface{}

func (m EnvMap) MustLoad() {
	if flNoEnv {
		return
	}

	for k, v := range m {
		MustLoadEnv(k, v)
	}
}

func MustLoadEnv(env string, val interface{}) {
	if flNoEnv {
		return
	}

	s := os.Getenv(env)
	if s == "" {
		return
	}

	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		ll.S.Panicf("Expect pointer for env: %v", env)
	}

	v = v.Elem()
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			ll.S.Fatalf("%v expects a boolean, got: %v", s)
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ll.S.Fatalf("%v expects an integer, got: %v", s)
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			ll.S.Fatalf("%v expects an unsigned integer, got: %v", s)
		}
		v.SetUint(i)
	case reflect.Float64:
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			ll.S.Fatalf("%v expects an float64, got: %v", s)
		}
		v.SetFloat(i)
	default:
		ll.S.Panicf("Unexpected type for env: %v", s)
	}
}

// PrintExample prints example config
func PrintExample(cfg interface{}) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		ll.Fatal(err.Error())
	}
	fmt.Println(string(data))
}

type Postgres = cmsql.ConfigPostgres

// DefaultPostgres ...
func DefaultPostgres() Postgres {
	return Postgres{
		Protocol:       "",
		Host:           "postgres",
		Port:           5432,
		Username:       "postgres",
		Password:       "postgres",
		Database:       "test",
		SSLMode:        "",
		Timeout:        15,
		GoogleAuthFile: "",
	}
}

func PostgresMustLoadEnv(c *cmsql.ConfigPostgres, prefix ...string) {
	p := "ET_POSTGRES"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	EnvMap{
		p + "_PORT":     &c.Port,
		p + "_HOST":     &c.Host,
		p + "_SSLMODE":  &c.SSLMode,
		p + "_USERNAME": &c.Username,
		p + "_PASSWORD": &c.Password,
		p + "_DATABASE": &c.Database,
		p + "_TIMEOUT":  &c.Timeout,
	}.MustLoad()
}

// HTTP ...
type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c *HTTP) MustLoadEnv(prefix string) {
	p := prefix
	EnvMap{
		p + "_HORT": &c.Host,
		p + "_PORT": &c.Port,
	}.MustLoad()
}

// Address ...
func (c *HTTP) Address() string {
	if c.Port == 0 {
		ll.Panic("Missing HTTP port")
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// DefaultRedis ...
func DefaultRedis() Redis {
	return Redis{
		Host:     "redis",
		Port:     "6379",
		Username: "",
		Password: "",
	}
}

func (c *Redis) MustLoadEnv(prefix ...string) {
	p := "ET_REDIS"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	EnvMap{
		p + "_PORT":     &c.Port,
		p + "_HOST":     &c.Host,
		p + "_USERNAME": &c.Username,
		p + "_PASSWORD": &c.Password,
	}.MustLoad()
}

// ConnectionString ...
func (c Redis) ConnectionString() string {
	s := ""
	if c.Username == "" || c.Password == "" {
		s = fmt.Sprintf("redis://%s:%s", c.Host, c.Port)
	} else {
		s = fmt.Sprintf("redis://%s:%s@%s:%s", c.Username, c.Password, c.Host, c.Port)
	}
	return s
}

type TelegramBot struct {
	Token string
	Chats map[string]int64
}

func (c *TelegramBot) MustLoadEnv(prefix ...string) {
	p := "ET_TELEGRAM_BOT"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	EnvMap{
		p + "_TOKEN": &c.Token,
	}.MustLoad()
}

func (c *TelegramBot) ConnectDefault() (*telebot.Channel, error) {
	return c.ConnectChannel("")
}

func (c *TelegramBot) ConnectChannel(channel string) (*telebot.Channel, error) {
	if c.Token == "" {
		ll.Warn("Disabled sending messages to telegram")
		return nil, nil
	}

	cfgName := channel
	if cfgName == "" {
		cfgName = "default"
	}
	ch, err := telebot.NewChannel(c.Token, c.Chats[cfgName])
	if err != nil {
		return nil, err
	}

	ll.Info("Enabled sending messages to telegram")
	telebot.RegisterChannel(channel, ch)
	return ch, err
}

func (c *TelegramBot) MustConnectChannel(channel string) *telebot.Channel {
	ch, err := c.ConnectChannel(channel)
	if err != nil {
		ll.Panic("Connect Telegram bot", l.String("channel", channel), l.Error(err))
	}
	return ch
}

type Kafka struct {
	Enabled     bool     `yaml:"enabled"`
	Brokers     []string `yaml:"brokers"`
	TopicPrefix string   `yaml:"topic_prefix"`
}

func DefaultKafka() Kafka {
	return Kafka{
		TopicPrefix: "d",
		Brokers: []string{
			"kafka:9092",
		},
		Enabled: true,
	}
}

type OnesignalConfig struct {
	ApiKey string `yaml:"api_key"`
	AppID  string `yaml:"app_id"`
}

func DefaultOnesignal() OnesignalConfig {
	return OnesignalConfig{
		ApiKey: "MmVjMTliNTItYzM5Yi00OWFlLThmODctMWM1YWE2YjY0OTEx",
		AppID:  "514a0d7d-2336-4ed8-80da-bc69ec35a19f",
	}
}

func (c *OnesignalConfig) MustLoadEnv(prefix ...string) {
	p := "ET_ONESIGNAL"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	EnvMap{
		p + "_API_KEY": &c.ApiKey,
		p + "_APP_ID":  &c.AppID,
	}.MustLoad()
}
