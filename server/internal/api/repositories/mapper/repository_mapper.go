package mapper

import (
	"github.com/toufiq-austcse/deployit/internal/api/repositories/dto/res"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github/api_res"
)

func ToRepoDetailsRes(providerRes *api_res.GithubRepoRes) res.RepoDetailsRes {
	return res.RepoDetailsRes{
		SvnUrl:        providerRes.SvnURL,
		DefaultBranch: providerRes.DefaultBranch,
		Name:          providerRes.Name,
	}
}

func ToRepoBranchesRes(githubBranchesRes []api_res.GithubGetBranchesAPIRes) []res.RepoBranchRes {
	branches := []res.RepoBranchRes{}
	for _, branch := range githubBranchesRes {
		branches = append(branches, res.RepoBranchRes{
			Name: branch.Name,
		})
	}
	return branches
}
