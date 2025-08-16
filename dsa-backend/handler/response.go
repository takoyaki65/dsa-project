package handler

type initCheckResponse struct {
	Initialized bool `json:"initialized"`
}

func newInitCheckResponse(initialized bool) initCheckResponse {
	return initCheckResponse{
		Initialized: initialized,
	}
}

type createUserSuccess struct {
	Msg string `json:"message"`
}

func newCreateUserSuccess(msg string) createUserSuccess {
	return createUserSuccess{
		Msg: msg,
	}
}
