package config

type AuthTimeout struct {
	LoginTimeout           int // seconds
	RegisterTimeout        int // seconds
	RecoverPasswordTimeout int // seconds
	ResetPasswordTimeout   int // seconds
}

type EmailTimeout struct {
	EmailSendTimeout   int // seconds
	VerifyEmailTimeout int // seconds
}

type UserTimeout struct {
	EmailAlreadyExistsTimeout int // seconds
	UsernameExistsTimeout     int // seconds
	CreateUserTimeout         int // seconds
	GetUserTimeout            int // seconds
	TokenVersionTimeout       int // seconds
}
