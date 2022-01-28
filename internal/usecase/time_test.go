package usecase

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todoNote/internal/model"
)

func TestTimeConvert(t *testing.T) {
	now := time.Now().UTC()
	format := "yyyy-MM-MM-dd hh:mm:ss"

	tts := []struct{
		in time.Time
		zone model.TimeZone
		out time.Time
	}{
		{now, model.UTC, now},
		{now, model.UTCp3, now.Add(3 * time.Hour)},
		{now, model.UTCm3, now.Add(-3 * time.Hour)},
	}

	for _, tt := range tts {
		assert.Equal(t, tt.out.Format(format), Convert(tt.in, tt.zone).Format(format))
	}
}

func TestValidateZone(t *testing.T) {
	tts := []struct{
		in string
		out bool
	}{
		{"UTC", true},
		{"UTC+1", true},
		{"UTC+3", true},
		{"UTC-6", true},
		{"UtC+1", false},
		{"utc+3", false},
		{"UTC-", false},
	}

	for _, tt := range tts {
		_, out := ValidateZone(tt.in)
		assert.Equal(t, tt.out, out)
	}
}