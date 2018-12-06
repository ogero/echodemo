package types

const (
	Role_Administrator Role = iota - 1
	Role_User
	role_OOB
	Permission_Settings_List   = "settings.list"
	Permission_Settings_Update = "settings.update"
	Permission_JobRunner       = "settings.jobrunner"
	Permission_Users_List      = "users.list"
	Permission_Users_Create    = "users.create"
	Permission_Users_Update    = "users.update"
	Permission_Users_Delete    = "users.delete"
)

type Role int

func (o Role) FriendlyValue() string {
	if o < -1 || o >= role_OOB {
		return "?"
	}
	return [...]string{"Administrator", "User"}[o+1]
}

func (o Role) RBACValue() string {
	if o < -1 || o >= role_OOB {
		return "?"
	}
	return [...]string{"administrator", "user"}[o+1]
}

func (o Role) List() []Role {
	return []Role{
		Role_Administrator, Role_User,
	}
}
