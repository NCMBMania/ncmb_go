package NCMB

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func SetUpPush() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func TearDownPush(ncmb NCMB) {
	fmt.Printf("TearDown %s\n", ncmb)
}

func TestPushSave(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	push.Set("message", "Hello, NCMB!")
	push.Set("target", []string{"ios", "android"})
	push.Set("immediateDeliveryFlag", true)
	bol, err := push.Save()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Save Error")
	}
	TearDownFile(ncmb)
}

func TestPushSaveWithDeliveryTime(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	push.Set("message", "Hello, NCMB!")
	push.Set("target", []string{"ios", "android"})
	t1 := time.Now()
	push.Set("deliveryTime", t1.Add(time.Hour*10))
	bol, err := push.Save()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Save Error")
	}
	TearDownFile(ncmb)
}

func TestPushSaveFailed(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	push.Set("message", "Hello, NCMB!")
	push.Set("immediateDeliveryFlag", true)
	if _, err := push.Save(); err == nil {
		t.Error("Save Error")
	} else if err.Error() != "target is required. %!s(<nil>)" {
		fmt.Println(err.Error())
	}
	TearDownFile(ncmb)
}

func TestPushSaveFailed2(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	push.Set("message", "Hello, NCMB!")
	// push.Set("immediateDeliveryFlag", true)
	push.Set("target", []string{"ios", "android"})
	if _, err := push.Save(); err == nil {
		t.Error("Save Error")
	} else if err.Error() != "immediateDeliveryFlag or deliveryTime is required" {
		t.Error(err.Error())
	}
	TearDownFile(ncmb)
}

func TestPushFetch(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	msg := "Hello, NCMB!"
	push.Set("message", msg)
	push.Set("target", []string{"ios", "android"})
	push.Set("immediateDeliveryFlag", true)
	bol, err := push.Save()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Save Error")
	}
	push2 := ncmb.Push()
	push2.ObjectId = push.ObjectId
	bol, err = push2.Fetch()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Fetch Error")
	}
	message, err := push2.GetString("message")
	if err != nil {
		t.Error(err)
	}
	if message != msg {
		t.Error("Fetch Error")
	}
	TearDownFile(ncmb)
}

func TestPushUpdate(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	push.Set("message", "Hello, NCMB!")
	push.Set("target", []string{"ios", "android"})
	push.Set("immediateDeliveryFlag", true)
	bol, err := push.Save()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Save Error")
	}
	push.Set("message", "Hello, NCMB! Update")
	bol, err = push.Save()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Save Error")
	}
	TearDownFile(ncmb)
}

func TestPushDelete(t *testing.T) {
	ncmb := SetUpFile()
	push := ncmb.Push()
	push.Set("message", "Hello, NCMB!")
	push.Set("target", []string{"ios", "android"})
	push.Set("immediateDeliveryFlag", true)
	bol, err := push.Save()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Save Error")
	}
	bol, err = push.Delete()
	if err != nil {
		t.Error(err)
	}
	if bol != true {
		t.Error("Delete Error")
	}
	TearDownFile(ncmb)
}

func TestPushFetchAll(t *testing.T) {
	ncmb := SetUpFile()
	query := ncmb.Query("push")
	query.Limit(5)
	query.EqualTo("message", "Hello, NCMB!")
	ary, err := query.FetchAll()
	if err != nil {
		t.Error(err)
	}
	if len(ary) != 5 {
		t.Errorf("FetchAll Error. Expected 5, but %d", len(ary))
	}
	TearDownFile(ncmb)
}
