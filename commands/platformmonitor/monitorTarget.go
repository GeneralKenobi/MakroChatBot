package platformmonitor

type monitorTarget struct {
	TargetName string

	Listeners map[string]struct{}
}
