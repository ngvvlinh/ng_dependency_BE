package model

import (
	"time"

	"o.o/api/top/types/etc/status4"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

type ImportType string

const (
	ImportTypeShopOrder   ImportType = "import_shop_order"
	ImportTypeShopProduct ImportType = "import_shop_product"
)

func (t ImportType) String() string { return string(t) }

// +sqlgen
type ImportAttempt struct {
	ID           dot.ID
	UserID       dot.ID
	AccountID    dot.ID
	OriginalFile string
	StoredFile   string
	Type         ImportType
	NCreated     int
	NUpdated     int
	NError       int
	Status       status4.Status
	ErrorType    string
	Errors       []*Error
	DurationMs   int
	CreatedAt    time.Time `sq:"create"`
}

type CreateImportAttemptCommand struct {
	ImportAttempt *ImportAttempt
}

// +sqlgen
type ExportAttempt struct {
	ID        string
	UserID    dot.ID
	AccountID dot.ID

	ExportType   string
	FileName     string `sq:"'filename'"`
	StoredFile   string
	DownloadURL  string
	RequestQuery string
	MimeType     string

	Status    status4.Status
	Errors    []*Error
	Error     *Error
	NTotal    int
	NExported int
	NError    int

	CreatedAt time.Time `sq:"create"`
	DeletedAt time.Time
	StartedAt time.Time
	DoneAt    time.Time
	ExpiresAt time.Time
}

func (e *ExportAttempt) GetAbortedError() error {
	if e.Status != status4.N {
		return nil
	}
	if len(e.Errors) == 0 {
		return cm.Errorf(cm.Unknown, nil, "Lỗi không xác định khi export")
	}
	return e.Errors[len(e.Errors)-1]
}
