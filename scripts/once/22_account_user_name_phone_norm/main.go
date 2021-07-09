package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	model2 "o.o/backend/com/etelecom/model"
	"o.o/backend/com/main/identity/model"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/scripts/once/22_account_user_name_phone_norm/config"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll         = l.New()
	cfg        config.Config
	DBMain     *cmsql.Database
	DBEtelecom *cmsql.Database
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	if DBMain, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while loading new database", l.Error(err))
	}
	if DBEtelecom, err = cmsql.Connect(cfg.PostgresEtelecom); err != nil {
		ll.Fatal("Error while loading new database", l.Error(err))
	}

	mapUser := make(map[dot.ID]*model.User)
	mapExtension := make(map[string]*model2.Extension)
	mapAccountUser := make(map[string]*model.AccountUser)
	count, errCount, updatedCount := 0, 0, 0
	{
		var fromCreatedAt time.Time
		for {
			// Get Account User and map account user by created_at
			accountUsers, err := scanAccountUsers(fromCreatedAt)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			count += len(accountUsers)
			if len(accountUsers) == 0 {
				break
			}

			userIDs := make([]dot.ID, 0, len(accountUsers))
			accountIDs := make([]dot.ID, 0, len(accountUsers))
			for _, accountUser := range accountUsers {
				accountUserID := fmt.Sprintf("%v_%v", accountUser.AccountID.String(), accountUser.UserID.String())
				mapAccountUser[accountUserID] = accountUser
				if accountUser.UserID != 0 {
					userIDs = append(userIDs, accountUser.UserID)
				}
				if accountUser.AccountID != 0 {
					accountIDs = append(accountIDs, accountUser.AccountID)
				}
			}

			// Get users by IDs and map user
			users, err := getUsersByIDs(userIDs)
			if err != nil {
				ll.Fatal("Error in GetUsersByIDs", l.Error(err))
			}
			for _, user := range users {
				if _, ok := mapUser[user.ID]; !ok {
					mapUser[user.ID] = user
				}
			}

			// Get users by IDs and map user
			extensions, err := getExtensionsByAccountIDsAndUserIDs(accountIDs, userIDs)
			if err != nil {
				ll.Fatal("Error in GetExtensionsByAccountIDsAndUserIDs", l.Error(err))
			}
			for _, ext := range extensions {
				accountUserID := fmt.Sprintf("%v_%v", ext.AccountID.String(), ext.UserID.String())
				if _, ok := mapExtension[accountUserID]; !ok {
					mapExtension[accountUserID] = ext
				}
			}
			fromCreatedAt = accountUsers[len(accountUsers)-1].CreatedAt
		}
	}

	{

		var (
			mu sync.Mutex
			wg sync.WaitGroup
		)
		maxGoroutines := 8
		ch := make(chan string, maxGoroutines)
		wg.Add(count)
		for accountUserID, accountUser := range mapAccountUser {
			ch <- accountUserID
			var phone, fullNameNorm, phoneNorm, extensionNumberNorm string
			if _, ok := mapUser[accountUser.UserID]; ok {
				user := mapUser[accountUser.UserID]
				phone = user.Phone
				phoneNorm = validate.NormalizeSearchCharacter(phone)
				fullNameNorm = validate.NormalizeSearchCharacter(user.FullName)
			}

			if _, ok := mapExtension[accountUserID]; ok {
				extensionNumberNorm = validate.NormalizeSearchCharacter(mapExtension[accountUserID].ExtensionNumber)
			}

			arrIDs := strings.Split(accountUserID, "_")
			accountID := arrIDs[0]
			userID := arrIDs[1]

			go func(accountID, userID, phone, fullNameNorm, phoneNorm, extensionNumberNorm string) {
				defer func() {
					<-ch
					wg.Done()
				}()
				err = updateAccountUser(accountID, userID, phone, fullNameNorm, phoneNorm, extensionNumberNorm)
				if err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}(accountID, userID, phone, fullNameNorm, phoneNorm, extensionNumberNorm)
		}
		wg.Wait()
		updatedCount = count - errCount
		ll.S.Infof("Update phone, full_name_norm, phone_norm for account_user: updated %v/%v, error %v/%v", updatedCount, count, errCount, count)
	}

}

func scanAccountUsers(fromCreatedAt time.Time) (accountUsers model.AccountUsers, err error) {
	err = DBMain.Where("created_at > ?", fromCreatedAt).
		Where("(phone IS NULL OR full_name_norm IS NULL OR phone_norm IS NULL OR extension_number_norm IS NULL) AND (created_at IS NOT NULL) AND (deleted_at IS NULL)").
		OrderBy("created_at").
		Limit(1000).
		Find(&accountUsers)
	return
}

func getUsersByIDs(ids []dot.ID) (users model.Users, err error) {
	err = DBMain.
		From("user").
		In("id", ids).
		Find(&users)
	return
}

func getExtensionsByAccountIDsAndUserIDs(accountIDs, userIDs []dot.ID) (extensions model2.Extensions, err error) {
	err = DBEtelecom.
		From("extension").
		Where("deleted_at IS NULL").
		In("account_id", accountIDs).
		In("user_id", userIDs).
		Find(&extensions)
	return
}

func updateAccountUser(accountID, userID, phone, fullNameNorm, phoneNorm, extensionNumberNorm string) (err error) {
	accountUser := model.AccountUser{
		FullNameNorm:        fullNameNorm,
		Phone:               phone,
		PhoneNorm:           phoneNorm,
		ExtensionNumberNorm: extensionNumberNorm,
	}

	err = DBMain.
		Table("account_user").
		Where("account_id=? AND user_id=?", accountID, userID).
		ShouldUpdate(&accountUser)
	return err
}
