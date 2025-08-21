package handler

import "time"

type RequestSuccess struct {
	Msg string `json:"message"`
}

type ErrorResponse struct {
	Errors struct {
		Body string `json:"body"`
	} `json:"errors"`
}

func newErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Errors: struct {
			Body string `json:"body"`
		}{
			Body: msg,
		},
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

type userResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

type userLoginResponse struct {
	Token     string       `json:"access_token"` // DO NOT modify json name, 'access_token' is required in Swagger UI
	TokenType string       `json:"token_type"`   // DO NOT modify json name, 'token_type' is required in Swagger UI
	ExpiredAt int64        `json:"exp"`
	User      userResponse `json:"user"`
}

type LectureResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	Deadline  time.Time `json:"deadline"`
}

type ProblemResponse struct {
	LectureID int64  `json:"lecture_id"`
	ProblemID int64  `json:"problem_id"`
	Title     string `json:"title"`
}
