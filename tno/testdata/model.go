//go:generate tno -type=DealerBusinessmanHistory

package model

import (
	"database/sql"
	"time"

	"github.com/YunzhanghuOpen/null.v3"
	"github.com/jmoiron/sqlx"
)

type DbInterface interface {
	Rebind(query string) string
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

type tableI interface {
	schema() string
	fields() []string
	getByID(id int64) interface{}
	selectByIDs(ids []int64) []interface{}
}

var iDealerBusinessmanHistory tableI

type DealerBusinessmanHistory struct {
	ID                int64     `db:"I_ID"`
	DealerID          string    `db:"CH_DEALER_ID"`
	BusinessManID     int64     `db:"I_BUSINESS_MAN_ID"`
	Type              int64     `db:"I_TYPE"`
	StartServiceAt    null.Time `db:"D_START_SERVICE_AT"`
	EndServiceAt      null.Time `db:"D_END_SERVICE_AT"`
	CreatedAt         null.Time `db:"D_CREATED_AT"`
	CreateManagerName string    `db:"CH_CREATE_MANAGE_NAME"`
	UpdatedAt         null.Time `db:"D_UPDATED_AT"`
	ModifyManagerName string    `db:"CH_MODIFY_MANAGE_NAME"`
	IsDeleted         bool      `db:"B_IS_DELETED"`
}

func (DealerBusinessmanHistory) schema() string {
	return `dealer_businessman_history`
}

func (DealerBusinessManRelation) fields() []string {
	return iDealerBusinessmanHistory.fields()
}

type DealerBusinessManRelation struct {
	ID               int64     `db:"I_ID"`
	DealerID         string    `db:"CH_DEALER_ID"`
	BusinessManID    int64     `db:"I_BUSINESS_MAN_ID"`
	Type             int32     `db:"I_TYPE"`
	Status           int32     `db:"I_STATUS"`
	CreatedAt        null.Time `db:"D_CREATED_AT"`
	UpdatedAt        null.Time `db:"D_UPDATED_AT"`
	StartServiceAt   time.Time `db:"D_START_SERVICE_AT"`
	CreateManageName string    `db:"CH_CREATE_MANAGE_NAME"`
	ModifyManageName string    `db:"CH_MODIFY_MANAGE_NAME"`
}
