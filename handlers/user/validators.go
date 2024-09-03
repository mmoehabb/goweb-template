package user

func ValidateCreds(creds *Credentials) (bool, map[string]string) {
  ok := true
  errs := make(map[string]string)
  if len(creds.Username) < 8 {
    errs["username"] = "username should contain at least 8 characters."
    ok = false
  }
  if len(creds.Password) < 9 {
    errs["password"] = "password should contain at least 9 characters."
    ok = false
  }
  return ok, errs
}

