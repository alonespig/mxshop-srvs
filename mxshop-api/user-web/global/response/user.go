package response

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	BirthDay string `json:"birthday"`
}
