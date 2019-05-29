package model

import (
	"time"

	cm "etop.vn/backend/pkg/common"
)

type ImportType string

const (
	ImportTypeShopOrder   ImportType = "shop_order"
	ImportTypeShopProduct ImportType = "shop_product"
)

var _ = sqlgenImportAttempt(&ImportAttempt{})

type ImportAttempt struct {
	ID           int64
	UserID       int64
	AccountID    int64
	OriginalFile string
	StoredFile   string
	Type         ImportType
	NCreated     int
	NUpdated     int
	NError       int
	Status       Status4
	ErrorType    string
	Errors       []*Error
	DurationMs   int
	CreatedAt    time.Time `sq:"create"`
}

type CreateImportAttemptCommand struct {
	ImportAttempt *ImportAttempt
}

var _ = sqlgenExportAttempt(&ExportAttempt{})

type ExportAttempt struct {
	ID        string
	UserID    int64
	AccountID int64

	ExportType   string
	FileName     string `sq:"'filename'"`
	StoredFile   string
	DownloadURL  string
	RequestQuery string
	MimeType     string

	Status    Status4
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
	if e.Status != S4Negative {
		return nil
	}
	if len(e.Errors) == 0 {
		return cm.Errorf(cm.Unknown, nil, "Lỗi không xác định khi export")
	}
	return e.Errors[len(e.Errors)-1]
}