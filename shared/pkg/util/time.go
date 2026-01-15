package util

import "time"

// Now returns the current time.
func Now() time.Time {
	return time.Now()
}

// Today returns today's date as YYYY-MM-DD.
func Today() string {
	return time.Now().Format("2006-01-02")
}

// NowRFC3339 returns the current time in RFC3339 format.
func NowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

// NowUnix returns the current Unix timestamp.
func NowUnix() int64 {
	return time.Now().Unix()
}

// ParseDate parses a YYYY-MM-DD date string.
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// ParseRFC3339 parses an RFC3339 timestamp string.
func ParseRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// DaysSince returns the number of days since a given time.
func DaysSince(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

// IsExpired checks if a time plus TTL days has passed.
func IsExpired(t time.Time, ttlDays int) bool {
	return DaysSince(t) > ttlDays
}

// FormatDuration formats a duration as human-readable string.
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	return d.Round(time.Hour).String()
}
