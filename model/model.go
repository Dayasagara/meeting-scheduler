package model

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//DBConfig has information required to connect to DB
type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}
