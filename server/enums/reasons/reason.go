package reasons

const (
	ENV_UPDATED             string = "Env update"
	TRIGGERED_VIA_DASHBOARD string = "Triggered via dashboard"
	SETTINGS_UPDATE         string = "Settings Update"
)

func GetReasonPtr(reason string) *string {
	return &reason
}
