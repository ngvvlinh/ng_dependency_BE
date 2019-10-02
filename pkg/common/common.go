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

// Environment constants
const (
	EnvDev     = "dev"
	EnvStag    = "stag"
	EnvSandbox = "sandbox"
	EnvProd    = "prod"

	// Environment Partner constants
	PartnerEnvTest = "test"
	PartnerEnvProd = "prod"
)

var (
	ll      = l.New()
	isDev   = false
	notProd = true
	env     string

	// https://etop.vn
	mainSiteBaseUrl string

	// Commit stores commit
	commit string
)

func init() {
	commit = strings.ReplaceAll(commit, "·", " ")
	commit = strings.ReplaceAll(commit, "⮐", "\n")
}

func CommitMessage() string {
	return commit
}

func SetEnvironment(e string) {
	if env != "" {
		ll.Panic("Already initialize environment")
	}

	env = e
	switch e {
	case EnvDev:
		isDev = true
		jsonx.EnableValidation(jsonx.Warning)

	case EnvStag:
		jsonx.EnableValidation(jsonx.Warning)

	case EnvSandbox:

	case EnvProd:
		notProd = false
	default:
		ll.S.Panicf("Invalid environment: %v", e)
	}
}

func Env() string {
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
