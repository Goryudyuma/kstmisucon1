package models

import "time"

type User struct {
	ID         uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	ScreenName string     `json:"screen_name" form:"name"`
	UserName   string     `json:"user_name" form:"emails"`
	Password   string     `json:"password" form:"password"`
	CreatedAt  *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" form:"updated_at"`
}
