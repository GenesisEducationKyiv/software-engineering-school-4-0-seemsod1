package helpers

import (
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"net/http"
	"strconv"
	"time"
)

// NewRateResponse is a helper function that creates a new RateResponse struct
func NewRateResponse(price float64) interface{} {
	type RateResponse struct {
		Price float64 `json:"price"`
	}
	return RateResponse{Price: price}
}

// GetTimezoneDiff is a helper function that returns the difference between the local time and the time needed
func GetTimezoneDiff(needHour int) int {
	localHour := time.Now().Hour()
	timeZoneDiff := needHour - localHour
	if timeZoneDiff > 12 || timeZoneDiff < -12 {
		if timeZoneDiff < 0 {
			timeZoneDiff *= -1
		}
		if localHour > needHour {
			return 24 - timeZoneDiff
		} else if localHour < needHour {
			return timeZoneDiff - 24
		} else {
			return 0
		}
	} else {
		return timeZoneDiff
	}
}

// ProcessTimezoneHeader is a helper function that processes the timezone header
func ProcessTimezoneHeader(r *http.Request) (int, error) {
	userTimezone := r.Header.Get("Accept-Timezone")
	var offsetStr string
	var offset int
	if userTimezone != "" {
		if userTimezone == "UTC" {
			offsetStr = "0"
		} else if len(userTimezone) == 5 {
			offsetStr = userTimezone[3:]
		} else {
			return 0, customerrors.InvalidTimezone
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset > 12 || offset < -12 {
			return 0, customerrors.InvalidTimezone
		}
	} else {
		offset = 0
	}
	return offset, nil
}
