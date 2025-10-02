package forms

type ShopCartItemForm struct {
	GoodsId int32 `form:"goods" json:"goods" binding:"required"`
	Nums    int32 `form:"nums" json:"nums" binding:"required"`
}

type ShopCartItemUpdateForm struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"`
}
