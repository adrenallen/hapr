package helpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
)

//Reads the body of a request and parses as JSON
//Returns the gabs container of value
//Note this will use the buffer from the body and render is gone forever
func GetGabsContainerFromRequest(r *http.Request) (*gabs.Container, error) {
	jsonString, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return GetGabsContainerFromBytes(jsonString)
}

func GetGabsContainerFromBytes(b []byte) (*gabs.Container, error) {
	requestGabsContainer, err := gabs.ParseJSON(b)
	if err != nil {
		return nil, err
	}
	return requestGabsContainer, nil
}

func GetValueFromRequestPayload(r *gabs.Container, key string) (string, error) {
	valContainer, err := GetContainerFromRequestPayload(r, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", valContainer.Data()), nil
}

func GetBoolValueFromRequestPayload(r *gabs.Container, key string) (bool, error) {
	valContainer, err := GetContainerFromRequestPayload(r, key)
	if err != nil {
		return false, err
	}
	return valContainer.Data().(bool), nil
}

func GetIntValueFromRequestPayload(r *gabs.Container, key string) (int, error) {
	valContainer, err := GetContainerFromRequestPayload(r, key)
	if err != nil {
		return 0, err
	}
	return int(valContainer.Data().(float64)), nil
}

func GetRawValueFromRequestPayload(r *gabs.Container, key string) (interface{}, error) {
	valContainer, err := GetContainerFromRequestPayload(r, key)
	if err != nil {
		return nil, err
	}
	return valContainer.Data(), nil
}

func GetContainerFromRequestPayload(r *gabs.Container, key string) (*gabs.Container, error) {
	// val, ok := r.Path(key).Data()
	ok := r.ExistsP(key)
	if !ok {
		return nil, errors.New("request payload helper: value not found in request")
	}
	return r.Path(key), nil
}
