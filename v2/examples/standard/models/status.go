package models

import (
	"runtime"
	"time"
)

type Status struct {
	StartTime time.Time `json:"-"`
	Uptime    string    `json:"uptime"`
	Commit    string    `json:"commit"`
	Built     string    `json:"built"`
	GoVersion string    `json:"-"`
}

func NewStatus(commit, built string) Status {
	return Status{
		StartTime: time.Now(),
		Commit:    commit,
		Built:     built,
		GoVersion: runtime.Version(),
	}
}

func (s Status) Snapshot() Status {
	s.Uptime = time.Since(s.StartTime).Truncate(time.Second).String()
	return s
}
