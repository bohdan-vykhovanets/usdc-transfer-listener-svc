package requests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/urlval"

	"net/http"
)

const (
	DefaultLimit  uint64 = 20
	MaxLimit      uint64 = 100
	DefaultOffset uint64 = 0
)

type TransfersRequest struct {
	ToFilter           *string `filter:"to"`
	FromFilter         *string `filter:"from"`
	CounterpartyFilter *string `filter:"counterparty"`
	PageLimit          uint64  `page:"limit"`
	PageOffset         uint64  `page:"offset"`
}

func NewTransfersRequest(r *http.Request) (TransfersRequest, error) {
	var request TransfersRequest
	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return request, errors.Wrap(err, "failed to decode query")
	}

	if request.PageLimit <= 0 || request.PageLimit > MaxLimit {
		request.PageLimit = DefaultLimit
	}
	if request.PageOffset <= 0 {
		request.PageOffset = DefaultOffset
	}

	if request.ToFilter != nil && !common.IsHexAddress(*request.ToFilter) {
		request.ToFilter = nil
	}
	if request.FromFilter != nil && !common.IsHexAddress(*request.FromFilter) {
		request.FromFilter = nil
	}
	if request.CounterpartyFilter != nil && !common.IsHexAddress(*request.CounterpartyFilter) {
		request.CounterpartyFilter = nil
	}

	return request, nil
}
