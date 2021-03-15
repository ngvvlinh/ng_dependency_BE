package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	identitymodel "o.o/backend/com/main/identity/model"
	fulfillmentmodel "o.o/backend/com/main/shipping/model"
	ticketmodel "o.o/backend/com/supporting/ticket/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/scripts/once/17_migration_ticket/config"
	"o.o/backend/scripts/once/17_migration_ticket/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll    = l.New()
	cfg   config.Config
	newDB *cmsql.Database
	oldDB *cmsql.Database
)

const (
	layout                = "2006-01-02 15:04:05"
	layoutWithMilliSecond = "2006-01-02 15:04:05.123"
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	if newDB, err = cmsql.Connect(cfg.NewDatabase); err != nil {
		ll.Fatal("Error while loading new database", l.Error(err))
	}

	if oldDB, err = cmsql.Connect(cfg.OldDatabase); err != nil {
		ll.Fatal("Error while loading old database", l.Error(err))
	}

	var (
		oldTickets        []*model.Ticket
		oldTicketComments []*model.TicketComment
	)

	mapFFm := make(map[string]*fulfillmentmodel.Fulfillment)     // key: ffmCode
	mapTicket := make(map[string]*ticketmodel.Ticket)            // key: oldTicketID value: newTicket
	mapTicketComments := make(map[string][]*model.TicketComment) // key: ticketID value: ticketComments
	mapCategory := map[string]dot.ID{
		"Thay đổi Tên KH": 1143772113554218820,
		"Giục lấy hàng":   1143772035038117662,
		"Thay đổi SDT":    1143772113550994045,
		"Yêu cầu khác":    1143772113555730296,
		"Thay đổi COD":    1143772035039984275,
		"Giục giao hàng":  1143772035039707243,
		"":                1143772113555730296,
	}
	mapVtigerAccountAndUserID := make(map[string]dot.ID) // key: vtigerAccountID value: userID
	mapUser := make(map[string]dot.ID)
	mapAccount := make(map[string]dot.ID)

	// load user from oldDB (crm)
	{
		vtigerAccounts, err := scanVtigerAccounts()
		if err != nil {
			ll.Fatal("Error when scan users from crm", l.Error(err))
		}

		var vtigerAccountEmail []string
		mapVtigerAccount := make(map[string]string) // key: email1 value: id
		for _, vtigerAccount := range vtigerAccounts {
			vtigerAccountEmail = append(vtigerAccountEmail, vtigerAccount.Email1)
			mapVtigerAccount[vtigerAccount.Email1] = vtigerAccount.ID
		}

		users, err := scanUsersByEmails(vtigerAccountEmail)
		if err != nil {
			ll.Fatal("Error when scan users by emails", l.Error(err))
		}
		for _, user := range users {
			if vtigerAccount, ok := mapVtigerAccount[user.Email]; ok {
				mapVtigerAccountAndUserID[vtigerAccount] = user.ID
			}
		}
	}

	// load old tickets
	fromOffset := uint64(0)
	for {

		tickets, err := scanTickets(fromOffset)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}

		if len(tickets) == 0 {
			break
		}

		oldTickets = append(oldTickets, tickets...)
		fromOffset += uint64(len(tickets))

		var ffmCodes []string

		for _, ticket := range tickets {
			ffmCodes = append(ffmCodes, ticket.FfmCode)
		}

		ffms, err := scanFulfillmentsByShippingCodes(ffmCodes)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}

		for _, ffm := range ffms {
			mapFFm[ffm.ShippingCode] = ffm
		}
	}

	// load old ticketComments
	fromOffset = uint64(0)
	for {
		ticketComments, err := scanTicketComments(fromOffset)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}

		if len(ticketComments) == 0 {
			break
		}

		oldTicketComments = append(oldTicketComments, ticketComments...)
		fromOffset += uint64(len(ticketComments))

		var (
			accountIDs, userIDs []string
		)
		for _, ticketComment := range ticketComments {
			mapTicketComments[ticketComment.TicketID] = append(mapTicketComments[ticketComment.TicketID], ticketComment)
			if ticketComment.CreatorID == "GHN" {
				continue
			}

			userID := strings.Split(ticketComment.CreatorID, "-")[0]
			accountID := strings.Split(ticketComment.CreatorID, "-")[1]

			accountIDs = append(accountIDs, accountID)
			userIDs = append(userIDs, userID)
		}

		if len(accountIDs) != 0 {
			accounts, err := scanAccountsByIDs(accountIDs)
			if err != nil {
				ll.Fatal("Error when scan accounts", l.Error(err))
			}
			for _, account := range accounts {
				mapAccount[account.ID.String()] = account.ID
			}
		}

		if len(userIDs) != 0 {
			users, err := scanUsersByIDs(userIDs)
			if err != nil {
				ll.Fatal("Error when scan users", l.Error(err))
			}
			for _, user := range users {
				mapUser[user.ID.String()] = user.ID
			}
		}

	}

	// create tickets
	ticketsCreated := 0
	for _, oldTicket := range oldTickets {
		ffm := mapFFm[oldTicket.FfmCode]
		if ffm == nil {
			ll.Error(fmt.Sprintf("ticket %v: ffm_code %v not found", oldTicket.ID, oldTicket.FfmCode))
			continue
		}
		if oldTicket.AssignedUserID != "" && mapVtigerAccountAndUserID[oldTicket.AssignedUserID] == 0 {
			ll.Error(fmt.Sprintf("ticket %v: assignedUserID %v not found", oldTicket.ID, oldTicket.AssignedUserID))
			continue
		}

		ID := cm.NewID()
		newTicket := &ticketmodel.Ticket{
			ID:              ID,
			Code:            gencode.GenerateCodeWithChar("TK", time.Now()),
			AssignedUserIDs: []dot.ID{mapVtigerAccountAndUserID[oldTicket.AssignedUserID]},
			AccountID:       ffm.ShopID,
			LabelIDs:        []dot.ID{mapCategory[oldTicket.TicketCategories]},
			ExternalID:      oldTicket.ExternalID,
			Title:           oldTicket.TicketTitle,
			Description:     oldTicket.Description,
			Note:            oldTicket.Note,
			RefID:           ffm.ID,
			RefType:         ticket_ref_type.FFM,
			Source:          ticket_source.System,
			State:           0,
			Status:          0,
			CreatedSource:   account_type.Etop,
			UpdatedBy:       mapVtigerAccountAndUserID[oldTicket.ModifiedBy],
			ConfirmedBy:     mapVtigerAccountAndUserID[oldTicket.ConfirmedBy],
			ClosedBy:        mapVtigerAccountAndUserID[oldTicket.ClosedBy],
			CreatedAt:       oldTicket.CreatedAt,
			UpdatedAt:       oldTicket.ModifiedAt,
			ConfirmedAt:     oldTicket.ConfirmedAt,
			ClosedAt:        oldTicket.ClosedAt,
			ConnectionID:    ffm.ConnectionID,
		}
		if oldTicket.EtopID != "" {
			createdBy, err := dot.ParseID(oldTicket.EtopID)
			if err != nil {
				ll.Info(fmt.Sprintf("ticket %v: etop_id %v parse ID error", oldTicket.ID, oldTicket.EtopID))
				continue
			}
			newTicket.CreatedBy = createdBy
		}

		switch oldTicket.TicketStatus {
		case "Open":
			newTicket.State = ticket_state.New
		case "In Progress":
			newTicket.State = ticket_state.Processing
		case "Confirmed":
			newTicket.State = ticket_state.Received
		case "Closed":
			if oldTicket.SubStatus == "success" || oldTicket.SubStatus == "" {
				newTicket.State = ticket_state.Success
			} else {
				newTicket.State = ticket_state.Fail
			}
		}
		newTicket.Status = newTicket.State.ToStatus5()

		mapTicket[oldTicket.ID] = newTicket
		if err := newDB.ShouldInsert(newTicket); err != nil {
			ll.Error(fmt.Sprintf("ticket %v: error when creating ticket", oldTicket.ID))
			continue
		}

		ll.Info(fmt.Sprintf("old ticket: %v | new ticket: %v | shop: %v", oldTicket.ID, newTicket.ID, newTicket.AccountID))

		ticketsCreated += 1
	}
	ll.Info(fmt.Sprintf("ticket success: %v/%v", ticketsCreated, len(oldTickets)))

	ticketCommentCreated := 0
	for _, oldTicketComment := range oldTicketComments {
		if oldTicketComment.CreatorID == "GHN" {
			continue
		}

		userID := strings.Split(oldTicketComment.CreatorID, "-")[0]
		accountID := strings.Split(oldTicketComment.CreatorID, "-")[1]

		if mapTicket[oldTicketComment.TicketID] == nil {
			continue
		}

		if mapAccount[accountID] == 0 {
			ll.Error(fmt.Sprintf("ticketcomment %v: accountID %v not found", oldTicketComment.ID, accountID))
			continue
		}

		if mapUser[userID] == 0 {
			ll.Error(fmt.Sprintf("ticketcomment %v: userID %v not found", oldTicketComment.ID, userID))
			continue
		}
		newTicketComment := &ticketmodel.TicketComment{
			ID:                cm.NewID(),
			TicketID:          mapTicket[oldTicketComment.TicketID].ID,
			AccountID:         mapTicket[oldTicketComment.TicketID].AccountID,
			CreatedBy:         mapUser[userID],
			ExternalCreatedAt: oldTicketComment.CreatedAt,
			Message:           oldTicketComment.Content,
			CreatedAt:         oldTicketComment.CreatedAt,
			UpdatedAt:         oldTicketComment.CreatedAt,
		}
		if accountID == "101" {
			newTicketComment.CreatedSource = account_type.Etop
		} else {
			newTicketComment.CreatedSource = account_type.Shop
		}

		_ = newTicketComment

		if err := newDB.ShouldInsert(newTicketComment); err != nil {
			ll.Info(fmt.Sprintf("ticket %v: error when creating ticket", oldTicketComment.ID))
			continue
		}

		ll.Info(fmt.Sprintf("ticket comment %v: ticket_id %v accountID %v", newTicketComment.ID, newTicketComment.TicketID, mapAccount[accountID]))
		ticketCommentCreated += 1
	}

	ll.Info(fmt.Sprintf("ticket_comment success: %v/%v", ticketCommentCreated, len(oldTicketComments)))
}

