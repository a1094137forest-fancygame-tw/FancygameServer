package v1

import (
	"log"

	"gopkg.in/olahol/melody.v1"

	"BackendServer/constant"
)

var (
	ConnPool *ConnectionPool

	funcMap = map[int]func(*Client, []byte){
		// common
		constant.REQUEST_PING:     Pong,
		constant.REQUEST_REGISTER: Register,
		constant.REQUEST_LOGIN:    Login,
		constant.REQUEST_LOGOUT:   Logout,

		// lobby
		constant.REQUEST_MEMBER_LIST: GetMemberList,
		constant.REQUEST_SET_MEMBER:  SetMember,
	}
)

type Client struct {
	Account string
	Avatar  int32
	Gender  int32
	Balance int64
	IsWait  bool
	IsAdmin bool
	Session *melody.Session
}

type ConnectionPool struct {
	admins map[string]*Client
	users  map[string]*Client
}

func Initialize() {
	ConnPool = SetPool()
}

func GetMessage(c *Client, messageCode int, msg []byte) {
	log.Println("messageCode:", messageCode)

	if fn, ok := funcMap[messageCode]; ok {
		log.Println("get api")
		fn(c, msg)
	} else {
		log.Println("unknown api")
		Error(c, []byte{})
	}
}

func SetPool() *ConnectionPool {
	return &ConnectionPool{
		admins: make(map[string]*Client),
		users:  make(map[string]*Client),
	}
}

func (cp *ConnectionPool) UpsertClient(client *Client, isAdmin bool) error {
	if isAdmin {
		if _, ok := cp.admins[client.Account]; ok {
			err := KickOut(cp.admins[client.Account].Session)
			if err != nil {
				return err
			}
		}
		cp.admins[client.Account] = client
		log.Println("[info] New connect set in admin pool:", client.Account)
		return nil
	} else {
		if _, ok := cp.users[client.Account]; ok {
			err := KickOut(cp.users[client.Account].Session)
			if err != nil {
				return err
			}
		}
		cp.users[client.Account] = client
		log.Println("[info] New connect set in user pool:", client.Account)

		return nil
	}
}

func (cp *ConnectionPool) DeleteClient(client *Client, isAdmin bool) error {
	if isAdmin {
		err := KickOut(cp.admins[client.Account].Session)
		if err != nil {
			return err
		}
		delete(cp.admins, client.Account)
		return nil
	}
	err := KickOut(cp.users[client.Account].Session)
	if err != nil {
		return err
	}
	delete(cp.users, client.Account)
	return nil
}
