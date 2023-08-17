package NCMB

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func SetUpAcl() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func TearDownAcl(ncmb NCMB) {
	fmt.Printf("TearDown %s\n", ncmb)
}

func TestItemSaveWithAcl(t *testing.T) {
	ncmb := SetUpAcl()
	userName, password := "testGoAcl", "testGoAcl"
	user, err := ncmb.SignUpByAccount(userName, password)
	if err != nil {
		t.Errorf("ncmb.SignUpByAccount() = %T, %s", user, err)
	}
	if user.ObjectId == "" {
		t.Errorf("user.ObjectId = %s, want not empty", user.ObjectId)
	}
	item := ncmb.Item("AclData")
	item.Set("message", "Hello, NCMB!")
	acl := ncmb.Acl()
	acl.SetPublicReadAccess(false)
	acl.SetUserReadAccess(user, true)
	acl.SetUserWriteAccess(user, true)
	item.SetAcl(acl)
	bol, err := item.Save()
	if err != nil {
		t.Errorf("item.Save() = %T, %s", item, err)
	}
	if item.ObjectId == "" {
		t.Errorf("item.ObjectId = %s, want not empty", item.ObjectId)
	}
	sessionToken := ncmb.SessionToken
	ncmb.SessionToken = ""
	bol, err = item.Fetch()
	if err == nil {
		t.Errorf("item.Fetch() = %t, %s", bol, err)
	}
	ncmb.SessionToken = sessionToken
	item.Fetch()
	acl, err = item.GetAcl()
	if err != nil {
		t.Errorf("item.GetAcl() = %T, %s", acl, err)
	}
	if acl.GetUserReadAccess(user) != true {
		t.Errorf("acl.GetUserReadAccess = %t, want true", acl.GetUserReadAccess(user))
	}
	if acl.GetUserWriteAccess(user) != true {
		t.Errorf("acl.GetUserWriteAccess = %t, want true", acl.GetUserWriteAccess(user))
	}
	item.Delete()
	bol, err = user.Delete()
	if err != nil {
		t.Errorf("user.Delete() = %T, %s", bol, err)
	}
}
