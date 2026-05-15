package user

import (
	"testing"

	"goweb/handlers/user"
)

func TestValidateCreds_Valid(t *testing.T) {
	creds := &user.Credentials{
		Username: "testuser",
		Password: "testpassword",
	}

	ok, errs := user.ValidateCreds(creds)

	if !ok {
		t.Errorf("expected valid credentials, got errors: %v", errs)
	}

	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d errors: %v", len(errs), errs)
	}
}

func TestValidateCreds_ShortUsername(t *testing.T) {
	creds := &user.Credentials{
		Username: "test",
		Password: "testpassword",
	}

	ok, errs := user.ValidateCreds(creds)

	if ok {
		t.Error("expected validation to fail for short username")
	}

	if errs["username"] != "username should contain at least 8 characters." {
		t.Errorf("expected username error message, got: %s", errs["username"])
	}
}

func TestValidateCreds_ShortPassword(t *testing.T) {
	creds := &user.Credentials{
		Username: "testuser",
		Password: "test",
	}

	ok, errs := user.ValidateCreds(creds)

	if ok {
		t.Error("expected validation to fail for short password")
	}

	if errs["password"] != "password should contain at least 9 characters." {
		t.Errorf("expected password error message, got: %s", errs["password"])
	}
}

func TestValidateCreds_BothInvalid(t *testing.T) {
	creds := &user.Credentials{
		Username: "test",
		Password: "test",
	}

	ok, errs := user.ValidateCreds(creds)

	if ok {
		t.Error("expected validation to fail for both invalid fields")
	}

	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}

	if _, ok := errs["username"]; !ok {
		t.Error("expected username error")
	}

	if _, ok := errs["password"]; !ok {
		t.Error("expected password error")
	}
}

func TestValidateCreds_ExactMinUsername(t *testing.T) {
	creds := &user.Credentials{
		Username: "12345678",
		Password: "testpassword",
	}

	ok, _ := user.ValidateCreds(creds)

	if !ok {
		t.Error("expected username with exactly 8 characters to be valid")
	}
}

func TestValidateCreds_ExactMinPassword(t *testing.T) {
	creds := &user.Credentials{
		Username: "testuser",
		Password: "123456789",
	}

	ok, _ := user.ValidateCreds(creds)

	if !ok {
		t.Error("expected password with exactly 9 characters to be valid")
	}
}
