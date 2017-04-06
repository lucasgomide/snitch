package snitch

type Hook interface {
	CallHook(deploy []Deploy) error
	SetWebHookURL(url string)
}

type Tsuru interface {
	FindLastDeploy(deploy *[]Deploy) error
}

type Deploy struct {
	App       string
	Timestamp string
	Commit    string
	User      string
}
