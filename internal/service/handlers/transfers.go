package handlers

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func IndexTransfers(w http.ResponseWriter, r *http.Request) {
	var response *[]data.Transfer

	response, err := Db(r).Transfer().Select()
	if err != nil {
		Log(r).Error("Error in IndexTransfers: ", err)
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Title:  "Error in IndexTransfers",
			Detail: err.Error(),
		})
	}

	ape.Render(w, response)
}
