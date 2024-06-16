package timezone_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/timezone"
	"github.com/stretchr/testify/require"
)

func TestProcessTimezoneHeader(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    int
		wantErr error
	}{
		{"valid timezone", "UTC+10", 10, nil},
		{"valid timezone", "UTC+3", 3, nil},
		{"no header", "", 0, nil},
		{"invalid timezone - more than 12 hours", "UTC+13", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - less than -12 hours", "UTC-13", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - not in the format UTC+/-<number>", "UTC+ab", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - not in the format UTC+/-<number>", "UTC+3:00", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - not in the format UTC+/-<number>", "UTC+", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - not in the format UTC+/-<number>", "UTC", 0, nil},
		{"invalid timezone - not in the format UTC+/-<number>", "UTC-", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - not in the format UTC+/-<number>", "ABC", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - length is more than 6", "UTC+1233", 0, customerrors.ErrInvalidTimezone},
		{"invalid timezone - length is less than 3", "UT", 0, customerrors.ErrInvalidTimezone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", http.NoBody)
			r.Header.Set("Accept-Timezone", tt.header)
			offset, err := timezone.ProcessTimezoneHeader(r)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, offset)
			}
		})
	}
}

func TestGetTimezoneDiff(t *testing.T) {
	tests := []struct {
		name     string
		local    int
		need     int
		wantDiff int
	}{
		{"same time", 10, 10, 0},
		{"need is ahead", 10, 12, 2},
		{"need is behind", 10, 8, -2},
		{"need is ahead - midnight", 23, 1, 2},
		{"need is behind - midnight", 1, 23, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := timezone.GetTimezoneDiff(tt.local, tt.need)
			require.Equal(t, tt.wantDiff, diff)
		})
	}
}

func TestValidateTimezoneDiff(t *testing.T) {
	tests := []struct {
		name    string
		diff    int
		wantErr error
	}{
		{"valid diff", 0, nil},
		{"valid diff", 12, nil},
		{"valid diff", -12, nil},
		{"invalid diff - more than 12 hours", 13, customerrors.ErrInvalidTimezone},
		{"invalid diff - less than -12 hours", -13, customerrors.ErrInvalidTimezone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := timezone.ValidateTimezoneDiff(tt.diff)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
