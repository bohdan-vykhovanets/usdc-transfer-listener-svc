package handlers

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultLimit  uint64 = 20
	MaxLimit      uint64 = 100
	DefaultOffset uint64 = 0
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
	query := r.URL.Query()

	q := Db(r).Transfer()

	if rawFrom := query.Get("from"); rawFrom != "" {
		if !common.IsHexAddress(rawFrom) {
			logger.Errorf("invalid address: %s", rawFrom)
			ape.RenderErr(w, baseError("Invalid address format in from param", "400"))
		}
		from := strings.ToLower(rawFrom)
		logger.Infof("filter by from: %s", from)
		q = q.FilterByFrom(from)
	}
	if rawTo := query.Get("to"); rawTo != "" {
		if !common.IsHexAddress(rawTo) {
			logger.Errorf("invalid address: %s", rawTo)
			ape.RenderErr(w, baseError("Invalid address format in to param", "400"))
		}
		to := strings.ToLower(rawTo)
		logger.Infof("filter by to: %s", to)
		q = q.FilterByTo(to)
	}
	if rawCounterparty := query.Get("counterparty"); rawCounterparty != "" {
		if !common.IsHexAddress(rawCounterparty) {
			logger.Errorf("invalid address: %s", rawCounterparty)
			ape.RenderErr(w, baseError("Invalid address format in counterparty param", "400"))
		}
		counterparty := strings.ToLower(rawCounterparty)
		logger.Infof("filter by counterparty: %s", counterparty)
		q = q.FilterByCounterparty(counterparty)
	}

	limit := DefaultLimit
	offset := DefaultOffset
	var err error

	if rawLimit := query.Get("limit"); rawLimit != "" {
		limit, err = strconv.ParseUint(rawLimit, 10, 64)
		if err != nil || limit <= 0 || limit > MaxLimit {
			logger.Info("Invalid limit value: %v, kept default value %v", rawLimit, DefaultLimit)
		}
	}
	if rawOffset := query.Get("offset"); rawOffset != "" {
		offset, err = strconv.ParseUint(rawOffset, 10, 64)
		if err != nil || offset < 0 {
			logger.Info("Invalid offset value: %v, kept default value %v", rawOffset, DefaultOffset)
		}
	}

	response, err = q.Paginate(limit, offset).Select()
	if err != nil {
		logger.Error("Error in IndexTransfers: ", err)
		ape.RenderErr(w, baseError(err.Error(), "500"))
	}

	ape.Render(w, response)
}
