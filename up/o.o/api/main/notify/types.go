package notify

import "o.o/capi/dot"

type UserNotiSetting struct {
	UserID        dot.ID
	DisableTopics []string
}
