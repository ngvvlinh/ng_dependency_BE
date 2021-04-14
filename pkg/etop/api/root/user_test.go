package root

/*
func TestCreateLoginResponse(t *testing.T) {
	ctx := context.Background()
	defaultTime, err := time.Parse(time.Stamp, "Jan 1 10:11:12")
	if err != nil {
		panic(err)
	}

	Convey("CreateLoginResponse", t, func() {
		token := "sample_token"
		userID := int64(123456789)
		user := &model.User{
			ID:        userID,
			CreatedAt: defaultTime,
			UpdatedAt: defaultTime,
		}
		user.FullName = "Josie"

		bus.MockHandler(func(query *model.GetUserByIDQuery) error {
			query.Result = user
			return nil
		})
		bus.MockHandler(func(query *model.GetAllAccountRolesQuery) error {
			return nil
		})

		loginResp := &pb.LoginResponse{
			AccessToken: token,
			ExpiresIn:   tokens.DefaultAccessTokenTTL,
			User: &pb.User{
				ID:       userID,
				FullName: "Josie",
			},
			AvailableAccounts: []*pb.LoginAccount{},
		}
		Convey("Login user", func() {
			resp, err := CreateLoginResponse(ctx, nil, token, userID, nil, 0, 0, false, 0)
			So(err, ShouldBeNil)
			So(resp, ShouldDeepEqual, loginResp)
		})

		accounts := []*model.AccountUserExtended{
			{
				Account: &model.Account{
					ID:   cm.NewIDWithTag(model.TagSupplier),
					Name: "Alice",
					Type: model.TypeSupplier,
				},
				AccountUser: &model.AccountUser{},
				User:        &model.User{},
			}, {
				Account: &model.Account{
					ID:   cm.NewIDWithTag(model.TagShop),
					Name: "Johny",
					Type: model.TypeShop,
				},
				AccountUser: &model.AccountUser{},
				User:        &model.User{},
			},
		}
		for _, acc := range accounts {
			acc.User.ID = cm.NewIDWithTag(model.TagUser)
			acc.AccountUser.Name = acc.Account.ID
			acc.AccountUser.UserID = acc.User.ID
		}

		accountsResp := []*pb.LoginAccount{
			{
				ID:   accounts[0].Account.ID,
				Name: "Alice",
				Type: pb.AccountType_supplier,
			},
			{
				ID:   accounts[1].Account.ID,
				Name: "Johny",
				Type: pb.AccountType_shop,
			},
		}
		bus.MockHandler(func(query *model.GetAllAccountRolesQuery) error {
			query.Result = accounts
			return nil
		})
		Convey("Login user with available accounts", func() {
			resp, err := CreateLoginResponse(ctx, nil, token, userID, nil, 0, 0, false, 0)
			So(err, ShouldBeNil)

			loginResp.AvailableAccounts = accountsResp
			So(resp, ShouldDeepEqual, loginResp)
		})

		Convey("Login user with preferred Name (shop)", func() {
			shop := &model.ShopExtended{
				Shop: &model.Shop{
					ID:        accounts[1].Account.ID,
					Name:      "Johny",
					Status:    model.StatusActive,
					CreatedAt: defaultTime,
					UpdatedAt: defaultTime,
				},
			}
			bus.MockHandler(func(query *model.GetShopExtendedQuery) error {
				query.Result = shop
				return nil
			})

			preferAccountID := shop.ID
			resp, err := CreateLoginResponse(ctx, nil, token, userID, nil, preferAccountID, 0, false, 0)
			So(err, ShouldBeNil)

			loginResp.AvailableAccounts = accountsResp
			loginResp.Account = &pb.LoginAccount{
				ID:   shop.ID,
				Name: "Johny",
				Type: pb.AccountType_shop,
			}

			loginResp.Shop = &pb.Shop{
				ID:     shop.ID,
				Name:   "Johny",
				Status: status3.Status_P,
			}
			So(resp, ShouldDeepEqual, loginResp)
		})

		Convey("Login user with preferred Name (supplier)", func() {
			supplier := &model.SupplierExtended{
				Supplier: &model.Supplier{
					ID:        accounts[0].Account.ID,
					Name:      "Alice",
					Status:    model.StatusActive,
					CreatedAt: defaultTime,
					UpdatedAt: defaultTime,
				},
			}
			bus.MockHandler(func(query *model.GetSupplierExtendedQuery) error {
				query.Result = supplier
				return nil
			})

			preferAccountID := supplier.ID
			resp, err := CreateLoginResponse(ctx, nil, token, userID, nil, preferAccountID, 0, false, 0)
			So(err, ShouldBeNil)

			loginResp.AvailableAccounts = accountsResp
			loginResp.Account = &pb.LoginAccount{
				ID:   supplier.ID,
				Name: "Alice",
				Type: pb.AccountType_supplier,
			}
			loginResp.Supplier = &pb.Supplier{
				ID:             supplier.ID,
				Name:           "Alice",
				Status:         status3.Status_P,
				ContactPersons: []*pb.ContactPerson{},
			}
			So(resp, ShouldDeepEqual, loginResp)
		})
	})
}
*/
