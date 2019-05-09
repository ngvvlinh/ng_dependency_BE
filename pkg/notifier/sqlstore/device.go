package sqlstore

import (
	"fmt"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/notifier/model"
	"github.com/lib/pq"
)

type DeviceStore struct {
	db cmsql.Database
}

func NewDeviceStore(db cmsql.Database) *DeviceStore {
	return &DeviceStore{
		db: db,
	}
}

func (s *DeviceStore) CreateDevice(args *model.CreateDeviceArgs) (*model.Device, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if args.DeviceID == "" || args.DeviceName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Device ID or Device Name")
	}
	if args.ExternalDeviceID == "" {
		return nil, cm.Errorf(cm.Internal, nil, "Missing External Device ID")
	}
	var dbDevice = new(model.Device)
	ok, err := s.db.Table("device").Where("device_id = ? AND account_id = ?", args.DeviceID, args.AccountID).Get(dbDevice)
	if err != nil {
		return nil, err
	}
	device := &model.Device{
		AccountID:        args.AccountID,
		DeviceID:         args.DeviceID,
		DeviceName:       args.DeviceName,
		ExternalDeviceID: args.ExternalDeviceID,
	}
	var id int64
	if ok && dbDevice.ID != 0 {
		// update this device
		id = dbDevice.ID
		if err := s.db.Table("device").Where("id = ?", id).ShouldUpdate(device); err != nil {
			return nil, err
		}
	} else {
		id = cm.NewID()
		device.ID = id
		// Use Onesignal by default
		device.ExternalServiceID = model.ExternalServiceOneSignalID
		if err := s.db.Table("device").ShouldInsert(device); err != nil {
			return nil, err
		}
	}
	res, err := s.GetDevice(&model.GetDeviceArgs{
		AccountID: args.AccountID,
		ID:        id,
	})
	return res, err
}

func (s *DeviceStore) GetDevice(args *model.GetDeviceArgs) (*model.Device, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	var device = new(model.Device)
	err := s.db.Table("device").Where("account_id = ? AND id = ?", args.AccountID, args.ID).ShouldGet(device)
	return device, err
}

func (s *DeviceStore) GetDevices(args *model.GetDevicesArgs) ([]*model.Device, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if args.ExternalServiceID == 0 {
		args.ExternalServiceID = model.ExternalServiceOneSignalID
	}

	var res []*model.Device
	if err := s.db.Table("device").Where("account_id = ? AND external_service_id = ?", args.AccountID, args.ExternalServiceID).Find((*model.Devices)(&res)); err != nil {
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
	if device.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if device.DeviceID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing Device ID")
	}
	err := s.db.Table("device").Where("device_id = ? AND account_id = ?", device.DeviceID, device.AccountID).ShouldDelete(&model.Device{})
	return err
}

func (s *DeviceStore) GetAllAccounts() ([]int64, error) {
	var accountIDs []int64
	x := s.db.SQL(`SELECT DISTINCT account_id FROM device`)
	sql, args, err := x.Build()
	if err != nil {
		return nil, err
	}

	sql2 := fmt.Sprintf(
		"SELECT array_agg(account_id) FROM (%v) AS s",
		sql,
	)
	err = s.db.QueryRow(sql2, args...).Scan((*pq.Int64Array)(&accountIDs))
	return accountIDs, err
}
