package domain

type User struct {
	Id             int64
	Login          string
	Password       string
	HashedPassword string
}
