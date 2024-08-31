package github

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github/api_res"
	"net/http"
	"strings"
)

type GithubHttpClient struct {
	restyReq *resty.Request
}

func NewGithubHttpClient() *GithubHttpClient {
	return &GithubHttpClient{
		restyReq: resty.
			New().
			SetBaseURL(config.AppConfig.GITHUB_API_BASE_URL).
			R().
			SetHeader("Content-Type", "application/json"),
	}
}

func (httpClient *GithubHttpClient) ValidateRepository(repoOwner, repoName *string) (*api_res.GithubRepoRes, int, error) {
	var githubRes *api_res.GithubRepoRes

	getRes, err := httpClient.restyReq.SetResult(&githubRes).Get("repos/" + *repoOwner + "/" + *repoName)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if getRes == nil {
		return nil, http.StatusInternalServerError, errors.New("error in repo validation")
	}
	if getRes.IsSuccess() {
		return githubRes, getRes.StatusCode(), nil
	}

	var githubErrorRes *api_res.GithubErrorRes
	unmarshalErr := json.Unmarshal(getRes.Body(), &githubErrorRes)
	if unmarshalErr != nil {
		return nil, http.StatusInternalServerError, unmarshalErr
	}
	return nil, getRes.StatusCode(), errors.New(githubErrorRes.Message)

}
func (httpClient *GithubHttpClient) ValidateRepositoryByUrl(repoUrl *string) (*api_res.GithubRepoRes, int, error) {
	repoName, owner, err := httpClient.ParseRepoUrl(repoUrl)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return httpClient.ValidateRepository(owner, repoName)

}

func (httpClient *GithubHttpClient) ParseRepoUrl(repoUrl *string) (repoOwner *string, repoName *string, err error) {
	urlWithoutHttps := strings.TrimLeft(*repoUrl, "https://")
	parts := strings.Split(urlWithoutHttps, "/")
	if len(parts) != 3 {
		return nil, nil, errors.New("invalid repo url")
	}
	name := strings.TrimRight(parts[len(parts)-1], ".git")
	owner := parts[len(parts)-2]
	return &name, &owner, err

}
