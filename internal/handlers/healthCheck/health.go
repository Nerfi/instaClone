package healthCheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthCheck struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (*HealthCheck) Check(w http.ResponseWriter, r *http.Request) {
	response := HealthCheck{
		Status:    "OK",
		Timestamp: time.Now(),
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("error encoding response")
		panic(err)
	}

}
