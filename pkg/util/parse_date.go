package util

import (
	"fmt"
	"time"
)

func ParseDateTime(s string) (t time.Time, err error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
		"2006-01-02",
		"02 Jan 2006",
		"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05Z07:00", // RFC3339 without T
		"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000",
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		"02/01/2006 15:04:05",     // indonesian date time
		"02/01/2006 15:04:05.000", // indonesian date time
		"02/01/2006",              // indonesian date
	}

	for _, layout := range layouts {
		if t, err = time.Parse(layout, s); err == nil {
			return
		}
	}

	return t, fmt.Errorf("unable to parse date: %s", s)
}
