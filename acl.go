package NCMB

type Acl struct {
	permissions map[string]map[string]bool
}

func (acl *Acl) Sets(valueMap map[string]map[string]bool) {
	for k, v := range valueMap {
		for k2, v2 := range v {
			acl.setAccess(k, k2, v2)
		}
	}
}

func (acl *Acl) SetPublicReadAccess(b bool) *Acl {
	return acl.setAccess("*", "read", b)
}

func (acl *Acl) GetPublicReadAccess() bool {
	return acl.getAccess("*", "read")
}

func (acl *Acl) GetPublicWriteAccess() bool {
	return acl.getAccess("*", "write")
}

func (acl *Acl) GetUserReadAccess(user *User) bool {
	return acl.getAccess(user.ObjectId, "read")
}

func (acl *Acl) GetUserWriteAccess(user *User) bool {
	return acl.getAccess(user.ObjectId, "write")
}

func (acl *Acl) GetRoleReadAccess(role string) bool {
	return acl.getAccess(role, "read")
}

func (acl *Acl) GetRoleWriteAccess(role string) bool {
	return acl.getAccess(role, "write")
}

func (acl *Acl) getAccess(key string, access string) bool {
	if acl.permissions == nil {
		return false
	}
	if acl.permissions[key] == nil {
		return false
	}
	return acl.permissions[key][access]
}

func (acl *Acl) SetPublicWriteAccess(b bool) *Acl {
	return acl.setAccess("*", "write", b)
}

func (acl *Acl) SetUserReadAccess(user *User, b bool) *Acl {
	return acl.setAccess(user.ObjectId, "read", b)
}

func (acl *Acl) SetUserWriteAccess(user *User, b bool) *Acl {
	return acl.setAccess(user.ObjectId, "write", b)
}

func (acl *Acl) SetRoleReadAccess(role string, b bool) *Acl {
	return acl.setAccess(role, "read", b)
}

func (acl *Acl) SetRoleWriteAccess(role string, b bool) *Acl {
	return acl.setAccess(role, "write", b)
}

func (acl *Acl) setAccess(key string, access string, b bool) *Acl {
	if acl.permissions == nil {
		acl.permissions = make(map[string]map[string]bool)
	}
	if acl.permissions[key] == nil {
		acl.permissions[key] = make(map[string]bool)
	}
	if b == true {
		acl.permissions[key][access] = b
	} else {
		delete(acl.permissions[key], access)
	}
	return acl
}

func (acl Acl) ToJSON() (map[string]map[string]bool, error) {
	permissions := make(map[string]map[string]bool)
	if acl.permissions == nil {
		acl.permissions = make(map[string]map[string]bool)
	}
	for k, v := range acl.permissions {
		for k2, v2 := range v {
			if permissions[k] == nil {
				permissions[k] = make(map[string]bool)
			}
			if v2 {
				permissions[k][k2] = v2
			}
		}
	}
	return permissions, nil
}
