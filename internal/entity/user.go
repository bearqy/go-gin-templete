package entity

import "time"

type User struct {
	ID        int64     `xorm:"pk autoincr" json:"id"`
	Name      string    `xorm:"varchar(100) notnull" json:"name"`
	Email     string    `xorm:"varchar(255) unique" json:"email"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}
