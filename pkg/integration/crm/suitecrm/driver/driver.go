package driver

import (
	"context"
	"strings"

	"o.o/api/main/contact"
	"o.o/api/main/shipping"
	"o.o/api/supporting/ticket"
	carriertypes "o.o/backend/com/supporting/ticket/provider/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	suitecrmclient "o.o/backend/pkg/integration/crm/suitecrm/client"
	"o.o/common/l"
)

var ll = l.New()
var _ carriertypes.TicketProvider = &SuiteCRMTicketDriver{}

type SuiteCRMTicketDriver struct {
	client        *suitecrmclient.Client
	shippingQuery shipping.QueryBus
	contactQuery  contact.QueryBus
}

func New(
	env string,
	cfg *suitecrmclient.SuiteCRMCfg,
	shippingQuery shipping.QueryBus,
	contactQuery contact.QueryBus,
) *SuiteCRMTicketDriver {
	return &SuiteCRMTicketDriver{
		client:        suitecrmclient.New(env, cfg),
		shippingQuery: shippingQuery,
		contactQuery:  contactQuery,
	}
}

func (s SuiteCRMTicketDriver) Ping(ctx context.Context) error {
	return nil
}

func (s SuiteCRMTicketDriver) CreateTicket(ctx context.Context, ticket *ticket.Ticket) (*ticket.Ticket, error) {
	getContactQuery := &contact.GetContactByIDQuery{
		ID:     ticket.RefID,
		ShopID: ticket.AccountID,
	}
	if err := s.contactQuery.Dispatch(ctx, getContactQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Không tìm thấy thông tin liên lạc")
		}
		return nil, err
	}
	contact := getContactQuery.Result

	insertCaseReq := &suitecrmclient.InsertCaseRequest{
		PhoneMobile: contact.Phone,
		Subject:     ticket.Title,
		Description: ticket.Description,
		RefType:     "topship", // hardcode
		RefID:       ticket.ID.String(),
		RefCompany:  "topship", // hardcode
	}
	insertCaseResp, err := s.client.InsertCase(ctx, insertCaseReq)
	if err != nil {
		return nil, err
	}

	// handle create fail
	messageResp := insertCaseResp.Message.String()
	if strings.Contains(messageResp, "cannot create cases") {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không tạo ticket thành công từ SuiteCRM. Chúng tôi đang liên hệ với SuiteCRM để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	return ticket, nil
}

func (s SuiteCRMTicketDriver) CreateTicketComment(ctx context.Context, ticketComment *ticket.TicketComment) (*ticket.TicketComment, error) {
	return ticketComment, nil
}
