package models

type Users struct {
	BaseModel
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"not null;unique" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	Phone     string `json:"phone"`
	Username  string `gorm:"not null;unique" json:"username"`
	RoleID    int    `json:"role_id"`
}
