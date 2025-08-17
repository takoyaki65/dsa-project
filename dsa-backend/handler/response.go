package handler

type createUserSuccess struct {
	Msg string `json:"message"`
}

func newCreateUserSuccess(msg string) createUserSuccess {
	return createUserSuccess{
		Msg: msg,
	}
}

type userLoginResponse struct {
	Token string `json:"token"`
}
