package mapper

import (
	"github.com/toufiq-austcse/deployit/internal/api/repositories/dto/res"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github/api_res"
)

func ToRepoDetailsRes(providerRes *api_res.GithubRepoRes) res.RepoDetailsRes {
	return res.RepoDetailsRes{
		GitUrl:        providerRes.GitURL,
		DefaultBranch: providerRes.DefaultBranch,
		FullName:      providerRes.FullName,
	}

}
