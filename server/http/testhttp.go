package http

import "fmt"

func BusinessRouter() map[string]func(any) any {
	return map[string]func(any) any{
		"/login": func(req any) any {
			fmt.Println(req)
			return "/login"
		},
		"/createPlayer": func(req any) any {
			fmt.Println(req)
			return "/createPlayer"
		},
		"/selectGateway": func(req any) any {
			fmt.Println(req)
			return "/selectGateway"
		},
	}
}

func AdminRouter() map[string]func(any) any {
	return map[string]func(any) any{
		"/admin": func(req any) any {
			fmt.Println(req)
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
