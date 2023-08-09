package NCMB

import (
	"encoding/json"
)

type User struct {
	ncmb *NCMB
	Item
}

func (user *User) Save() (bool, error) {
	if user.ObjectId != "" && user.fields["password"] != nil {
		delete(user.fields, "password")
	}
	return user.Item.Save()
}

func (user *User) Login() (*User, error) {
	user.Item.ncmb = user.ncmb
	user.Item.ClassName = "users"
	path := "login"
	fields := map[string]interface{}{"userName": user.fields["userName"].(string), "password": user.fields["password"].(string)}
	params := ExecOptions{}
	params.Fields = &fields
	params.Path = &path
	return user.loginAndSignUp(params)
}

func (user *User) SignUpByAccount() (*User, error) {
	user.Item.ncmb = user.ncmb
	user.Item.ClassName = "users"
	fields := map[string]interface{}{"userName": user.fields["userName"].(string), "password": user.fields["password"].(string)}
	params := ExecOptions{}
	params.Fields = &fields
	return user.loginAndSignUp(params)
}

func (user *User) LoginWithMailAddress() (*User, error) {
	user.Item.ncmb = user.ncmb
	user.Item.ClassName = "users"
	fields := map[string]interface{}{"mailAddress": user.fields["mailAddress"].(string), "password": user.fields["password"].(string)}
	params := ExecOptions{}
	path := "login"
	params.Path = &path
	params.Fields = &fields
	return user.loginAndSignUp(params)
}

func (user *User) loginAndSignUp(params ExecOptions) (*User, error) {
	request := Request{ncmb: user.ncmb}
	params.ClassName = "users"
	data, err := request.Post(params)
	if err != nil {
		return nil, err
	}
	var hash map[string]interface{}
	err = json.Unmarshal(data, &hash)
	if err != nil {
		return nil, err
	}
	user.ncmb.SessionToken = hash["sessionToken"].(string)
	delete(hash, "sessionToken")
	user.Sets(hash)
	return user, nil
}

func (user *User) RequestSignUpEmail() (bool, error) {
	params := ExecOptions{}
	path := "requestMailAddressUserEntry"
	params.Path = &path
	return user.requestOrPasswordReset(params)
}

func (user *User) RequestPasswordReset() (bool, error) {
	params := ExecOptions{}
	path := "requestPasswordReset"
	params.Path = &path
	return user.requestOrPasswordReset(params)
}

func (user *User) requestOrPasswordReset(params ExecOptions) (bool, error) {
	params.ClassName = "users"
	fields := map[string]interface{}{"mailAddress": user.fields["mailAddress"].(string)}
	params.Fields = &fields
	user.ncmb.SessionToken = ""
	request := Request{ncmb: user.ncmb}
	data, err := request.Post(params)
	if err != nil {
		return false, err
	}
	var hash map[string]interface{}
	err = json.Unmarshal(data, &hash)
	if err != nil {
		return false, err
	}
	if _, ok := hash["createDate"]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (user *User) Logout() (bool, error) {
	params := ExecOptions{}
	path := "logout"
	params.Path = &path
	params.ClassName = "users"
	request := Request{ncmb: user.ncmb}
	_, err := request.Get(params)
	if err != nil {
		return false, err
	}
	user.ncmb.SessionToken = ""
	return true, nil
}
