package model

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFoo(&Foo{})

type Foo struct {
	ID        int64
	AccountID int64

	ABC       string `sq:"'abc_2'"`
	Def2      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenAccount(&Account{})

type Account struct {
	ID   int64
	Name string
}

var _ = sqlgenFooWithAccount(
	&FooWithAccount{}, &Foo{}, sq.AS("foo"),
	sq.JOIN, &Account{}, sq.AS("a"), "foo.account_id = a.id",
)

type FooWithAccount struct {
	Foo     *Foo
	Account *Account
}
