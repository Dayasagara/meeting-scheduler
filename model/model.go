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

type User struct {
	UserID   int    `json:"userID" gorm:"column:userID;primary_key;AUTO_INCREMENT"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
}

//CreateResponse defines response code and message with token
type CreateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Availability struct {
	UserID    int    `json:"userID" gorm:"column:userID"`
	Date      string `json:"date" gorm:"column:date"`
	StartSlot string `json:"startSlot" gorm:"column:fromTime"`
	EndSlot   string `json:"endSlot" gorm:"column:fromTime"`
	Duration  int    `json:"durationTime"`
}
