package cm

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var (
	ll = l.New()

	// https://etop.vn
	mainSiteBaseUrl string

	// Commit stores commit
	commit string
)

func init() {
	commit = strings.ReplaceAll(commit, "·", " ")
	commit = strings.ReplaceAll(commit, "¶", "\n")

	cmenv.SetHook(func(env cmenv.EnvType) error {
		switch env {
		case cmenv.EnvDev, cmenv.EnvStag:
			jsonx.EnableValidation(jsonx.Panicking)
		}
		return nil
	})
}

func CommitMessage() string {
	return commit
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
