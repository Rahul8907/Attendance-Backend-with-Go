package models

import "time"

type User struct {
	ID         int64     `xorm:"pk autoincr"`
	Name       string    `xorm:"varchar(100) notnull"`
	Grade      string    `xorm:"varchar(2) notnull"`
	LogInTime  time.Time `xorm:"logintime"`
	LogOutTime time.Time `xorm:"logouttime"`
}
