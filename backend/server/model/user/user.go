package model

type User struct {
	ID         int    	`json:"id" gorm:"primarykey;autoIncrement"`
	Name       string 	`json:"uname" gorm:"column:name; not null"`
	Email      string 	`json:"email" gorm:"unique;not null"`
	Password   string 	`json:"password" gorm:"not null"`
	IsVerified bool		`json:"isverified" gorm:"column:isverified"`
}
