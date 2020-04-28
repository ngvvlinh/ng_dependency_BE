package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type Foo struct {
	ID        dot.ID
	AccountID dot.ID

	ABC       string `sq:"'abc_2'"`
	Def2      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

// +sqlgen
type Account struct {
	ID   dot.ID
	Name string
}

// +sqlgen:      Foo as foo
// +sqlgen:join: Account as a on foo.account_id = a.id
type FooWithAccount struct {
	Foo     *Foo
	Account *Account
}
