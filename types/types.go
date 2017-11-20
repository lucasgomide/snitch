package types

import "time"

// Hook represents an interface for any hooks
type Hook interface {
	CallHook(deploy []Deploy) error
	ValidatesFields() error
}

// Tsuru represents an interface for tsuru
type Tsuru interface {
	FindLastDeploy(deploy *[]Deploy) error
}

// Deploy represents a struct for deploy
type Deploy struct {
	App       string
	Timestamp string
	Commit    string
	User      string
	Image     string
}

// ConvertTimestampToRFC822 converts the Timestamp's deploy
// to time.RFC822 date format
func (d Deploy) ConvertTimestampToRFC822() string {
	t, _ := time.Parse(time.RFC3339, d.Timestamp)
	return t.Format(time.RFC822)
}
