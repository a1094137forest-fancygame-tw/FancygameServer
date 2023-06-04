package constant

import (
	p "github.com/golang/protobuf/proto"
	"gopkg.in/olahol/melody.v1"

	"BackendServer/proto/common"
)

const (
	SUCCESS = 201
)

// common
const (
	RESPONSE_PONG  = 1
	RESPONSE_ERROR = 2

	RESPONSE_REGISTER_RES = iota + 41
	RESPONSE_LOGIN_RES
	RESPONSE_LOGOUT_RES
)

const (
	EVENT_KICKOUT = iota + 71
)

// lobby
const (
	RESPONSE_GET_MEMBER_LIST_RES = iota + 151
	RESPONSE_GET_GAME_LIST_RES
	RESPONSE_KICKOUT_MEMBER_RES
)

const (
	ERROR_OPERATION_TIMEOUT = iota + 501
	ERROR_OPERATION_FAILED
	ERROR_INVALID_OPERATION
	ERROR_KICKOUT
)

const (
	EVENT_UPDATE_MEMBER_DATA = iota + 201
)

var MsgFlags = map[int]string{
	RESPONSE_PONG: "Pong",

	ERROR_OPERATION_TIMEOUT: "Failed Time out",
	ERROR_OPERATION_FAILED:  "Failed operation",
	ERROR_INVALID_OPERATION: "Failed invalid operation",

	EVENT_KICKOUT: "Failed connect",
}

func ResponseWithData(s *melody.Session, method int, data *[]byte) {
	result := append([]byte{byte(0), byte(method)}, *data...)
	_ = s.WriteBinary(result)
}

func ResponseWithError(s *melody.Session, statusCode int, err string) {
	var data = common.Error{
		StatusCode: int64(statusCode),
		Message:    err,
	}

	byteData, _ := p.Marshal(&data)
	byteData = append([]byte{byte(0), byte(RESPONSE_ERROR)}, byteData...)
	_ = s.WriteBinary(byteData)
}

func ResponseWithKickOut(s *melody.Session, statusCode int, err string) {
	var data = common.KickOut{
		StatusCode: int64(statusCode),
		Message:    err,
	}

	byteData, _ := p.Marshal(&data)
	byteData = append([]byte{byte(0), byte(EVENT_KICKOUT)}, byteData...)
	_ = s.WriteBinary(byteData)
}
