package handlers

import (
	"net/http"
)

type Handler interface {
	http.Handler
	Close() error
}
