package transformer

import "time"

func dateToString(time time.Time) *string {
	if time.IsZero() {
		return nil
	}
	str := time.String()
	return &str
}
