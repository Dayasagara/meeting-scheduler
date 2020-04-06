package model

import "time"

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
	StartSlot string `json:"startSlot"`
	EndSlot   string `json:"endSlot"`
}

type AvailabilitySlots struct {
	StartTime    time.Time `json:"-" gorm:"column:startTime"`
	Availability bool      `json:"availability" gorm:"column:availability"`
	StartSlot    string    `json:"startSlot"`
	EndSlot      string    `json:"endSlot"`
}

type GetAvParameters struct {
	UserID int    `json:"userID"`
	Date   string `json:"date"`
}

type ScheduleEvent struct {
	UserID       int    `json:"userID" gorm:"column:userID"`
	EventID      int    `json:"eventID" gorm:"column:eventID;primary_key;AUTO_INCREMENT"`
	StartingFrom string `json:"startingFrom" gorm:"column:starting_from"`
	EndingTill   string `json:"endingTill" gorm:"column:ending_till"`
	Date         string `json:"date" gorm:"column:date"`
	EventName    string `json:"eventName" gorm:"column:event_name"`
	Description  string `json:"description" gorm:"column:description"`
	ScheduledBy  int    `json:"scheduledBy" gorm:"column:scheduled_by"`
	Location     string `json:"location" gorm:"column:location"`
	Sync         bool   `json:"-" gorm:"column:sync"`
}
