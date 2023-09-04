package NCMB

import (
	"encoding/json"
	"fmt"
)

type Role struct {
	ncmb        *NCMB
	Roles       []Role
	addRoles    []Role
	removeRoles []Role
	Users       []User
	addUsers    []User
	removeUsers []User
	Item
}

func (role *Role) AddRole(r Role) *Role {
	role.addRoles = append(role.Roles, r)
	return role
}

func (role *Role) AddUser(u User) *Role {
	role.addUsers = append(role.Users, u)
	return role
}

func (role *Role) Save() (bool, error) {
	if role.ObjectId == "" {
		return role.Create()
	} else {
		return role.Update()
	}
}

func (role *Role) Create() (bool, error) {
	params := ExecOptions{}
	fields, err := role.Fields()
	if err != nil {
		return false, err
	}
	params.Fields = &fields
	request := Request{ncmb: role.ncmb}
	params.ClassName = role.ClassName
	if data, err := request.Post(params); err != nil {
		return false, err
	} else {
		var hash map[string]interface{}
		if err = json.Unmarshal(data, &hash); err != nil {
			return false, err
		}
		role.Item.Sets(hash)
	}
	return true, nil
}

func (role *Role) Update() (bool, error) {
	params := ExecOptions{}
	if fields, err := role.Fields(); err != nil {
		return false, err
	} else {
		params.Fields = &fields
	}
	params.ObjectId = &role.ObjectId
	request := Request{ncmb: role.ncmb}
	params.ClassName = role.ClassName
	if data, err := request.Put(params); err != nil {
		return false, err
	} else {
		var hash map[string]interface{}
		if err = json.Unmarshal(data, &hash); err != nil {
			return false, err
		}
		role.Item.Sets(hash)
	}
	return true, nil
}

func (role *Role) GetObjectId() string {
	return role.Item.ObjectId
}

func (role *Role) GetClassName() string {
	return "role"
}

func (role *Role) FetchRole() ([]Role, error) {
	query := role.ncmb.Query("roles")
	query.RelatedTo(role, "belongRole")
	items, err := query.FetchAll()
	if err != nil {
		return nil, err
	}
	roles := []Role{}
	for _, item := range items {
		role := Role{ncmb: role.ncmb}
		role.Item = item
		roles = append(roles, role)
	}
	return roles, nil
}

func (role *Role) FetchUser() ([]User, error) {
	query := role.ncmb.Query("users")
	query.RelatedTo(role, "belongUser")
	items, err := query.FetchAll()
	if err != nil {
		return nil, err
	}
	users := []User{}
	for _, item := range items {
		user := User{ncmb: role.ncmb}
		user.Item = item
		users = append(users, user)
	}
	return users, nil
}

func (role *Role) Fields() (map[string]interface{}, error) {
	fields := map[string]interface{}{}
	if roleName, err := role.Item.GetString("roleName"); err == nil {
		fields["roleName"] = roleName
	}
	if acl, err := role.Item.GetAcl(); err == nil {
		fields["acl"] = acl
	}
	if len(role.addRoles) > 0 && len(role.removeRoles) > 0 {
		return nil, fmt.Errorf("Can not add and remove role in same time")
	}
	if len(role.addUsers) > 0 && len(role.removeUsers) > 0 {
		return nil, fmt.Errorf("Can not add and remove user in same time")
	}
	if len(role.addRoles) > 0 {
		fields["belongRole"] = map[string]interface{}{
			"__op":    "AddRelation",
			"objects": ToRoleObjects(role.addRoles),
		}
	}
	if len(role.removeRoles) > 0 {
		fields["belongRole"] = map[string]interface{}{
			"__op":    "RemoveRelation",
			"objects": ToRoleObjects(role.removeRoles),
		}
	}
	if len(role.addUsers) > 0 {
		fields["belongUser"] = map[string]interface{}{
			"__op":    "AddRelation",
			"objects": ToUserbjects(role.addUsers),
		}
	}
	if len(role.removeUsers) > 0 {
		fields["belongUser"] = map[string]interface{}{
			"__op":    "RemoveRelation",
			"objects": ToUserbjects(role.removeUsers),
		}
	}
	return fields, nil
}

func ToRoleObjects(ary []Role) []map[string]string {
	objects := []map[string]string{}
	for _, r := range ary {
		objects = append(objects, r.ToPointer())
	}
	return objects
}

func ToUserbjects(ary []User) []map[string]string {
	objects := []map[string]string{}
	for _, u := range ary {
		objects = append(objects, u.ToPointer())
	}
	return objects
}

func (role *Role) ToPointer() map[string]string {
	return map[string]string{
		"__type":    "Pointer",
		"className": "role",
		"objectId":  role.Item.ObjectId,
	}
}
