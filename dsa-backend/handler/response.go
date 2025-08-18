package handler

type createUserSuccess struct {
	Msg string `json:"message"`
}

func newCreateUserSuccess(msg string) createUserSuccess {
	return createUserSuccess{
		Msg: msg,
	}
}

type userResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

type userLoginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}
