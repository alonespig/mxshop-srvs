package forms

type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"`
	PassWord string `form:"password" json:"password" binding:"required"`
}
