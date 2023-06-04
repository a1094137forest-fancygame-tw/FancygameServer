package model

type RegisterReq struct {
	MessageCode int64 `json:"messageCode"`
	Data        struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Avatar   int64  `json:"avatar"`
		Gender   int64  `json:"gender"`
	}
}

type ServerRespInfo struct {
	GameList []GameListInfo
	Avatar   int
	Gender   int
	IsAdmin  bool
	Balance  int
}

type GameListInfo struct {
	GameId int
	Status bool
}

type MemberListReq struct {
	MessageCode int `json:"messageCode"`
}

type MemberListRespInfo struct {
	MemberList []MemberInfo
}

type MemberInfo struct {
	Account        string
	Password       string
	Avatar         int
	Gender         int
	LastLoginTime  int
	LastLogoutTime int
	Balance        int
}

type MemberSettingReq struct {
	MessageCode int               `json:"messageCode"`
	Data        MemberSettingInfo `json:"data"`
}

type MemberSettingInfo struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Avatar   int    `json:"avatar"`
	Gender   int    `json:"gender"`
}

type LoginReq struct {
	MessageCode int `json:"messageCode"`
	Data        struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
}

type LoginResp struct {
	ResponseCode int            `json:"responseCode"`
	Message      string         `json:"message"`
	Data         ServerRespInfo `json:"data"`
}
