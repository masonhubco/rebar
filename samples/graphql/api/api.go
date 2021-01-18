package api

import (
	"encoding/json"
	"net/http"
	"time"
)

var appStartTime time.Time

func init() {
	appStartTime = time.Now()
}

func Status() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type status struct {
			Status string `json:"status"`
			Redis  string `json:"redis"`
			Uptime string `json:"uptime"`
		}

		s := status{
			Status: func() string {
				return "up"
			}(),
			Redis: func() string {
				return "connected"
			}(),
			Uptime: time.Since(appStartTime).Truncate(time.Second).String(),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s)
	}
}
