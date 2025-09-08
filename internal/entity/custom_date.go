package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CustomDate time.Time

const customDateFormat = "01-2006"

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		*cd = CustomDate(time.Time{})
		return nil
	}
	t, err := time.Parse(customDateFormat, s)
	if err != nil {
		// Fallback to standard time format
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			// Fallback to date-only format
			t, err = time.Parse("2006-01-02", s)
			if err != nil {
				return fmt.Errorf("failed to parse time in any supported format: %s", s)
			}
		}
	}
	*cd = CustomDate(t)
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if time.Time(cd).IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(time.Time(cd).Format(customDateFormat))
}

func (cd *CustomDate) Scan(value any) error {
	if value == nil {
		*cd = CustomDate(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*cd = CustomDate(v)
		return nil
	}
	return fmt.Errorf("can't convert %T to CustomDate", value)
}

func (cd CustomDate) Value() (driver.Value, error) {
	if time.Time(cd).IsZero() {
		return nil, nil
	}
	return time.Time(cd), nil
}
