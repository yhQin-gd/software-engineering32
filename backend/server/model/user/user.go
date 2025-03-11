package model

type User struct {
	ID         int    	`json:"id" gorm:"primarykey;autoIncrement"`
	Name       string 	`json:"name" gorm:"column:name; not null"`
	Email      string 	`json:"email" gorm:"unique;not null"`
	Password   string 	`json:"password" gorm:"not null"`
	RoleId	int 		`json:"role_id" gorm:"column:role_id;default:2"`	// 1: admin, 2: user
	IsVerified bool		`json:"is_verified" gorm:"column:is_verified"`
}
