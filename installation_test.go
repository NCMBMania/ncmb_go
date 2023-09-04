package NCMB

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func SetUpInstallation() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func TearDownInstallation(ncmb NCMB) {
	fmt.Printf("TearDown %s\n", ncmb)
}

func TestSaveInstallation(t *testing.T) {
	ncmb := SetUpInstallation()
	installation := ncmb.Installation()
	installation.Set("deviceToken", "testDeviceToken")
	installation.Set("deviceType", "ios")
	if bol, err := installation.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	if bol, err := installation.Delete(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
}

func TestUpdateInstallation(t *testing.T) {
	ncmb := SetUpInstallation()
	installation := ncmb.Installation()
	installation.Set("deviceToken", "testDeviceToken")
	installation.Set("deviceType", "ios")
	if bol, err := installation.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	newToken := "testDeviceToken2"
	installation.Set("deviceToken", newToken)
	if bol, err := installation.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	if token, err := installation.GetString("deviceToken"); err != nil || token != newToken {
		t.Errorf("Expected: %s, Actual: %s", newToken, token)
	}
	if bol, err := installation.Delete(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
}

func TestFetchInstallation(t *testing.T) {
	ncmb := SetUpInstallation()
	installation := ncmb.Installation()
	token := "testDeviceToken"
	installation.Set("deviceToken", token)
	installation.Set("deviceType", "ios")
	if bol, err := installation.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	installation2 := ncmb.Installation()
	installation2.Set("objectId", installation.ObjectId)
	bol, err := installation2.Fetch()
	if err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	if t2, err := installation2.GetString("deviceToken"); err != nil || t2 != token {
		t.Errorf("Expected: %s, Actual: %s", token, t2)
	}
	if bol, err := installation.Delete(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
}

func TestFetchAllInstallation(t *testing.T) {
	ncmb := SetUpInstallation()
	installation := ncmb.Installation()
	token := "testDeviceToken"
	installation.Set("deviceToken", token)
	installation.Set("deviceType", "ios")
	if bol, err := installation.Save(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
	query := ncmb.Query("installations")
	query.EqualTo("deviceType", "ios")
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(items) != 1 {
		t.Errorf("Expected: %d, Actual: %d", 1, len(items))
	}
	if bol, err := installation.Delete(); err != nil {
		t.Errorf("Error: %s, %T", err, bol)
	}
}
