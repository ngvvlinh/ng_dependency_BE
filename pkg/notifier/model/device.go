package model

type GetDeviceArgs struct {
	AccountID int64
	ID        int64
}

type GetDevicesArgs struct {
	AccountID         int64
	ExternalServiceID int
}

type CreateDeviceArgs struct {
	AccountID        int64
	DeviceID         string
	DeviceName       string
	ExternalDeviceID string
}
