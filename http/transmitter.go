package http

import (
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/platypus/http/handlers/routers/router"
)

func NewRouterWriter() router.Writer {
	return wsutil.WriteServerBinary
}
