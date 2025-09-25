package sqlt

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// AnyTime is a struct that matches any time type, ignoring the value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type MockDB struct {
	*sql.DB
	sqlmock.Sqlmock
}

type MockGorm struct {
	*MockDB
	DB *gorm.DB
}

func NewSqlmock() (*MockDB, error) {
	db, mock, err := sqlmock.New() // mock db
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlmock: %v", err)
	}
	return &MockDB{DB: db, Sqlmock: mock}, nil
}

func (m *MockDB) Gorm() (*MockGorm, error) {
	// create gorm.DB
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      m.DB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %v", err)
	}
	return &MockGorm{MockDB: m, DB: db}, nil
}
