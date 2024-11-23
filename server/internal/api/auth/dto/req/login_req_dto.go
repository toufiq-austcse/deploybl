package req

type LoginReqDto struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}
