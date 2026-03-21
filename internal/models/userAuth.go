package models

type UserAuth struct {
	Id       string
	Password string
	Roles    []string
}
