package cm

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"etop.vn/api/meta"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

type EnvType int

// Environment constants
const (
	EnvDev EnvType = iota + 1
	EnvStag
	EnvSandbox
	EnvProd

	// Environment Partner constants
	PartnerEnvTest = "test"
	PartnerEnvDev  = "dev"
	PartnerEnvProd = "prod"
)

var envValues = map[string]EnvType{
	"dev":     EnvDev,
	"stag":    EnvStag,
	"sandbox": EnvSandbox,
	"prod":    EnvProd,
}

var envNames = map[EnvType]string{
	EnvDev:     "dev",
	EnvStag:    "stag",
	EnvSandbox: "sandbox",
	EnvProd:    "prod",
}

func (e EnvType) String() string {
	return envNames[e]
}

var (
	ll      = l.New()
	isDev   = false
	notProd = true
	env     EnvType

	// https://etop.vn
	mainSiteBaseUrl string

	// Commit stores commit
	commit string
)

func init() {
	commit = strings.ReplaceAll(commit, "·", " ")
	commit = strings.ReplaceAll(commit, "¶", "\n")
}

func CommitMessage() string {
	return commit
}

func SetEnvironment(e string) {
	if env != 0 {
		ll.Panic("Already initialize environment")
	}

	env = envValues[e]
	switch env {
	case EnvDev:
		isDev = true
		jsonx.EnableValidation(jsonx.Panicking)

	case EnvStag:
		jsonx.EnableValidation(jsonx.Panicking)

	case EnvSandbox:

	case EnvProd:
		notProd = false

	default:
		ll.S.Panicf("invalid environment: %v", e)
	}
}

func Env() EnvType {
	return env
}

func IsDev() bool {
	return isDev
}

func IsDevOrStag() bool {
	return notProd && env != EnvSandbox
}

func NotProd() bool {
	return notProd
}

func IsProd() bool {
	return !notProd
}

func PartnerEnv() string {
	switch env {
	case EnvDev, EnvStag, EnvSandbox:
		return PartnerEnvTest
	case EnvProd:
		return PartnerEnvProd
	default:
		return env.String()
	}
}

func SetMainSiteBaseURL(s string) {
	if mainSiteBaseUrl != "" {
		ll.Panic("already init base url")
	}

	re := regexp.MustCompile(`^https?://[^/]+$`)
	if !re.MatchString(s) {
		ll.Panic("Invalid base url", l.String("s", s))
	}
	mainSiteBaseUrl = s
}

func MainSiteBaseURL() string {
	return mainSiteBaseUrl
}

func Parallel(ctx context.Context, jobs ...func(context.Context) error) []error {
	wg := sync.WaitGroup{}
	wg.Add(len(jobs))
	errors := make([]error, len(jobs))

	for i := range jobs {
		job := jobs[i]
		go func(i int) {
			defer wg.Done()
			errors[i] = job(ctx)
		}(i)
	}

	wg.Wait()
	return errors
}

type Paging = meta.Paging
type Filter = meta.Filter
