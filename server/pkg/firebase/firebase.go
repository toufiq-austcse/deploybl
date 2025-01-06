package firebase

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/pkg/firebase/api_res"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"google.golang.org/api/option"
)

type Client struct {
	AuthClient *auth.Client
	restyReq   *resty.Request
}

func NewFirebaseClient() (*Client, error) {
	opt := option.WithCredentialsFile(config.AppConfig.FIREBASE_CONFIG_FILE_PATH)
	newApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error in creating firebase app ", err.Error())
		return nil, err
	}
	authClient, err := newApp.Auth(context.Background())
	if err != nil {
		fmt.Println("error in getting firebase auth ", err.Error())
		return nil, err
	}
	return &Client{
		AuthClient: authClient,
		restyReq: resty.
			New().
			SetBaseURL(config.AppConfig.GOOGLE_API_BASE_URL).
			SetHeader("Content-Type", "application/json").
			R(),
	}, nil
}

func (client *Client) VerifyCustomToken(customToken string) (*api_res.VerifyCustomTokenRes, error) {
	var customTokenRes *api_res.VerifyCustomTokenRes

	response, err := client.restyReq.SetBody(map[string]interface{}{
		"token":             customToken,
		"returnSecureToken": true,
	}).SetResult(&customTokenRes).Post("/identitytoolkit/v3/relyingparty/verifyCustomToken?key=" + config.AppConfig.FIREBASE_API_KEY)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("error in verifying custom token")
	}
	return customTokenRes, nil
}
