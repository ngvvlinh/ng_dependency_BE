package types

import (
	"context"
	"time"

	"o.o/api/main/connectioning"
)

type Driver interface {
	GetTelecomDriver(
		env string,
		connection *connectioning.Connection,
		shopConnection *connectioning.ShopConnection,
	) (TelecomDriver, error)
}

type TelecomDriver interface {
	Ping(ctx context.Context) error
	GenerateToken(ctx context.Context) (*GenerateTokenResponse, error)
	CreateExtension(ctx context.Context, req *CreateExtensionRequest) (*CreateExtensionResponse, error)
	GetCallLogs(ctx context.Context, req *GetCallLogsRequest) (*GetCallLogsResponse, error)
}

type TelecomSync interface {
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type GenerateTokenResponse struct {
	AccessToken string
	ExpiresAt   time.Time
	TokenType   string
	ExpiresIn   int
}

type CreateExtensionRequest struct {
	ExtensionNumber   string
	ExtensionPassword string
	Profile           *ProfileExtension
	Hotline           string
}

type ProfileExtension struct {
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Description string
}

type CreateExtensionResponse struct {
	ID string
}

type GetCallLogsRequest struct {
	StartedAt time.Time
	EndedAt   time.Time
	ScrollID  string // VHT pagination
}

type GetCallLogsResponse struct {
	CallLogs []*CallLog
	ScrollID string // VHT pagination
}

type CallLog struct {
	CallID     string
	CallStatus string
	Caller     string
	Callee     string
	Direction  string
	StartedAt  time.Time
	EndedAt    time.Time
	Duration   int
	AudioURLs  []string
}
