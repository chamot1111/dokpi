package dokpi

type appStatus struct {
	name    string
	running bool
}

type msg struct {
	cmdSuccess bool
	errorMsg   string
	appsStatus []appStatus
}

func (m *msg) setError(msg string) {
	m.cmdSuccess = false
	m.errorMsg = msg
}
