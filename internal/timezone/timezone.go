package timezone

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ErrInvalidTimezone is returned when the timezone is invalid.
var ErrInvalidTimezone = errors.New("invalid timezone")

// GetTimezoneDiff returns the difference between the local time and the specified hour.
func GetTimezoneDiff(localHour, needHour int) int {
	timeZoneDiff := calculateTimeZoneDiff(localHour, needHour)
	return adjustTimeZoneDiff(timeZoneDiff)
}

func calculateTimeZoneDiff(localHour, needHour int) int {
	return needHour - localHour
}

// adjustTimeZoneDiff adjusts the time difference if it exceeds 12 hours in either direction.
func adjustTimeZoneDiff(timeZoneDiff int) int {
	if timeZoneDiff > 12 {
		return timeZoneDiff - 24
	} else if timeZoneDiff < -12 {
		return 24 + timeZoneDiff
	}
	return timeZoneDiff
}

// ProcessTimezoneHeader processes the timezone header and returns the offset.
func ProcessTimezoneHeader(r *http.Request) (int, error) {
	userTimezone := r.Header.Get("Accept-Timezone")
	if userTimezone == "" {
		return 0, nil
	}

	offsetStr, err := extractOffsetString(userTimezone)
	if err != nil {
		return 0, fmt.Errorf("extracting offset string: %w", err)
	}

	offset, err := parseOffset(offsetStr)
	if err != nil {
		return 0, fmt.Errorf("parsing offset: %w", err)
	}

	return offset, nil
}

func extractOffsetString(userTimezone string) (string, error) {
	switch {
	case userTimezone == "UTC":
		return "0", nil
	case len(userTimezone) > 4 &&
		len(userTimezone) < 7 &&
		(strings.HasPrefix(userTimezone, "UTC+") || strings.HasPrefix(userTimezone, "UTC-")):
		return userTimezone[3:], nil
	default:
		return "", ErrInvalidTimezone
	}
}

func parseOffset(offsetStr string) (int, error) {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset > 12 || offset < -12 {
		return 0, ErrInvalidTimezone
	}
	return offset, nil
}

func ValidateTimezoneDiff(diff int) error {
	if diff > 12 || diff < -12 {
		return ErrInvalidTimezone
	}
	return nil
}
