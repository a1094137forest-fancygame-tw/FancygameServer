package v1

import (
	"context"
	"time"

	p "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"gopkg.in/olahol/melody.v1"

	"BackendServer/config"
	"BackendServer/constant"
	"BackendServer/proto/backend/fancygame"
	"BackendServer/proto/common"
)

func Error(c *Client, msg []byte) {
	data := common.Error{
		StatusCode: constant.ERROR_INVALID_OPERATION,
		Message:    constant.MsgFlags[constant.ERROR_INVALID_OPERATION],
	}

	byteData, _ := p.Marshal(&data)
	constant.ResponseWithData(c.Session, constant.RESPONSE_ERROR, &byteData)
}

func Pong(c *Client, msg []byte) {
	resp := &common.Pong{
		StatusCode: constant.SUCCESS,
		Message:    constant.MsgFlags[constant.RESPONSE_PONG],
		Data: &common.PongInfo{
			TimeStamp: int64(time.Now().Unix()),
		},
	}

	byteMsg, _ := p.Marshal(resp)

	constant.ResponseWithData(c.Session, constant.RESPONSE_PONG, &byteMsg)
}

func Register(c *Client, msg []byte) {
	var req common.Register

	if err := p.Unmarshal(msg, &req); err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_INVALID_OPERATION, constant.MsgFlags[constant.ERROR_INVALID_OPERATION])
		return
	}

	byteMsg, _ := p.Marshal(&req)

	var userReq fancygame.RegisterInfo

	_ = p.Unmarshal(byteMsg, &userReq)

	conn, err := grpc.Dial(config.UserServerUrl, grpc.WithInsecure())
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	client := fancygame.NewUserClient(conn)

	userResp, err := client.Register(context.Background(), &userReq)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	/*
		conn, err := grpc.Dial(config.GameServerUrl, grpc.WithInsecure())
		if err != nil {
			constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
			return
		}

		client := backendGame.NewGameClient(conn)

		gameResp, err := client.GetGameList(context.Background(), &common.EmptyReq{})
		if err != nil {
			constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
			return
		}
	*/

	c.Account = req.Account
	c.Avatar = int32(req.Avatar)
	c.Gender = int32(req.Gender)
	c.IsAdmin = false
	c.Balance = 0

	err = ConnPool.UpsertClient(c, userResp.Data.IsAdmin)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	var resp = common.RegisterRes{
		StatusCode: userResp.StatusCode,
		Message:    userResp.Message,
		Data: &common.RegisterResInfo{
			//GameList: gameResp.Data.GameList,
			GameList: []*common.GameList{},
			Avatar:   common.AvatarEnum(userResp.Data.Avatar),
			Gender:   common.GenderEnum(userResp.Data.Gender),
			IsAdmin:  userResp.Data.IsAdmin,
			Balance:  userResp.Data.Balance,
		},
	}

	byteMsg, _ = p.Marshal(&resp)

	constant.ResponseWithData(c.Session, constant.RESPONSE_REGISTER_RES, &byteMsg)
}

func Login(c *Client, msg []byte) {
	var req common.Login

	if err := p.Unmarshal(msg, &req); err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_INVALID_OPERATION, constant.MsgFlags[constant.ERROR_INVALID_OPERATION])
		return
	}

	account := req.Account
	password := req.Password

	validate := constant.ValidateString(account, true, true, true, constant.MAX_ACCOUNT, constant.MIN_ACCOUNT)
	if !validate {
		constant.ResponseWithError(c.Session, constant.ERROR_INVALID_OPERATION, constant.MsgFlags[constant.ERROR_INVALID_OPERATION])
		return
	}
	validate = constant.ValidateString(password, true, true, true, constant.MAX_PASSWORD, constant.MIN_PASSWORD)
	if !validate {
		constant.ResponseWithError(c.Session, constant.ERROR_INVALID_OPERATION, constant.MsgFlags[constant.ERROR_INVALID_OPERATION])
		return
	}

	var userReq = fancygame.LoginInfo{
		Account:  req.Account,
		Password: req.Password,
	}

	conn, err := grpc.Dial(config.UserServerUrl, grpc.WithInsecure())
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}
	client := fancygame.NewUserClient(conn)

	userResp, err := client.Login(context.Background(), &userReq)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	/*
		conn, err := grpc.Dial(config.GameServerUrl, grpc.WithInsecure())
		if err != nil {
			constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
			return
		}

		client := backendGame.NewGameClient(conn)

		gameResp, err := client.GetGameList(context.Background(), &common.EmptyReq{})
		if err != nil {
			constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
			return
		}
	*/

	c.Account = userReq.Account
	c.Avatar = int32(userResp.Data.Avatar)
	c.Gender = int32(userResp.Data.Gender)
	c.IsAdmin = userResp.Data.IsAdmin
	c.Balance = userResp.Data.Balance

	var resp = common.LoginRes{
		StatusCode: userResp.StatusCode,
		Message:    userResp.Message,
		Data: &common.LoginResInfo{
			//GameList: gameResp.Data.GameList,
			GameList: []*common.GameList{},
			Avatar:   common.AvatarEnum(userResp.Data.Avatar),
			Gender:   common.GenderEnum(userResp.Data.Gender),
			IsAdmin:  userResp.Data.IsAdmin,
			Balance:  userResp.Data.Balance,
		},
	}

	err = ConnPool.UpsertClient(c, userResp.Data.IsAdmin)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	byteMsg, _ := p.Marshal(&resp)

	constant.ResponseWithData(c.Session, constant.RESPONSE_REGISTER_RES, &byteMsg)
}

func Logout(c *Client, msg []byte) {
	conn, err := grpc.Dial(config.UserServerUrl, grpc.WithInsecure())
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}
	client := fancygame.NewUserClient(conn)

	userReq := fancygame.LogoutInfo{
		Account: c.Account,
		Balance: c.Balance,
	}

	userResp, err := client.Logout(context.Background(), &userReq)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	var logoutRes = common.LogoutRes{
		StatusCode: userResp.StatusCode,
		Message:    userResp.Message,
	}

	byteMsg, _ := p.Marshal(&logoutRes)

	constant.ResponseWithData(c.Session, constant.RESPONSE_REGISTER_RES, &byteMsg)
	ConnPool.DeleteClient(c, c.IsAdmin)
}

func KickOut(s *melody.Session) error {
	constant.ResponseWithKickOut(s, constant.EVENT_KICKOUT, constant.MsgFlags[constant.EVENT_KICKOUT])
	err := s.Close()
	if err != nil {
		return err
	}
	return nil
}
