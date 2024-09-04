package res

type RepoDetailsRes struct {
	GitUrl        string `json:"git_url"`
	DefaultBranch string `json:"default_branch"`
	FullName      string `json:"full_name"`
}
