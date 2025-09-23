package model

type Category struct {
	BaseModel
	Name             string    `gorm:"type:varchar(20);not null"`
	Level            int32     `gorm:"type:int;not null;default:1"`
	IsTab            bool      `gorm:"type:bool;default:false"`
	ParentCategoryID int32     `gorm:"type:int"`
	ParentCategory   *Category `gorm:"foreignKey:ParentCategoryID"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandID    int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands     Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;not null;default:1"`
}

func (Banner) TableName() string {
	return "banner"
}
