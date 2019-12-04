package model

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFoo(&Foo{})

type Foo struct {
	ID        dot.ID
	AccountID dot.ID

	ABC       string `sq:"'abc_2'"`
	Def2      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenAccount(&Account{})

type Account struct {
	ID   dot.ID
	Name string
}

var _ = sqlgenFooWithAccount(
	&FooWithAccount{}, &Foo{}, "foo",
	sq.JOIN, &Account{}, "a", "foo.account_id = a.id",
)

type FooWithAccount struct {
	Foo     *Foo
	Account *Account
}
