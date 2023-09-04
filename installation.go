package NCMB

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Installation struct {
	ncmb *NCMB
	Item
}

func (installation *Installation) Save() (bool, error) {
	if installation.Item.ObjectId == "" {
		return installation.Create()
	} else {
		return installation.Update()
	}
}

func (installation *Installation) Create() (bool, error) {
	if valid, err := installation.valid(); !valid {
		return false, err
	}
	return installation.Item.Create()
}

func (installation *Installation) Update() (bool, error) {
	if valid, err := installation.valid(); !valid {
		return false, err
	}
	return installation.Item.Update()
}

func (installation *Installation) valid() (bool, error) {
	deviceToken, err := installation.GetString("deviceToken")
	if err != nil {
		return false, err
	}
	if deviceToken == "" {
		return false, fmt.Errorf("deviceToken is required")
	}
	deviceType, err := installation.GetString("deviceType")
	if err != nil {
		return false, err
	}
	if slices.Index([]string{"ios", "android"}, deviceType) == -1 {
		return false, fmt.Errorf("deviceType accepts ios or android, not %s", deviceType)
	}
	return true, nil
}
