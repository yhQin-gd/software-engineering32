package model

type User struct {
	ID         int    	`json:"id" gorm:"primarykey;autoIncrement"`
	Name       string 	`json:"uname" gorm:"column:name; not null"`
	Email      string 	`json:"email" gorm:"unique;not null"`
	Password   string 	`json:"password" gorm:"not null"`
	RoleId	int 		`json:"roleid" gorm:"column:roleid;default:2"`	// 1: admin, 2: user
	IsVerified bool		`json:"isverified" gorm:"column:isverified"`
}
