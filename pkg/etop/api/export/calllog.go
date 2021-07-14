package export

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"o.o/api/etelecom"
	"o.o/api/etelecom/call_direction"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	apishop "o.o/api/top/int/shop"
	etelecommodel "o.o/backend/com/etelecom/model"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/capi/dot"
)

func (s *Service) ExportCallLogs(ctx context.Context,
	id string, shop *identitymodel.Shop,
	exportOpts ExportOption, output io.Writer,
	result chan<- *apishop.ExportStatusItem,
	total int, rows RowsInterface, opts core.Opts,
) (_err error) {
	var count, countError int
	maxErrors := total / 100
	if maxErrors > BaseRowsErrors {
		maxErrors = BaseRowsErrors
	}

	makeProgress := func() *apishop.ExportStatusItem {
		return &apishop.ExportStatusItem{
			Id:            id,
			ProgressMax:   total,
			ProgressValue: count,
			ProgressError: countError,
		}
	}
	handleError := func(err error) bool {
		switch err {
		case nil:
			return true
		case ErrAborted:
			return false
		default:
			countError++
			statusItem := makeProgress()
			statusItem.Error = cmapi.PbError(err)
			if countError > maxErrors {
				err = cm.Errorf(cm.Aborted, nil, "Quá nhiều lỗi xảy ra")
				statusItem.Error = cmapi.PbError(err)
				result <- statusItem
			}
			return false
		}
	}
	defer func() {
		handleError(_err)
		if _err == nil {
			result <- makeProgress()
		}
		close(result)
	}()

	hotlines, err := s.getAllHotlines(ctx, shop.OwnerID)
	if err != nil {
		return err
	}
	users, err := s.getAllUsers(ctx, shop.ID)
	if err != nil {
		return err
	}

	if _, err = WriteBOM(output); err != nil {
		return err
	}

	csvWriter := csv.NewWriter(output)
	csvWriter.Comma = exportOpts.Delimiter
	defer func() {
		csvWriter.Flush()
		if err = csvWriter.Error(); _err == nil && err != nil {
			_err = err
		}
	}()

	tableWrite := NewTableWriter(csvWriter, Noempty)
	var callLog etelecommodel.CallLog
	if err = buildTableCallLog(csvWriter, tableWrite, exportOpts, &count, &callLog, hotlines, users); err != nil {
		return err
	}

	result <- makeProgress()
	lastTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return ErrAborted
		default:
			if !rows.Next() {
				return rows.Err()
			}
			args := callLog.SQLScanArgs(opts)
			if err = rows.Scan(args...); err != nil {
				if !handleError(err) {
					return ErrAborted
				}
			}

			// increase count and line for each line written to csv file
			count++
			if err = tableWrite.WriteRow(); err != nil {
				if !handleError(err) {
					return ErrAborted
				}
			}

			if now := time.Now(); now.Sub(lastTime) > 100*time.Millisecond {
				result <- makeProgress()
				lastTime = now
			}
		}
	}
}

func (s *Service) getAllHotlines(ctx context.Context, ownerID dot.ID) (map[dot.ID]*etelecom.Hotline, error) {
	var res = make(map[dot.ID]*etelecom.Hotline)
	query := &etelecom.ListHotlinesQuery{
		OwnerID:      ownerID,
		ConnectionID: connectioning.DefaultDirectPortsipConnectionID,
	}
	if err := s.etelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	for _, hotline := range query.Result {
		res[hotline.ID] = hotline
	}
	return res, nil
}

func (s *Service) getAllUsers(ctx context.Context, shopID dot.ID) (mapUsers map[dot.ID]*identity.User, err error) {
	mapUsers = make(map[dot.ID]*identity.User)
	query := &identity.ListAccountUsersQuery{
		AccountID: shopID,
	}
	if err = s.identityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	userIDs := []dot.ID{}
	for _, accUser := range query.Result.AccountUsers {
		userIDs = append(userIDs, accUser.UserID)
	}

	queryUser := &identity.ListUsersByIDsAndNameNormQuery{
		IDs: userIDs,
	}
	if err = s.identityQuery.Dispatch(ctx, queryUser); err != nil {
		return nil, err
	}
	for _, user := range queryUser.Result {
		mapUsers[user.ID] = user
	}
	return
}

func buildTableCallLog(csvWriter *csv.Writer,
	w *TableWriter, exportOpts ExportOption,
	count *int, callLog *etelecommodel.CallLog,
	hotlines map[dot.ID]*etelecom.Hotline,
	users map[dot.ID]*identity.User,
) error {
	var formatAsText = func(s string) string { return s }
	if exportOpts.ExcelMode {
		formatAsText = FormatAsTextForExcel
	}

	w.AddColumn("STT", func() string {
		return strconv.Itoa(*count)
	})
	w.AddColumn("Số điện thoại", func() string {
		// số điện thoại khách hàng
		_, phoneNumber := getExtNumberAndPhoneNumber(callLog)
		return formatAsText(phoneNumber)
	})
	w.AddColumn("Hotline", func() string {
		hotline, ok := hotlines[callLog.HotlineID]
		if !ok {
			return ""
		}
		return formatAsText(hotline.Name + " - " + hotline.Hotline)
	})
	w.AddColumn("Máy nhánh", func() string {
		extNumber, _ := getExtNumberAndPhoneNumber(callLog)
		return formatAsText(extNumber)
	})
	w.AddColumn("Tên người dùng", func() string {
		// Tên nhân viên sử dụng máy nhánh
		user, ok := users[callLog.UserID]
		if !ok {
			return ""
		}
		return user.FullName
	})
	w.AddColumn("Hướng", func() string {
		return callLog.Direction.GetLabelRefName()
	})
	w.AddColumn("Trạng thái", func() string {
		return callLog.CallState.GetLabelRefName()
	})
	w.AddColumn("Thời điểm bắt đầu", func() string {
		// format: 2021-05-11 7:51:27
		return FormatDateTime(callLog.StartedAt)
	})
	w.AddColumn("Thời điểm kết thúc", func() string {
		// format: 2021-05-11 7:51:27
		return FormatDateTime(callLog.EndedAt)
	})
	w.AddColumn("Thời gian đàm thoại", func() string {
		// format: 0:01:52
		return FormatDuration(time.Duration(callLog.Duration) * time.Second)
	})
	w.AddColumn("File ghi âm", func() string {
		recordingURLs := callLog.AudioURLs
		if recordingURLs == nil || len(recordingURLs) == 0 {
			return ""
		}
		return recordingURLs[0]
	})

	return w.WriteHeader()
}

func getExtNumberAndPhoneNumber(callLog *etelecommodel.CallLog) (extNumber, phoneNumber string) {
	switch callLog.Direction {
	case call_direction.In, call_direction.ExtIn:
		phoneNumber = callLog.Caller
		extNumber = callLog.Callee
	case call_direction.Out, call_direction.ExtOut:
		phoneNumber = callLog.Callee
		extNumber = callLog.Caller
	default:
	}
	return
}
