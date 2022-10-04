package endpoints

import (
	"context"
	"encoding/json"
	"io"

	"github.com/go-kit/kit/endpoint"
	models "github.com/susinda/models"
	"github.com/susinda/services"
)


type RequestEndpoint struct{}

func (RequestEndpoint) InitRequest(service services.RequestService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		body := request.(io.ReadCloser)
		dec := json.NewDecoder(body)
		dec.DisallowUnknownFields()
		var requestInit models.InitRequestModel
		err := dec.Decode(&requestInit)
		if err != nil {
			return nil, err
		}
		return service.ProducerActiveMq(requestInit.Name)
	}
}