package v1

import (
	"context"

	p "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	"BackendServer/config"
	"BackendServer/constant"
	"BackendServer/proto/backend/fancygame"
	"BackendServer/proto/lobby"
)

func GetMemberList(c *Client, msg []byte) {
	conn, err := grpc.Dial(config.LobbyServerUrl, grpc.WithInsecure())
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}
	defer conn.Close()

	client := fancygame.NewLobbyClient(conn)

	lobbyResp, err := client.GetMemberList(context.Background(), &fancygame.EmptyReq{})
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	byteMsg, _ := p.Marshal(lobbyResp)

	constant.ResponseWithData(c.Session, constant.RESPONSE_GET_MEMBER_LIST_RES, &byteMsg)
}

func SetMember(c *Client, msg []byte) {
	var req fancygame.SetMember
	if err := p.Unmarshal(msg, &req); err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_INVALID_OPERATION, constant.MsgFlags[constant.ERROR_INVALID_OPERATION])
		return
	}

	conn, err := grpc.Dial(config.LobbyServerUrl, grpc.WithInsecure())
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}
	defer conn.Close()

	client := fancygame.NewLobbyClient(conn)

	lobbyResp, err := client.SetMemberData(context.Background(), &req)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	var memberData lobby.UpdateMemberData

	d, _ := p.Marshal(lobbyResp.Data)
	_ = p.Unmarshal(d, &memberData)

	err = UpdateMemberData(c, &memberData)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	byteMsg, _ := p.Marshal(lobbyResp)

	constant.ResponseWithData(c.Session, constant.RESPONSE_GET_MEMBER_LIST_RES, &byteMsg)
}

func KickOutMember(c *Client, msg []byte) {
	var req fancygame.KickOutMemberInfo
	if err := p.Unmarshal(msg, &req); err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_INVALID_OPERATION, constant.MsgFlags[constant.ERROR_INVALID_OPERATION])
		return
	}

	memberClient := ConnPool.users[req.Account]

	conn, err := grpc.Dial(config.LobbyServerUrl, grpc.WithInsecure())
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}
	client := fancygame.NewLobbyClient(conn)

	lobbyReq := fancygame.KickOutMemberInfo{
		Account: req.Account,
		Balance: memberClient.Balance,
	}

	userResp, err := client.KickOutMember(context.Background(), &lobbyReq)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	err = ConnPool.DeleteClient(memberClient, memberClient.IsAdmin)
	if err != nil {
		constant.ResponseWithError(c.Session, constant.ERROR_OPERATION_FAILED, err.Error())
		return
	}

	var resp = lobby.KickOutMemberRes{
		StatusCode: userResp.StatusCode,
		Message:    userResp.Message,
	}

	byteMsg, _ := p.Marshal(&resp)

	constant.ResponseWithData(c.Session, constant.RESPONSE_GET_MEMBER_LIST_RES, &byteMsg)
}

func UpdateMemberData(c *Client, d *lobby.UpdateMemberData) error {
	byteData, err := p.Marshal(d)
	if err != nil {
		return err
	}
	byteData = append([]byte{byte(0), byte(constant.EVENT_UPDATE_MEMBER_DATA)}, byteData...)

	for _, poolClient := range ConnPool.admins {
		if poolClient.Account != c.Account {
			err = poolClient.Session.WriteBinary(byteData)
			if err != nil {
				return err
			}
		}
	}

	err = ConnPool.users[d.Account].Session.WriteBinary(byteData)
	if err != nil {
		return err
	}
	return nil
}

/*
func UpdateGameList(){

}*/
