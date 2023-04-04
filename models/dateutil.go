package models

import (
	"time"
)

func ConvertToRFC3339(timestamp string) (string, error) {
	// Parse the timestamp string using PostgreSQL timestamp layout
	t, err := time.Parse("2006-01-02 15:04:05.999999999", timestamp)
	if err != nil {
		return "", err
	}

	// Convert the time to RFC3339 format
	rfc3339 := t.Format(time.RFC3339)
	return rfc3339, nil
}
