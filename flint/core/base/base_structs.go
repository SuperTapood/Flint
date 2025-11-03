package base

import "time"

type Deployment struct {
	Name     string
	Age      time.Duration
	Status   string
	Revision int
}
