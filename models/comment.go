package models

import "time"

type Comment struct {
	ID            uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	WriterID      uint       `json:"writer_id" form:"writer_id"`
	WriterName    uint       `json:"writer_name" form:"writer_name"`
	Comment       string     `json"comment" form:"comment"`
	ParentID      uint       `json"parent_id" form:"parent_id"`
	ParentComment string     `json:"parent_comment" form:"parent_comment"`
	CreatedAt     *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" form:"updated_at"`
}
