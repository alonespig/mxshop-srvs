package model

type Category struct {
	BaseModel
	Name             string    `gorm:"type:varchar(20);not null"`
	Level            int32     `gorm:"type:int;not null;default:1"`
	IsTab            bool      `gorm:"type:bool;default:false"`
	ParentCategoryID int32     `gorm:"type:int"`
	ParentCategory   *Category `gorm:"foreignKey:ParentCategoryID"`
}
