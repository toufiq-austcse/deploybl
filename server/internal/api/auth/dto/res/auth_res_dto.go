package res

type SignUpResDto struct {
	Token    string       `json:"token"`
	UserInfo AuthUserInfo `json:"user_info"`
}

type LoginResModel struct {
	Token    string       `json:"token"`
	UserInfo AuthUserInfo `json:"user_info"`
}

type AuthUserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
