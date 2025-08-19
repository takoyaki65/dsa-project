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
	Token     string       `json:"access_token"` // DO NOT modify json name, 'access_token' is required in Swagger UI
	TokenType string       `json:"token_type"`   // DO NOT modify json name, 'token_type' is required in Swagger UI
	User      userResponse `json:"user"`
}
