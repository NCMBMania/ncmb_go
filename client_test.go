package NCMB

import (
	// "fmt"
	"testing"
)

func TestInitialize(t *testing.T) {
	applicationKey, clientKey := "aaa", "bbb"
	ncmb := Initialize(applicationKey, clientKey)
	if ncmb.ApplicationKey != applicationKey {
		t.Errorf("ncmb.ApplicationKey = %s, want %s", ncmb.ApplicationKey, applicationKey)
	}
	if ncmb.ClientKey != clientKey {
		t.Errorf("ncmb.ClientKey = %s, want %s", ncmb.ClientKey, clientKey)
	}
}