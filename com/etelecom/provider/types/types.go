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
	GetCallLogs(ctx context.Context) ([]*CallLog, error)
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

type CallLog struct {
	CallID       string
	CallStatus   string
	Callee       string
	CalleeDomain string
	StartTime    string
	EndTime      string
	TaskDuration string
	AudioURLs    []string
}
