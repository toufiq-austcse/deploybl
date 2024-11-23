package mapper

import (
	"firebase.google.com/go/v4/auth"
	"github.com/toufiq-austcse/deployit/internal/api/auth/dto/res"
)

func ToSignUpResDto(token string, userInfo auth.UserInfo) res.SignUpResDto {
	return res.SignUpResDto{
		Token: token,
		UserInfo: res.AuthUserInfo{
			Name:  userInfo.DisplayName,
			Email: userInfo.Email,
		},
	}
}
