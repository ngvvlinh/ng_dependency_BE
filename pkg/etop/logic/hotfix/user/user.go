package user

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

type UserService struct {
	etelecomAggr  etelecom.CommandBus
	etelecomQuery etelecom.QueryBus
	dbEtelecom    *cmsql.Database
	dbMain        *cmsql.Database
	identityAggr  identity.CommandBus
	identityQuery identity.QueryBus
}

func New(etelecomA etelecom.CommandBus, etelecomQ etelecom.QueryBus, dbEtelecom com.EtelecomDB, dbMain com.MainDB, identityAggr identity.CommandBus, identityQuery identity.QueryBus) *UserService {
	return &UserService{
		etelecomAggr:  etelecomA,
		etelecomQuery: etelecomQ,
		dbEtelecom:    dbEtelecom,
		dbMain:        dbMain,
		identityAggr:  identityAggr,
		identityQuery: identityQuery,
	}
}

type CreateAndActiveHotlineArgs struct {
	Hotline      string
	OwnerID      dot.ID
	ConnectionID dot.ID
	TenantID     dot.ID
}

type Line struct {
	ID          string
	Name        string
	Email       string
	Phone       string
	CompanyName string
	Hotlines    []string
}

func (s *UserService) HandleImportUser(c *httpx.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
	// continue
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}

	file, err := files[0].Open()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not read file")
	}
	defer file.Close()

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file.")
	}
	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file")
	}

	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file")
	}
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung")
	}

	ctx := bus.Ctx()
	success := 0
	total := len(rows) - 1
	for i, row := range rows {
		if i == 0 {
			continue
		}
		line, _err := s.parseRow(ctx, row)
		if _err != nil {
			return _err
		}
		owner, err := s.checkAndRegisterUser(ctx, line)
		if err != nil {
			return nil
		}

		// Get and check tenant, if not exist, create tenant
		createTenantArgs := &etelecom.CreateTenantArgs{
			OwnerID:      owner.ID,
			ConnectionID: connectioning.DefaultDirectPortsipConnectionID,
		}
		tenant, err := s.checkAndCreateTenant(ctx, createTenantArgs)
		if err != nil {
			return err
		}

		if tenant.Status.Valid && tenant.Status.Enum != status3.P {
			cmd := &etelecom.ActivateTenantCommand{
				OwnerID:      owner.ID,
				TenantID:     tenant.ID,
				ConnectionID: connectioning.DefaultDirectPortsipConnectionID,
			}
			if err = s.etelecomAggr.Dispatch(ctx, cmd); err != nil {
				return err
			}
		}

		for _, hotline := range line.Hotlines {
			args := &CreateAndActiveHotlineArgs{
				Hotline:      hotline,
				OwnerID:      owner.ID,
				ConnectionID: tenant.ConnectionID,
				TenantID:     tenant.ID,
			}
			if err = s.createAndActiveHotline(ctx, args); err != nil {
				return err
			}
		}
		success += 1
	}

	c.SetResult(map[string]interface{}{
		"code":    "ok",
		"total":   total,
		"success": success,
	})
	return nil
}

func (s *UserService) createAndActiveHotline(ctx context.Context, args *CreateAndActiveHotlineArgs) error {
	getHotlineQuery := &etelecom.GetHotlineByHotlineNumberQuery{
		Hotline: args.Hotline,
		OwnerID: args.OwnerID,
	}
	err := s.etelecomQuery.Dispatch(ctx, getHotlineQuery)
	var hotline *etelecom.Hotline
	switch cm.ErrorCode(err) {
	case cm.NoError:
		hotline = getHotlineQuery.Result
	case cm.NotFound:
		createHotlineCmd := &etelecom.CreateHotlineCommand{
			OwnerID:      args.OwnerID,
			Hotline:      args.Hotline,
			ConnectionID: args.ConnectionID,
			Status:       status3.Z,
		}
		if err := s.etelecomAggr.Dispatch(ctx, createHotlineCmd); err != nil {
			return err
		}
		hotline = createHotlineCmd.Result
	default:
		return err
	}
	if hotline.Status == status3.P {
		return nil
	}
	activeHotlineForTenant := &etelecom.ActiveHotlineForTenantCommand{
		HotlineID: hotline.ID,
		OwnerID:   hotline.OwnerID,
		TenantID:  args.TenantID,
	}
	if err = s.etelecomAggr.Dispatch(ctx, activeHotlineForTenant); err != nil {
		return err
	}
	return nil
}
func (s *UserService) checkAndRegisterUser(ctx context.Context, line *Line) (*identity.User, error) {
	getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
		Phone: line.Phone,
	}
	err := s.identityQuery.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		// continue
	case cm.NotFound:
		cmd := &identity.RegisterSimplifyCommand{
			Phone:               line.Phone,
			Password:            line.Phone,
			FullName:            line.Name,
			Email:               line.Email,
			CompanyName:         line.CompanyName,
			IsCreateDefaultShop: true,
		}
		if err := s.identityAggr.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	// Get owner
	getUserQuery = &identity.GetUserByPhoneOrEmailQuery{
		Phone: line.Phone,
	}
	if err := s.identityQuery.Dispatch(ctx, getUserQuery); err != nil {
		return nil, err
	}
	return getUserQuery.Result, nil
}

func (s *UserService) checkAndCreateTenant(ctx context.Context, args *etelecom.CreateTenantArgs) (tenant *etelecom.Tenant, err error) {
	// Get and check tenant, if not exist, create tenant
	getTenantQuery := &etelecom.GetTenantByConnectionQuery{
		OwnerID:      args.OwnerID,
		ConnectionID: args.ConnectionID,
	}
	err = s.etelecomQuery.Dispatch(ctx, getTenantQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		tenant = getTenantQuery.Result
	case cm.NotFound:
		createTenantCmd := &etelecom.CreateTenantCommand{
			OwnerID:      args.OwnerID,
			AccountID:    args.AccountID,
			ConnectionID: args.ConnectionID,
		}
		if err = s.etelecomAggr.Dispatch(ctx, createTenantCmd); err != nil {
			return nil, err
		}
		tenant = createTenantCmd.Result
	default:
		return nil, err
	}
	return
}

func (s *UserService) parseRow(ctx context.Context, row []string) (*Line, error) {
	phone, ok := validate.NormalizePhone(row[3])
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phone does not valid")
	}
	strHotlines := strings.Split(row[5], ",")
	var hotlines = []string{}
	for _, hotline := range strHotlines {
		hotlines = append(hotlines, strings.ReplaceAll(hotline, " ", ""))
	}
	return &Line{
		ID:          row[0],
		Name:        row[1],
		Email:       row[2],
		Phone:       phone.String(),
		CompanyName: row[4],
		Hotlines:    hotlines,
	}, nil
}
