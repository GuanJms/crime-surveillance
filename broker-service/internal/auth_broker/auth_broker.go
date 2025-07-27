package authbroker

var ab *AuthBroker

type AuthBroker struct{}

func NewAuthBroker() *AuthBroker {
	if ab == nil {
		ab = &AuthBroker{}
	}
	return ab
}
