package config

type Config struct {
	Email    string
	Password string
	Address  string
}

var (
	Cfg = Config{
		Email:    "r_mukhamedzhanov@bk.ru",
		Password: "01102003rus",
		Address:  "smtp.mail.ru:587",
	}
)
