package interfaces

import "github.com/Dayasagara/meeting-scheduler/model"

//DBEngine is used to call DB Methods from handler
var DBEngine DBInterface

//DBInterface contains all the DB methods
type DBInterface interface {
	DBConnect(model.DBConfig) error
	CheckUser(string) error
	CreateUser(model.User) error
}
