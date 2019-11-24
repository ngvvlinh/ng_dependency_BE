package model

import "etop.vn/capi/dot"

type GetDeviceArgs struct {
	UserID           dot.ID
	ExternalDeviceID string
}

type GetDevicesArgs struct {
	UserID            dot.ID
	ExternalServiceID int
}

type CreateDeviceArgs struct {
	UserID            dot.ID
	AccountID         dot.ID
	DeviceID          string
	DeviceName        string
	ExternalDeviceID  string
	ExternalServiceID int
	Config            *DeviceConfig
}

type UpdateDeviceArgs struct {
	UserID           dot.ID
	ExternalDeviceID string
	Config           *DeviceConfig
}
