package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Dayasagara/meeting-scheduler/model"
	"github.com/jinzhu/gorm"
)

var tx *gorm.DB

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

//Authenticate user with email and password
func (dc *DBRepo) Authenticate(email, password string) (error, int) {
	var user model.User
	rows := dc.GormDB.Debug().Where(`"email"=? and "password" =? `, email, password).First(&user).RowsAffected
	if rows == 1 {
		return nil, user.UserID
	}
	return errors.New("Invalid user"), 0

}

//Check if slots are defined for a user on a given day
func (dc *DBRepo) CheckForDuplicate(userID int, date string) bool {
	var availabilities []model.Availability
	dc.GormDB.Debug().Table("availabilities").Where(`"userID" = ? and "date" = ?`, userID, date).Find(&availabilities)
	if len(availabilities) > 0 {
		return true
	}
	return false
}

//Define slots for a user for a given day
func (dc *DBRepo) DefineAvailability(availability model.Availability) error {
	tx = dc.GormDB.Begin()
	query := `INSERT INTO "availabilities" ("userID", "date", "startTime", "availability") VALUES`
	var values []string
	var startTime string

	endSlot, endConvErr := strconv.Atoi(strings.Split(availability.EndSlot, ":")[0])
	startSlot, startConvErr := strconv.Atoi(strings.Split(availability.StartSlot, ":")[0])
	if endConvErr != nil || startConvErr != nil {
		return errors.New("Time conversion error")
	}
	for startSlot < endSlot {
		startTime = strconv.Itoa(startSlot) + ":00:00"
		values = append(values, fmt.Sprintf(`('%d','%s', '%s', '%v')`,
			availability.UserID, availability.Date, startTime, true))
		startSlot++
	}
	if err := tx.Debug().Exec(query + strings.Join(values, ",")).Error; err != nil {
		dc.rollbackTransaction()
		return err
	}
	tx.Commit()
	return nil
}

func (dc *DBRepo) GetAvailability(userID int, date string) ([]model.AvailabilitySlots, error) {
	var avSlots []model.AvailabilitySlots
	err := dc.GormDB.Debug().Table("availabilities").Where(`"userID" = ? and "date" = ?`, userID, date).Find(&avSlots).Error
	return avSlots, err
}

//Check for availability of slot
func (dc *DBRepo) CheckAvailability(userID int, date string, from string) bool {
	var avSlots []model.AvailabilitySlots
	rows := dc.GormDB.Debug().Table("availabilities").Where(`"userID" = ? and "date" = ? and "startTime" = ? and "availability" = ?`, userID, date, from, true).Find(&avSlots).RowsAffected
	if rows >= 1 {
		return true
	}
	return false
}

//update the availability table and schedule an event
func (dc *DBRepo) ScheduleEvent(event model.ScheduleEvent) error {
	var err error
	tx = dc.GormDB.Begin()
	err = tx.Debug().Table("availabilities").Where(`"userID" = ? and date = ? and "startTime" = ?`, event.UserID, event.Date, event.StartingFrom).
		Updates(map[string]interface{}{"availability": false}).Error

	if err != nil {
		return err
	}

	rows := tx.Debug().Table("scheduled_events").Create(&event).RowsAffected
	if rows == 0 {
		dc.rollbackTransaction()
		return errors.New("Error in creating Event")
	}
	tx.Commit()
	return err
}

func (dc *DBRepo) GetEvents(userID int) ([]model.ScheduleEvent, error) {
	var events []model.ScheduleEvent
	err := dc.GormDB.Debug().Table("scheduled_events").Where(`"userID" = ? and "sync" = ?`, userID, false).Find(&events).Error
	if err != nil {
		return nil, err
	}
	err = dc.GormDB.Debug().Table("scheduled_events").Where(`"userID" = ? and sync = ?`, userID, false).
		Updates(map[string]interface{}{"sync": true}).Error

	if err != nil {
		return nil, err
	}
	return events, err
}

func (dc *DBRepo) rollbackTransaction() {
	tx.Rollback()
}
