package helpers

import "time"

func GenerateTime(month int, day int, year int, locationString string) time.Time {
	locationObj, err := time.LoadLocation(locationString)
	if err != nil {
		locationObj = time.UTC //default to UTC if we can't do it
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, locationObj)
}
