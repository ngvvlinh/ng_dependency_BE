package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"o.o/api/top/types/etc/ticket/ticket_type"
	ticketmodel "o.o/backend/com/supporting/ticket/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/scripts/once/18_migration_ticket_label/config"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	DB  *cmsql.Database
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	ll.S.Info("config", cfg)

	if DB, err = cmsql.Connect(cfg.Database); err != nil {
		ll.Fatal("Error while loading new database", l.Error(err))
	}

	var shopIDs []int64

	{
		mShopIDsHaveDefaultTicketLabels := make(map[int64]bool)

		offset := uint64(0)
		for {
			_shopIDs, err := loadShopIDsHaveDefaultTicketLabels(offset)
			if err != nil {
				log.Fatal(err)
			}

			for _, _shopID := range _shopIDs {
				mShopIDsHaveDefaultTicketLabels[_shopID] = true
			}

			if len(_shopIDs) < 1000 {
				break
			}
			offset += 1000
		}

		fmt.Println(len(mShopIDsHaveDefaultTicketLabels))
		fromID := int64(0)
		for {
			_shopIDs, err := loadShopIDs(fromID)
			if err != nil {
				log.Fatal(err)
			}

			if len(_shopIDs) == 0 {
				break
			}

			fromID = _shopIDs[len(_shopIDs)-1]

			for _, _shopID := range _shopIDs {
				if mShopIDsHaveDefaultTicketLabels[_shopID] {
					continue
				}
				shopIDs = append(shopIDs, _shopID)
			}
			fmt.Println(len(shopIDs))
		}
	}

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)

	success := 0

	wg.Add(len(shopIDs))

	for _, shopID := range shopIDs {
		ids := []dot.ID{cm.NewID(), cm.NewID(), cm.NewID()}
		go func(_shopID int64, _ids []dot.ID) {
			defer func() {
				wg.Done()
			}()

			newTickets := []*ticketmodel.TicketLabel{
				{
					ID:        _ids[0],
					ShopID:    dot.ID(_shopID),
					Type:      ticket_type.Internal,
					Name:      "Tư vấn",
					Code:      "support",
					Color:     "#1abc9c",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        _ids[1],
					ShopID:    dot.ID(_shopID),
					Type:      ticket_type.Internal,
					Name:      "Góp ý",
					Code:      "feedback",
					Color:     "#9b59b6",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        _ids[2],
					ShopID:    dot.ID(_shopID),
					Type:      ticket_type.Internal,
					Name:      "Khiếu nại",
					Code:      "complain",
					Color:     "#e74c3c",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			if err = DB.InTransaction(bus.Ctx(), func(tx cmsql.QueryInterface) error {
				for _, newTicket := range newTickets {
					if _, err = tx.Insert(newTicket); err != nil {
						return err
					}
				}

				return nil
			}); err == nil {
				mu.Lock()
				success += 1
				mu.Unlock()
			}
		}(shopID, ids)
	}

	wg.Wait()

	fmt.Println(fmt.Sprintf("success: %v/%v", success, len(shopIDs)))
}

func loadShopIDs(fromID int64) (shopIDs []int64, err error) {
	rows, err := DB.
		Select("id").
		From("shop").
		Where("id > ?", fromID).
		OrderBy("id").
		Query()
	if err != nil {
		return nil, err
	}

	var shopID int64
	for rows.Next() {
		err := rows.Scan(&shopID)
		if err != nil {
			return nil, err
		}
		shopIDs = append(shopIDs, shopID)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return shopIDs, nil
}

func loadShopIDsHaveDefaultTicketLabels(offset uint64) (shopIDs []int64, err error) {
	rows, err := DB.
		Select("shop_id").
		From("ticket_label").
		Where("color = '#1abc9c' and type = ?", ticket_type.Internal).
		OrderBy("id").
		Limit(1000).
		Offset(offset).
		Query()
	if err != nil {
		return nil, err
	}

	var shopID int64
	for rows.Next() {
		err := rows.Scan(&shopID)
		if err != nil {
			return nil, err
		}
		shopIDs = append(shopIDs, shopID)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return shopIDs, nil
}
