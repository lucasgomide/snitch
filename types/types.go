package snitch

import "time"

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

func (d Deploy) ConvertTimestampToRFC822() string {
	t, _ := time.Parse(time.RFC3339, d.Timestamp)
	return t.Format(time.RFC822)
}
