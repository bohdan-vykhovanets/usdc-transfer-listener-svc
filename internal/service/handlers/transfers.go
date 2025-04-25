package handlers

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/service/requests"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strings"
)

func baseError(detail string, status string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Detail: detail,
		Status: status,
	}
}

type IndexResponse struct {
	Transfers  []data.Transfer `json:"transfers"`
	TotalCount int             `json:"total_count"`
}

func IndexTransfers(w http.ResponseWriter, r *http.Request) {
	logger := Log(r)
	var response *[]data.Transfer

	request, err := requests.NewTransfersRequest(r)
	if err != nil {
		logger.Errorf("request is incorrect")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := Db(r).Transfer()

	if request.ToFilter != nil {
		q = q.FilterByTo(strings.ToLower(*request.ToFilter))
	}
	if request.FromFilter != nil {
		q = q.FilterByFrom(strings.ToLower(*request.FromFilter))
	}
	if request.CounterpartyFilter != nil {
		q = q.FilterByCounterparty(strings.ToLower(*request.CounterpartyFilter))
	}

	response, err = q.Paginate(request.PageLimit, request.PageOffset).Select()
	if err != nil {
		logger.Error("Error in IndexTransfers: ", err)
		ape.RenderErr(w, baseError(err.Error(), "500"))
	}

	ape.Render(w, response)
}
