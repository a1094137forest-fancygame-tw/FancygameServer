package model

type Ping struct {
	MessageCode int `json:"messageCode"`
	SystemTime  int `json:"systemTime,omitempty"`
}

type Pong struct {
	ResponseCode int `json:"responseCode"`
	SystemTime   int `json:"systemTime"`
}
