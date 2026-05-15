package users

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"goweb/db"
	"goweb/db/users"
)

func setupTestDB(t *testing.T) *gorm.DB {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := gorm.Open(sqlite.Open("file:"+dbPath+"?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&users.DataModel{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestAdd_CreatesUser(t *testing.T) {
	gormDB := setupTestDB(t)

	db.SetConnectionForTest(gormDB)
	defer db.ResetConnectionForTest()

	err := users.Add("testuser", "testpassword")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	var count int64
	gormDB.Model(&users.DataModel{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}
}

func TestAdd_DuplicateUsername(t *testing.T) {
	gormDB := setupTestDB(t)

	db.SetConnectionForTest(gormDB)
	defer db.ResetConnectionForTest()

	err := users.Add("testuser", "testpassword")
	if err != nil {
		t.Fatalf("first Add failed: %v", err)
	}

	err = users.Add("testuser", "anotherpassword")
	if err == nil {
		t.Error("expected error for duplicate username, got nil")
	}
}

func TestGet_ExistingUser(t *testing.T) {
	gormDB := setupTestDB(t)

	db.SetConnectionForTest(gormDB)
	defer db.ResetConnectionForTest()

	err := users.Add("testuser", "testpassword")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	user, err := users.Get("testuser")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got: %s", user.Username)
	}

	if user.Password != "testpassword" {
		t.Errorf("expected password 'testpassword', got: %s", user.Password)
	}
}

func TestGet_NonExistingUser(t *testing.T) {
	gormDB := setupTestDB(t)

	db.SetConnectionForTest(gormDB)
	defer db.ResetConnectionForTest()

	_, err := users.Get("nonexistent")
	if err == nil {
		t.Error("expected error for non-existing user, got nil")
	}
}
