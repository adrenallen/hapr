package users

import (
	"encoding/json"
	"net/http"
)

type UserPassDTO struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (userPassObj UserPassDTO) ImportPayloadData(r *http.Request) (UserPassDTO, error) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userPassObj)
	return userPassObj, err
}
