package model

import "time"

type Meta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}

type Document struct {
	ID        string    `json:"id,omitempty" gorm:"primaryKey;autoIncrement;column:id"`
	Name      string    `json:"name,omitempty" gorm:"column:name"`
	FileData  []byte    `json:"file_data,omitempty" gorm:"column:file_data"`
	Public    bool      `json:"public,omitempty" gorm:"column:public"`
	File      bool      `json:"file,omitempty" gorm:"column:file"`
	Mime      string    `json:"mime,omitempty" gorm:"column:mime"`
	Token     string    `json:"token,omitempty" gorm:"column:token"`
	Grant     string    `json:"grant,omitempty" gorm:"column:grant"`
	JSONData  string    `json:"json_data,omitempty" gorm:"column:json_data"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime;column:created_at"`
}
