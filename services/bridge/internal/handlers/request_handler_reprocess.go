package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/services/internal/bridge-compliance-shared/http/helpers"
	"github.com/stellar/go/services/internal/bridge-compliance-shared/protocols/bridge"
)

// Authorize implements /reprocess endpoint
func (rh *RequestHandler) Reprocess(w http.ResponseWriter, r *http.Request) {
	request := &bridge.ReprocessRequest{}
	err := helpers.FromRequest(r, request)
	if err != nil {
		log.Error(err.Error())
		helpers.Write(w, helpers.InvalidParameterError)
		return
	}

	err = request.Validate()
	if err != nil {
		errorResponse := err.(*helpers.ErrorResponse)
		// TODO
		// log.WithFields(errorResponse.LogData).Error(errorResponse.Error())
		helpers.Write(w, errorResponse)
		return
	}

	operation, err := rh.Horizon.LoadOperation(request.OperationID)
	if err != nil {
		helpers.Write(w, &bridge.ReprocessResponse{Status: "error", Message: err.Error()})
		return
	}

	err = rh.PaymentListener.ReprocessPayment(operation, request.Force)

	if err != nil {
		helpers.Write(w, &bridge.ReprocessResponse{Status: "error", Message: err.Error()})
		return
	}

	helpers.Write(w, &bridge.ReprocessResponse{Status: "ok"})
}
