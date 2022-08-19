package http

func BusinessRouter() map[string]func() any {
	return map[string]func() any{
		"/login": func() any {
			return "/login"
		},
		"/createPlayer": func() any {
			return "/createPlayer"
		},
		"/selectGateway": func() any {
			return "/selectGateway"
		},
	}
}

func AdminRouter() map[string]func() any {
	return map[string]func() any{
		"/admin": func() any {
			return LoginReq{
				PlayerId: 1,
				OpenId:   "Reds",
			}
		},
	}
}

type LoginReq struct {
	OpenId   string
	PlayerId uint64
}
