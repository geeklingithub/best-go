package login

import (
	besthttp "github.com/geeklingithub/best-go/server/http"
)

func RegisterRouter(router *besthttp.Router) error {

	router.AddRouter("/login", Login, LoginReq{})
	router.AddRouter("/createPlayer", CreatePlayer, CreatePlayerReq{})
	router.AddRouter("/selectGateway", SelectGateway, SelectGatewayReq{})
	return nil
}

// Login 业务逻辑
func Login(req LoginReq, ctx besthttp.NewContext) {
	ctx.SendResp("Login")
}

// CreatePlayer 业务逻辑
func CreatePlayer(req CreatePlayerReq, ctx besthttp.NewContext) {
	ctx.SendResp("CreatePlayer")
}

// SelectGateway 业务逻辑
func SelectGateway(req SelectGatewayReq, ctx besthttp.NewContext) {
	ctx.SendResp("CreatePlayer")
}
