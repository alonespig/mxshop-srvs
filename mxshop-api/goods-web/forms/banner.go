package forms

type BannerForm struct {
	Image string `form:"image" json:"image" binding:"required"`
	Index int32  `form:"index" json:"index" binding:"required"`
	Url   string `form:"url" json:"url" binding:"required"`
}
