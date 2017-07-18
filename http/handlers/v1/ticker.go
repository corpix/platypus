package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/datasources"
	"github.com/cryptounicorns/market-fetcher-http/logger"
)

func Mount(r *mux.Router, d *datasources.Datasources, l logger.Logger) {
	v := r.PathPrefix("/v1").Subrouter()

	v.HandleFunc(
		"/tickers",
		func(rw http.ResponseWriter, r *http.Request) {
			// FIXME: Very dirty, requires refactoring
			conn, _, _, err := ws.UpgradeHTTP(r, rw, nil)
			if err != nil {
				l.Error(err)
				return
			}

			go func() {
				var (
					res []byte
					err error
				)
				defer conn.Close()

				res, err = json.Marshal(
					d.Ticker.Store.FindPrefix([]string{}),
				)
				if err != nil {
					l.Error(err)
					return
				}
				err = wsutil.WriteServerMessage(
					conn,
					ws.OpBinary,
					res,
				)
				if err != nil {
					l.Error(err)
					return
				}

				for v := range d.Ticker.Feed {
					res, err = json.Marshal(v)
					if err != nil {
						l.Error(err)
						return
					}

					// FIXME: Should we read from client to mitigate
					// high memory consumption attacks?
					err = wsutil.WriteServerMessage(
						conn,
						ws.OpBinary,
						res,
					)
					if err != nil {
						l.Error(err)
						return
					}
				}
			}()
		},
	)
}
