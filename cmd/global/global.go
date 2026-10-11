package global

var (
	Version, Commit, Date = "dev", "dev", ""

	CfgFile string
	NoColor bool
	NoStyle bool
	Verbose bool
	// Fresh starts without restoring the last session, see configuration.UiConfig.RestoreSession
	Fresh bool
)
