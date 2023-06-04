package constant

import (
	"encoding/json"

	"gopkg.in/olahol/melody.v1"
)

// common
const (
	REQUEST_PING = 0

	REQUEST_REGISTER = iota + 10
	REQUEST_LOGIN
	REQUEST_LOGOUT
)

// lobby
const (
	REQUEST_MEMBER_LIST = iota + 101
	REQUEST_SET_MEMBER
)

const (
	METHOD_REGISTER = iota + 151
)

type RequestMessage struct {
	MessageCode int         `json:"messageCode"`
	Data        interface{} `json:"data"`
}

func RequestWithData(s *melody.Session, messageCode int, data interface{}) error {

	req := RequestMessage{
		MessageCode: messageCode,
		Data:        data,
	}

	jsonD, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = s.Write(jsonD)
	if err != nil {
		return err
	}
	return nil
}
