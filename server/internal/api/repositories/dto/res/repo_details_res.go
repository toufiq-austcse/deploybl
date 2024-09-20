package res

type RepoDetailsRes struct {
	SvnUrl        string `json:"svn_url"`
	DefaultBranch string `json:"default_branch"`
	Name          string `json:"name"`
}
