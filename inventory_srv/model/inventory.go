package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;not null;index"`
	Stocks  int32 `gorm:"type:int;not null"`
	Version int32 `gorm:"type:int;not null"` // 分布式锁的乐观锁
}

