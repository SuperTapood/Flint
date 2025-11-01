package base

import "time"

type Deployment struct {
	Name     string
	Duration time.Duration
	Status   string
	Revision int
}
