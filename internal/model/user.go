package model

type User struct {
	ID      int
	Login   string
	HashPas []byte
	Level   int
}
