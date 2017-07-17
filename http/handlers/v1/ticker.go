package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/datasources"
	"github.com/cryptounicorns/market-fetcher-http/logger"
)

func Mount(r *mux.Router, d *datasources.Datasources, l logger.Logger) {
	v := r.PathPrefix("/v1").Subrouter()

	v.HandleFunc(
		"/tickers",
		func(rw http.ResponseWriter, r *http.Request) {
			res, err := json.Marshal(
				d.Ticker.Store.FindPrefix([]string{}),
			)
			if err != nil {
				l.Error(err)
			}
			rw.Write(res)
		},
	)
}
