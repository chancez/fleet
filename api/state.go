package api

import (
	"errors"
	"net/http"
	"path"

	log "github.com/coreos/fleet/Godeps/_workspace/src/github.com/golang/glog"

	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/schema"
)

func wireUpStateResource(mux *http.ServeMux, prefix string, cAPI client.API) {
	base := path.Join(prefix, "state")
	sr := stateResource{cAPI, base}
	mux.Handle(base, &sr)
}

type stateResource struct {
	cAPI     client.API
	basePath string
}

func (sr *stateResource) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		sendError(rw, http.StatusMethodNotAllowed, errors.New("only GET supported against this resource"))
		return
	}

	sr.list(rw, req)
}

func (sr *stateResource) list(rw http.ResponseWriter, req *http.Request) {
	states, err := sr.cAPI.UnitStates()
	if err != nil {
		log.Errorf("Failed fetching UnitStates: %v", err)
		sendError(rw, http.StatusInternalServerError, nil)
		return
	}

	page := schema.UnitStatePage{States: states}
	sendResponse(rw, http.StatusOK, &page)
}