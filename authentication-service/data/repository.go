package data

type Repository interface {
	CreateUser(user *User) error
	AuthenticateUser(credentials *UserLoginCredentials) (bool, error)
	GetUserInfo(username string) (*User, error)
	UpdateLoginTime(user *User) error
	ChangeUserRoleTo(id string, role string) error
	GetAllUsers() ([]*User, error)
}
