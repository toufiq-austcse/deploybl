package reasons

const (
	ENV_UPDATED        string = "Env updated"
	RESTART_DEPLOYMENT string = "Restart deployment"
	REBUILD_DEPLOYMENT string = "Rebuild deployment"
	STOP_DEPLOYMENT    string = "Stop deployment"
)

func GetReasonPtr(reason string) *string {
	return &reason
}
