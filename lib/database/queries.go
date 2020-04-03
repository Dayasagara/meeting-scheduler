package database

import (
	"errors"

	"github.com/Dayasagara/meeting-scheduler/model"
)

//Check if the email already exists
func (dc *DBRepo) CheckUser(email string) error {
	rows := dc.GormDB.Debug().Exec(`SELECT * FROM "users" where "email"=?`, email).RowsAffected
	if rows == 0 {
		return nil
	}
	return errors.New("Email already exists")
}

//Create user
func (dc *DBRepo) CreateUser(user model.User) error {
	var err error
	rows := dc.GormDB.Debug().Create(&user).RowsAffected
	if rows == 0 {
		return errors.New("Error in creating user")
	}
	return err
}

func (dc *DBRepo) Authenticate(email, password string) (error, int) {
	var user model.User
	rows := dc.GormDB.Debug().Where(`"email"=? and "password" =? `, email, password).First(&user).RowsAffected
	if rows == 1 {
		return nil, user.UserID
	}
	return errors.New("Invalid user"), 0

}
