package model

type GetDeviceArgs struct {
	UserID           int64
	ExternalDeviceID string
}

type GetDevicesArgs struct {
	UserID            int64
	ExternalServiceID int
}

type CreateDeviceArgs struct {
	UserID            int64
	AccountID         int64
	DeviceID          string
	DeviceName        string
	ExternalDeviceID  string
	ExternalServiceID int
	Config            *DeviceConfig
}

type UpdateDeviceArgs struct {
	UserID           int64
	ExternalDeviceID string
	Config           *DeviceConfig
}
