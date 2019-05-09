// Package ssm (sync state management) provides logic for managing sync state. Currently it's used only by Kiotviet.
package ssm

import (
	"time"

	"etop.vn/backend/pkg/common/l"
)

var ll = l.New()

// SyncState represents `since` and `page`, which is a common pattern to work
// with offset-and-time-based API.
type SyncState struct {
	Since time.Time `json:"since"`
	Page  int       `json:"page"`
}

// IsValid reports whether the state is valid.
func (s SyncState) IsValid() bool {
	return !s.Since.IsZero() && s.Page >= 1
}

// AdvanceArgs represents arguments for AdvanceState.
type AdvanceArgs struct {
	Done  bool
	Size  int
	Start time.Time
	End   time.Time
}

// AdvanceState calculates next SyncState from current SyncState and response
// data.
func AdvanceState(s SyncState, delta time.Duration, args AdvanceArgs) (next SyncState) {
	done, size, start, end := args.Done, args.Size, args.Start, args.End
	since, page := s.Since, s.Page

	if since.IsZero() {
		ll.Panic("Unexpected zero time", l.Any("s", s), l.Any("args", args))
	}
	defer func() {
		ll.Info("AdvanceState", l.Object("s", s), l.Object("args", args), l.Object("next", next))
	}()

	switch {
	case done && size == 0 && page == 1:
		// If the request returns nothing, it's safe to repeat the request next
		// time.
		return s
	case done && size == 0 && page > 1:
		if end.IsZero() {
			return SyncState{since.Add(delta), 1}
		}
		return SyncState{end.Add(delta), 1}
	case done && size > 0:
		// Due to a bug from Kiotviet, which returns zero date even if it's
		// requested with `modifiedDateFrom`, we will only advance `end` if it's
		// not zero.
		//
		// It's `since`, not `since.Add(delta)`, because Kiotviet may return
		// zero `modifiedDate`. In that case, it's better to just repeat the
		// request instead of advance `since`. Otherwise, if we advance `since`,
		// the result is not idempotent (when we repeatly call `AdvanceState()`,
		// `since` will keep increasing).
		//
		// This also works with Haravan, because Haravan does not allow to sort
		// by `updated_at`, only allow to query by `updated_at_min`. Therefore
		// when `start` and `end` are not provided, the only option for us is
		// increasing `page` (when not done).
		if end.IsZero() {
			return SyncState{since, 1}
		}
		return SyncState{end.Add(delta), 1}
	// Due to a bug from Kiotviet, which returns zero date even if it's
	// requested with modifiedDateFrom, we will always advance page until the
	// first item has non-zero date.
	case !done && start.IsZero():
		return SyncState{since, page + 1}
	case !done && start.Equal(end):
		// This also works in case both `start` and `end` are zero (or not
		// provided). The only option for us is increasing `page`.
		return SyncState{since, page + 1}
	case !done && start.Before(end):
		return SyncState{end, 1}
	}

	ll.Panic("Unexpected", l.Any("s", s), l.Any("args", args))
	return SyncState{}
}
