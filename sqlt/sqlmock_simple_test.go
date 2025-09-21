package sqlt

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestAnyTime tests the AnyTime matcher
func TestAnyTime(t *testing.T) {
	anyTime := AnyTime{}

	// Test with time.Time
	now := time.Now()
	if !anyTime.Match(now) {
		t.Error("AnyTime should match time.Time")
	}

	// Test with non-time value
	if anyTime.Match("not a time") {
		t.Error("AnyTime should not match non-time values")
	}

	if anyTime.Match(123) {
		t.Error("AnyTime should not match non-time values")
	}
}

// TestNewSqlmock tests the NewSqlmock function
func TestNewSqlmock(t *testing.T) {
	mockDB := NewSqlmock()

	if mockDB == nil {
		t.Error("NewSqlmock should return a non-nil MockDB")
	}

	if mockDB.DB == nil {
		t.Error("NewSqlmock should return a non-nil database")
	}

	if mockDB.Sqlmock == nil {
		t.Error("NewSqlmock should return a non-nil mock")
	}
}

// TestGorm tests the Gorm function
func TestGorm(t *testing.T) {
	mockDB := NewSqlmock()

	gormMock := mockDB.Gorm()
	if gormMock == nil {
		t.Error("Gorm should return a non-nil MockGorm")
	}

	if gormMock.DB == nil {
		t.Error("Gorm should return a non-nil GORM database")
	}
}

// TestMockDBOperations tests basic mock operations
func TestMockDBOperations(t *testing.T) {
	mockDB := NewSqlmock()

	// Test basic query expectation
	mockDB.Sqlmock.ExpectQuery("SELECT \\* FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John"))

	rows, err := mockDB.DB.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Query should not return error, got: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Expected at least one row")
	}

	if err := mockDB.Sqlmock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}

// TestMockDBExec tests exec operations
func TestMockDBExec(t *testing.T) {
	mockDB := NewSqlmock()

	// Test exec expectation
	mockDB.Sqlmock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := mockDB.DB.Exec("INSERT INTO users (name) VALUES (?)", "John")
	if err != nil {
		t.Errorf("Exec should not return error, got: %v", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		t.Errorf("LastInsertId should not return error, got: %v", err)
	}

	if lastInsertID != 1 {
		t.Errorf("Expected last insert ID 1, got %d", lastInsertID)
	}

	if err := mockDB.Sqlmock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}
