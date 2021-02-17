package dtos

import (
	"fmt"
	"time"

	"github.com/Jeffail/gabs"
)

type DateDTO struct {
	Day        int    `json:"day"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
	Nanosecond int    `json:"nanosecond"`
	Timezone   string `json:"timezone"`
}

func (d *DateDTO) ConvertToTime() time.Time {
	locObj, err := time.LoadLocation(d.Timezone)
	if err != nil {
		locObj = time.UTC
	}

	return time.Date(d.Year, time.Month(d.Month), d.Day, d.Hour, d.Minute, d.Second, d.Nanosecond, locObj)
}

func (d *DateDTO) LoadDateFromMap(container *gabs.Container) error {

	if err := getDataFromContainerToInt(container, &d.Day, "day"); err != nil {
		return err
	}

	if err := getDataFromContainerToInt(container, &d.Year, "year"); err != nil {
		return err
	}

	if err := getDataFromContainerToInt(container, &d.Month, "month"); err != nil {
		return err
	}

	//We dont require the following so ignore errors
	getDataFromContainerToInt(container, &d.Hour, "hour")

	getDataFromContainerToInt(container, &d.Minute, "minute")

	getDataFromContainerToInt(container, &d.Second, "second")

	getDataFromContainerToInt(container, &d.Nanosecond, "nanosecond")

	err := getDataFromContainerToString(container, &d.Timezone, "timezone")
	if err != nil {
		//we don't care just set to UTC
		d.Timezone = time.UTC.String()
	}
	return nil

}

func getDataFromContainerToInt(container *gabs.Container, attribute *int, key string) error {
	ok := container.ExistsP(key)
	if !ok {
		return fmt.Errorf("key was not found %v", key)
	}

	valConverted, ok := container.Path(key).Data().(float64)

	if !ok {
		return fmt.Errorf("Value found at %s was not a valid int", key)
	}

	*attribute = int(valConverted)
	return nil
}

func getDataFromContainerToString(container *gabs.Container, attribute *string, key string) error {
	ok := container.ExistsP(key)
	if !ok {
		return fmt.Errorf("key was not found %v", key)
	}

	valConverted, ok := container.Path(key).Data().(string)

	if !ok {
		return fmt.Errorf("Value found at %s was not a valid string", key)
	}
	*attribute = valConverted
	return nil
}
