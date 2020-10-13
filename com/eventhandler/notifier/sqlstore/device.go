package sqlstore

import (
	"fmt"
	"time"

	"github.com/lib/pq"

	"o.o/backend/com/eventhandler/notifier/model"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/capi/util"
)

type DeviceStore struct {
	db *cmsql.Database
}

type M map[string]interface{}

func NewDeviceStore(db com.NotifierDB) *DeviceStore {
	model.SQLVerifySchema(db)
	return &DeviceStore{
		db: db,
	}
}

func (s *DeviceStore) CreateDevice(args *model.CreateDeviceArgs) (*model.Device, error) {
	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing User ID")
	}
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if args.ExternalDeviceID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing External Device ID")
	}
	externalServiceID := args.ExternalServiceID
	if externalServiceID == 0 {
		// Use Onesignal by default
		externalServiceID = model.ExternalServiceOneSignalID
	}

	var dbDevice = new(model.Device)
	ok, err := s.db.Table("device").Where("external_device_id = ? AND user_id = ? AND external_service_id = ?", args.ExternalDeviceID, args.UserID, externalServiceID).Get(dbDevice)
	if err != nil {
		return nil, err
	}
	device := &model.Device{
		AccountID:         args.AccountID,
		UserID:            args.UserID,
		DeviceName:        args.DeviceName,
		DeviceID:          args.DeviceID,
		ExternalDeviceID:  args.ExternalDeviceID,
		ExternalServiceID: externalServiceID,
	}
	defaultConfig := &model.DeviceConfig{
		SubcribeAllShop: true,
		Mute:            false,
	}
	var id dot.ID
	if ok && dbDevice.ID != 0 {
		if !dbDevice.DeactivatedAt.IsZero() {
			// active this device
			if err := s.db.Table("device").Where("id = ?", dbDevice.ID).ShouldUpdateMap(M{"deactivated_at": nil}); err != nil {
				return nil, err
			}
		}
		// update it
		if err := s.db.Table("device").Where("id = ?", dbDevice.ID).ShouldUpdate(device); err != nil {
			return nil, err
		}
	} else {
		// create new device and make sure only one external_device_id is actived at a time
		if _, err := s.db.Table("device").Where("external_device_id = ? AND deactivated_at IS NULL", args.ExternalDeviceID).UpdateMap(M{"deactivated_at": time.Now()}); err != nil {
			return nil, err
		}
		id = cm.NewID()
		device.ID = id
		device.Config = defaultConfig
		if err := s.db.Table("device").ShouldInsert(device); err != nil {
			return nil, err
		}
	}
	res, err := s.GetDevice(&model.GetDeviceArgs{
		UserID:           args.UserID,
		ExternalDeviceID: device.ExternalDeviceID,
	})
	return res, err
}

func (s *DeviceStore) UpdateDevice(args *model.UpdateDeviceArgs) (*model.Device, error) {
	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing User ID")
	}
	if args.ExternalDeviceID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing External Device ID")
	}

	device := &model.Device{
		ExternalDeviceID: args.ExternalDeviceID,
		UserID:           args.UserID,
		Config:           args.Config,
	}
	if err := s.db.Table("device").Where("user_id = ? AND external_device_id = ?", args.UserID, args.ExternalDeviceID).
		ShouldUpdate(device); err != nil {
		return nil, err
	}
	res, err := s.GetDevice(&model.GetDeviceArgs{
		UserID:           args.UserID,
		ExternalDeviceID: args.ExternalDeviceID,
	})
	return res, err
}

func (s *DeviceStore) GetDevice(args *model.GetDeviceArgs) (*model.Device, error) {
	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing User ID")
	}
	if args.ExternalDeviceID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing External Device ID")
	}
	var device = new(model.Device)
	err := s.db.Table("device").Where("user_id = ? AND external_device_id = ? AND deactivated_at IS NULL", args.UserID, args.ExternalDeviceID).ShouldGet(device)
	return device, err
}

func (s *DeviceStore) GetDevices(args *model.GetDevicesArgs) ([]*model.Device, error) {
	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing User ID")
	}
	if args.ExternalServiceID == 0 {
		args.ExternalServiceID = model.ExternalServiceOneSignalID
	}

	var res []*model.Device
	if err := s.db.Table("device").Where("user_id = ? AND external_service_id = ? AND deactivated_at IS NULL", args.UserID, args.ExternalServiceID).Find((*model.Devices)(&res)); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *DeviceStore) GetExternalDeviceIDs(arg *model.GetDevicesArgs) ([]string, error) {
	devices, err := s.GetDevices(arg)
	if err != nil {
		return nil, err
	}
	var deviceIDs = make([]string, 0, len(devices))
	for _, device := range devices {
		if device.ExternalDeviceID != "" {
			deviceIDs = append(deviceIDs, device.ExternalDeviceID)
		}
	}
	return deviceIDs, nil
}

func (s *DeviceStore) DeleteDevice(device *model.Device) error {
	if device.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing User ID")
	}
	// deprecated `device.DeviceID` soon
	if device.ExternalDeviceID == "" && device.DeviceID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing External ID")
	}
	x := s.db.Table("device").Where("user_id = ?", device.UserID)
	if device.DeviceID != "" {
		x = x.Where("device_id = ?", device.DeviceID)
	}
	if device.ExternalDeviceID != "" {
		x = x.Where("external_device_id = ?", device.ExternalDeviceID)
	}
	err := x.ShouldUpdate(&model.Device{
		DeactivatedAt: time.Now(),
	})
	return err
}

func (s *DeviceStore) DeleteDeviceByExternalID(externalDeviceID string, externalServiceID int) error {
	_, err := s.db.Table("device").
		Where("external_device_id = ? AND external_service_id = ?", externalDeviceID, externalServiceID).
		Delete(&model.Device{})
	return err
}

func (s *DeviceStore) GetAllUsers() ([]dot.ID, error) {
	var userIDs []int64
	x := s.db.SQL(`SELECT DISTINCT user_id FROM device`)
	sql, args, err := x.Build()
	if err != nil {
		return nil, err
	}

	sql2 := fmt.Sprintf(
		"SELECT array_agg(user_id) FROM (%v) AS s",
		sql,
	)
	err = s.db.QueryRow(sql2, args...).Scan((*pq.Int64Array)(&userIDs))
	return util.Int64ToIDs(userIDs), err
}
