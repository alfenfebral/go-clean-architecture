package timeutil

import "time"

// GetTimeNow - get current time with timezone
func GetTimeNow() time.Time {
	// Init the loc
	loc, _ := time.LoadLocation("Asia/Jakarta")

	// Set timezone,
	return time.Now().In(loc)
}
