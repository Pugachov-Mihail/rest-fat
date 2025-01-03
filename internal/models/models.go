package models

type Permission int

const (
	PermissionUser    = "PermissionUser"
	PermissionAdmin   = "PermissionAdmin"
	PermissionTrainer = "PermissionTrainer"
)

var Permissions = map[Permission]string{
	0: PermissionUser,
	2: PermissionAdmin,
	4: PermissionTrainer,
}

type User struct {
	Id        int64
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      Role
}

type Role struct {
	Id         int64
	Name       string
	Value      int
	Permission Permission
}