func scanTickets(fromOffset uint64) (tickets model.Tickets, err error) {
	err = oldDB.
		From("ticket").
		Where("created_at >= '2021-01-01 00:00:01'").
		OrderBy("id").
		Offset(fromOffset).
		Limit(1000).
		Find(&tickets)

	return
}

func scanTicketComments(fromOffset uint64) (ticketComments model.TicketComments, err error) {
	err = oldDB.
		From("ticket_comment").
		Where("creator_id <> 'GHN'").
		Where("created_at >= '2021-01-01 00:00:01'").
		OrderBy("id").
		Offset(fromOffset).
		Limit(1000).
		Find(&ticketComments)

	return
}

func scanFulfillmentsByShippingCodes(shippingCodes []string) (ffms fulfillmentmodel.Fulfillments, err error) {
	err = newDB.
		From("fulfillment").
		In("shipping_code", shippingCodes).
		Find(&ffms)

	return
}

func scanVtigerAccounts() (vtigerAccounts model.VtigerAccounts, err error) {
	err = oldDB.
		From("vtiger_account").
		Find(&vtigerAccounts)

	return
}

func scanUsersByEmails(emails []string) (users identitymodel.Users, err error) {
	err = newDB.
		From("user").
		In("email", emails).
		Find(&users)

	return
}

func scanUsersByIDs(ids []string) (users identitymodel.Users, err error) {
	err = newDB.
		From("user").
		In("id", ids).
		Find(&users)

	return
}

func scanAccountsByIDs(ids []string) (accounts identitymodel.Accounts, err error) {
	err = newDB.
		From("account").
		In("id", ids).
		Find(&accounts)

	return
}
