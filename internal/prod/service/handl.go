package service

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Register(router *httprouter.Router)
}
