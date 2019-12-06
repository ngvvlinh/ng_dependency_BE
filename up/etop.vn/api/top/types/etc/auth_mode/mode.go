package auth_mode

// +enum
type AuthMode int

const (
	// +enum=default
	Default AuthMode = 0

	// +enum=manual
	Manual AuthMode = 1
)
