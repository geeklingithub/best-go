package login

type LoginReq struct {
	OpenId   string
	PlayerId uint64
}

type CreatePlayerReq struct {
	OpenId   string
	PlayerId uint64
}

type SelectGatewayReq struct {
	OpenId   string
	PlayerId uint64
}
