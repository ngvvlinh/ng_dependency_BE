package types

import (
	"context"
	"time"

	"o.o/api/etelecom/call_state"
	"o.o/api/main/connectioning"
)

type Driver interface {
	GetTelecomDriver(
		connection *connectioning.Connection,
		shopConnection *connectioning.ShopConnection,
	) (TelecomDriver, error)
}

type TelecomDriver interface {
	Ping(ctx context.Context) error
	GenerateToken(ctx context.Context) (*GenerateTokenResponse, error)
	CreateExtension(ctx context.Context, req *CreateExtensionRequest) (*CreateExtensionResponse, error)

	CreateOutboundRule(context.Context, *CreateOutboundRuleRequest) error

	DestroyCallSession(ctx context.Context, request *DestroyCallSessionRequest) error
}

type TelecomSync interface {
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// use for account admin
type TelecomAdminDriver interface {
	GenerateToken(ctx context.Context) (*GenerateTokenResponse, error)

	CreateTenant(ctx context.Context, req *CreateTenantRequest) (*CreateTenantResponse, error)

	// gán tenant vào trunk provider (aaranet) để biết tenant này trunk qua provider nào và với số hotline nào
	AddHotlineToTenantInTrunkProvider(context.Context, *AddHotlineToTenantInTrunkProviderRequest) error

	// find Trunk provider with tenant_id, then remove hotline out of this tenant
	RemoveHotlineOutOfTenantInTrunkProvider(context.Context, *RemoveHotlineOutOfTenantInTrunkProviderRequest) error
}

type AdministratorTelecom struct {
	Driver                 TelecomAdminDriver
	TokenExpiresAt         time.Time
	TrunkProviderDefaultID string
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
	CallID        string
	CallStatus    string
	Caller        string
	Callee        string
	Direction     string
	StartedAt     time.Time
	EndedAt       time.Time
	Duration      int
	AudioURLs     []string
	CallTargets   []*CallTarget
	CallState     call_state.CallState
	HotlineNumber string
	SessionID     string
}

type CallTarget struct {
	TargetNumber string
	TalkDuration int
	CallState    call_state.CallState
	AnsweredTime time.Time
	EndedTime    time.Time
}

type CreateTenantRequest struct {
	Name     string
	Domain   string
	Password string
	Enable   bool
	Profile  TenantProfile
}

type TenantProfile struct {
	FirstName string
	LastName  string
	Email     string
}

type CreateTenantResponse struct {
	ID string
}

type CreateOutboundRuleRequest struct {
	// trunk provider: aarenat provider id in portsip
	TrunkProviderID string
}

type DestroyCallSessionRequest struct {
	SessionID int
}
type AddHotlineToTenantInTrunkProviderRequest struct {
	// trunk provider: aarenat provider id in portsip
	TrunkProviderID string
	// external tenant ID
	TenantID string
	// hotline number
	Hotline string
}

type RemoveHotlineOutOfTenantInTrunkProviderRequest struct {
	// trunk provider: aarenat provider id in portsip
	TrunkProviderID string
	// external tenant ID
	TenantID string
	// hotline number
	Hotline string
}
