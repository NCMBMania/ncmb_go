package NCMB

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func SetUpUser() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func TearDownUser(ncmb NCMB) {
	fmt.Printf("TearDown %s\n", ncmb)
}

func TestLogin(t *testing.T) {
	ncmb := SetUpUser()
	userName, password := os.Getenv("USER_NAME"), os.Getenv("PASSWORD")
	user, err := ncmb.Login(userName, password)
	if err != nil {
		t.Errorf("ncmb.Login() = %T, %s", user, err)
	}
	if user.ObjectId == "" {
		t.Errorf("user.ObjectId = %s, want not empty", user.ObjectId)
	}
	TearDownUser(ncmb)
}

func TestRegister(t *testing.T) {
	ncmb := SetUpUser()
	userName, password := "testGo", "testGo"
	user, err := ncmb.SignUpByAccount(userName, password)
	if err != nil {
		t.Errorf("ncmb.SignUpByAccount() = %T, %s", user, err)
	}
	if user.ObjectId == "" {
		t.Errorf("user.ObjectId = %s, want not empty", user.ObjectId)
	}
	bol, err := user.Delete()
	if err != nil {
		t.Errorf("user.Delete() = %T, %s", bol, err)
	}
	TearDownUser(ncmb)
}

func TestLoginWithMailAddress(t *testing.T) {
	ncmb := SetUpUser()
	mailAddress, password := os.Getenv("EMAIL"), os.Getenv("PASSWORD")
	user, err := ncmb.LoginWithMailAddress(mailAddress, password)
	if err != nil {
		t.Errorf("ncmb.LoginWithMailAddress() = %T, %s", user, err)
	}
	if user.ObjectId == "" {
		t.Errorf("user.ObjectId = %s, want not empty", user.ObjectId)
	}
	TearDownUser(ncmb)
}

func TestRequestSignUpEmail(t *testing.T) {
	ncmb := SetUpUser()
	mailAddress := os.Getenv("EMAIL")
	bol, err := ncmb.RequestSignUpEmail(mailAddress)
	if err != nil {
		t.Errorf("ncmb.RequestSignUpEmail() = %T, %s", bol, err)
	}
	if bol == false {
		t.Errorf("ncmb.RequestSignUpEmail() = %T", bol)
	}
	TearDownUser(ncmb)
}

func TestLogout(t *testing.T) {
	ncmb := SetUpUser()
	mailAddress, password := os.Getenv("EMAIL"), os.Getenv("PASSWORD")
	user, err := ncmb.LoginWithMailAddress(mailAddress, password)
	if err != nil {
		t.Errorf("ncmb.LoginWithMailAddress() = %T, %s", user, err)
	}
	if user.ObjectId == "" {
		t.Errorf("user.ObjectId = %s, want not empty", user.ObjectId)
	}
	if ncmb.SessionToken == "" {
		t.Errorf("ncmb.SessionToken = %s, want not empty", ncmb.SessionToken)
	}
	bol, err := ncmb.Logout()
	if err != nil {
		t.Errorf("ncmb.Logout() = %T, %s", bol, err)
	}
	if bol == false {
		t.Errorf("ncmb.Logout() = %T", bol)
	}
	if ncmb.SessionToken != "" {
		t.Errorf("ncmb.SessionToken = %s, want empty", ncmb.SessionToken)
	}
	TearDownUser(ncmb)
}

func TestUpdate(t *testing.T) {
	ncmb := SetUpUser()
	userName := "testGo2"
	user, err := ncmb.SignUpByAccount(userName, "testGo")
	if err != nil {
		t.Errorf("ncmb.SignUpByAccount() = %T, %s", user, err)
	}
	if user.ObjectId == "" {
		t.Errorf("user.ObjectId = %s, want not empty", user.ObjectId)
	}
	val, err := user.GetString("userName")
	if err != nil {
		t.Errorf("user.GetString(userName) = %T, %s", val, err)
	}
	if val != userName {
		t.Errorf("user.GetString(userName) = %s, want %s", val, userName)
	}
	displayName := "testGo3"
	user.Set("displayName", displayName)
	bol, err := user.Save()
	if err != nil {
		t.Errorf("user.Save() = %T, %s", bol, err)
	}
	if bol == false {
		t.Errorf("user.Save() = %T", bol)
	}
	bol, err = user.Fetch()
	if err != nil {
		t.Errorf("user.Fetch() = %T, %s", bol, err)
	}
	if bol == false {
		t.Errorf("user.Fetch() = %T", bol)
	}
	val, err = user.GetString("displayName")
	if err != nil {
		t.Errorf("user.GetString(displayName) = %T, %s", val, err)
	}
	if val != displayName {
		t.Errorf("user.GetString(displayName) = %s, want %s", val, displayName)
	}
	bol, err = user.Delete()
	if err != nil {
		t.Errorf("user.Delete() = %T, %s", bol, err)
	}
	TearDownUser(ncmb)
}
