package userrole

type Role int64

const (
	Admin Role = iota
	Manager
	Student
)
