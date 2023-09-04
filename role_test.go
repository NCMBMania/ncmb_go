package NCMB

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func SetUpRole() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func TearDownRole(ncmb NCMB) {
	fmt.Printf("TearDown %s\n", ncmb)
}

func TestSaveRole(t *testing.T) {
	ncmb := SetUpRole()
	role := ncmb.Role("TestRole")
	bol, err := role.Save()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if bol == false {
		t.Errorf("Error: %s", err)
	}
	if role.ObjectId == "" {
		t.Errorf("Error: %s", err)
	} else {
		fmt.Printf("objectId: %s\n", role.ObjectId)
	}
	if bol, err := role.Delete(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	TearDownRole(ncmb)
}

func TestAddRole(t *testing.T) {
	ncmb := SetUpRole()
	role1 := ncmb.Role("TestRole1")
	if _, err := role1.Save(); err != nil {
		t.Errorf("Error: %s", err)
	}
	role2 := ncmb.Role("TestRole2")
	if _, err := role2.Save(); err != nil {
		t.Errorf("Error: %s", err)
	}
	role3 := ncmb.Role("TestRole3")
	if _, err := role3.Save(); err != nil {
		t.Errorf("Error: %s", err)
	}
	role1.AddRole(role2)
	role1.AddRole(role3)
	if bol, err := role1.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
		role1.Delete()
		role2.Delete()
		role3.Delete()
	}
	ary, err := role1.FetchRole()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(ary) != 2 {
		t.Errorf("Expected: %d, Actual: %d", 2, len(ary))
	}
	if ary[1].ObjectId != role2.ObjectId {
		t.Errorf("Expected: %s, Actual: %s", role2.ObjectId, ary[0].ObjectId)
	}
	role1.Delete()
	role2.Delete()
	role3.Delete()
	TearDownRole(ncmb)
}

func TestAddUser(t *testing.T) {
	ncmb := SetUpRole()
	role1 := ncmb.Role("TestRole1")
	if _, err := role1.Save(); err != nil {
		t.Errorf("Error: %s", err)
	}
	userName, password := "roleGo", "roleGo"
	user, err := ncmb.SignUpByAccount(userName, password)
	role1.AddUser(*user)
	if bol, err := role1.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
		role1.Delete()
		user.Delete()
	}
	ary, err := role1.FetchUser()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(ary) != 1 {
		t.Errorf("Error: %s, %d", err, len(ary))
	}
	if ary[0].ObjectId != user.ObjectId {
		t.Errorf("Expected: %s, Actual: %s", user.ObjectId, ary[0].ObjectId)
	}
	role1.Delete()
	user.Delete()
	TearDownRole(ncmb)
}

func TestUpdateRole(t *testing.T) {
	ncmb := SetUpRole()
	role := ncmb.Role("TestRole1")
	role.Save()
	roleName := "TestRole2"
	role.Set("roleName", roleName)
	role.Save()
	if name, err := role.GetString("roleName"); err != nil || name != roleName {
		t.Errorf("Expected: %s, Actual: %s", roleName, name)
	}
	role.Delete()
}
