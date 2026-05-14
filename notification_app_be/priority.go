// priority.go
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"time"
)

// Notification matches the Test Server API Response
type NotificationItem struct {
	ID        string `json:"ID"`
	Type      string `json:"Type"`
	Message   string `json:"Message"`
	Timestamp string `json:"Timestamp"`
}

// Weight mapping
var priorityWeights = map[string]int{
	"Placement": 3,
	"Result":    2,
	"Event":     1,
}

// FetchNotifications calls the external evaluation server
func FetchNotifications(accessToken string) ([]NotificationItem, error) {
	req, err := http.NewRequest("GET", "http://4.224.186.213/evaluation-service/notifications", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Notifications []NotificationItem `json:"notifications"`
	}
	err = json.Unmarshal(body, &result)
	return result.Notifications, err
}

// GetTopPriority filters and sorts to find the top N
func GetTopPriority(notifications []NotificationItem, n int) []NotificationItem {

	sort.Slice(notifications, func(i, j int) bool {
		weightI := priorityWeights[notifications[i].Type]
		weightJ := priorityWeights[notifications[j].Type]

		if weightI != weightJ {
			return weightI > weightJ
		}
		timeI, err1 := time.Parse(time.RFC3339, notifications[i].Timestamp)
		if err1 != nil {
			timeI, _ = time.Parse("2006-01-02 15:04:05", notifications[i].Timestamp)
		}
		timeJ, err2 := time.Parse(time.RFC3339, notifications[j].Timestamp)
		if err2 != nil {
			timeJ, _ = time.Parse("2006-01-02 15:04:05", notifications[j].Timestamp)
		}

		return timeI.After(timeJ)
	})

	// Return top N
	if len(notifications) > n {
		return notifications[:n]
	}
	return notifications
}
