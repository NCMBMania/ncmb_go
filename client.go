package NCMB

// import "fmt"

type NCMB struct {
	ApplicationKey string
	ClientKey      string
	SessionToken   string
}

func Initialize(applicationKey string, clientKey string) NCMB {
	ncmb := NCMB{ApplicationKey: applicationKey, ClientKey: clientKey}
	return ncmb
}

func (ncmb *NCMB) Item(className string) Item {
	return Item{ncmb: ncmb, ClassName: className}
}

func (ncmb *NCMB) Query(className string) Query {
	return Query{ncmb: ncmb, className: className}
}

func (ncmb *NCMB) File(fileName string) File {
	file := File{ncmb: ncmb}
	file.Set("fileName", fileName)
	file.Item.ClassName = "file"
	file.Item.ncmb = ncmb
	return file
}

func (ncmb *NCMB) Push() Push {
	push := Push{ncmb: ncmb}
	push.Item.ClassName = "push"
	push.Item.ncmb = ncmb
	return push
}

func (ncmb *NCMB) GeoPoint(latitude float64, longitude float64) GeoPoint {
	return GeoPoint{Type: "GeoPoint", Latitude: latitude, Longitude: longitude}
}

func (ncmb *NCMB) Login(userName string, password string) (*User, error) {
	user := User{ncmb: ncmb}
	user.Set("userName", userName)
	user.Set("password", password)
	return user.Login()
}

func (ncmb *NCMB) SignUpByAccount(userName string, password string) (*User, error) {
	user := User{ncmb: ncmb}
	user.Set("userName", userName)
	user.Set("password", password)
	return user.SignUpByAccount()
}

func (ncmb *NCMB) LoginWithMailAddress(mailAddress string, password string) (*User, error) {
	user := User{ncmb: ncmb}
	user.Set("mailAddress", mailAddress)
	user.Set("password", password)
	return user.LoginWithMailAddress()
}

func (ncmb *NCMB) RequestSignUpEmail(mailAddress string) (bool, error) {
	user := User{ncmb: ncmb}
	user.Set("mailAddress", mailAddress)
	return user.RequestSignUpEmail()
}

func (ncmb *NCMB) Logout() (bool, error) {
	user := User{ncmb: ncmb}
	return user.Logout()
}

func (ncmb *NCMB) RequestPasswordReset(mailAddress string) (bool, error) {
	user := User{ncmb: ncmb}
	user.Set("mailAddress", mailAddress)
	return user.RequestPasswordReset()
}

func (ncmb *NCMB) Acl() Acl {
	return Acl{}
}

func (ncmb *NCMB) Role(name string) Role {
	role := Role{ncmb: ncmb}
	role.Set("roleName", name)
	role.Item.ClassName = "roles"
	role.Item.ncmb = ncmb
	role.ClassName = "roles"
	return role
}

func (ncmb *NCMB) Installation() Installation {
	installation := Installation{ncmb: ncmb}
	installation.Item.ClassName = "installations"
	installation.Item.ncmb = ncmb
	installation.ClassName = "installations"
	return installation
}
