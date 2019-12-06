package try_on

// +enum
type TryOnCode int

const (
	// +enum=unknown
	Unknown TryOnCode = 0

	// +enum=none
	None TryOnCode = 1

	// +enum=open
	Open TryOnCode = 2

	// +enum=try
	Try TryOnCode = 3
)
