package transports

import (
	"context"
	"encoding/json"

	ll "log"
	"net/http"
	https "net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/susinda/endpoints"
	models "github.com/susinda/models/request"
	"github.com/susinda/services"
)

type HttpTransport struct {
	router *mux.Router
}

func (self *HttpTransport) Init() {
	logger := log.NewLogfmtLogger(os.Stderr)
	log.With(logger)

	// Initialize endpoints
	responseEndPoint := endpoints.RequestEndpoint{}

	// Initialize cors headers
	headersOk := handlers.AllowedHeaders([]string{"x-requested-with", "access-control-allow-origin", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Initialize Routes
	self.router = mux.NewRouter()

	self.post("/request/input", responseEndPoint.InitRequest(services.RequestService{}))
	// Initialize Server
	logger.Log("INFO", "Server started and listening at http://localhost:8080")

	
	ll.Fatal(https.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(self.router)))
}


func (self *HttpTransport) get(path string, endpoint endpoint.Endpoint, decoder httptransport.DecodeRequestFunc) {
	self.router.Methods("GET").Path(path).Handler(httptransport.NewServer(endpoint, decoder, encodeResponse))
}
func (self *HttpTransport) post(path string, endpoint endpoint.Endpoint) {
	self.router.Methods("POST").Path(path).Handler(httptransport.NewServer(endpoint, decodeSaveRequest, encodeResponse))
}

func (self *HttpTransport) put(path string, endpoint endpoint.Endpoint) {
	self.router.Methods("PUT").Path(path).Handler(httptransport.NewServer(endpoint, decodeSaveRequest, encodeResponse))
}

func (self *HttpTransport) patch(path string, endpoint endpoint.Endpoint) {
	self.router.Methods("PATCH").Path(path).Handler(httptransport.NewServer(endpoint, decodeSaveRequest, encodeResponse))
}

func (self *HttpTransport) delete(path string, endpoint endpoint.Endpoint) {
	self.router.Methods("DELETE").Path(path).Handler(httptransport.NewServer(endpoint, decodeDeleteOneRequest, encodeResponse))
}



func decodeFindOneRequest(_ context.Context, req *http.Request) (interface{}, error) {
	var findOneReq models.FindOneRequest = models.FindOneRequest{}
	findOneReq.ID, _ = getId(req.URL.Path)
	return findOneReq, nil
}


func decodeExclusionRequest(_ context.Context, req *http.Request) (interface{}, error) {
	searchReq := models.ExclusionRequest{}
	var err error
	searchReq.Type = req.URL.Query().Get("type")
	return searchReq, err
}


func decodeSaveRequest(_ context.Context, req *http.Request) (interface{}, error) {
	return req.Body, nil
}


func decodeDeleteOneRequest(_ context.Context, req *http.Request) (interface{}, error) {
	var deleteOneReq models.DeleteRequest = models.DeleteRequest{}
	deleteOneReq.ID, _ = getId(req.URL.Path)
	return deleteOneReq, nil
}



func encodeResponse(_ context.Context, w https.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}



func getId(path string) (int, error) {
	p := strings.Split(path, "/")
	return strconv.Atoi(p[len(p)-1])
}

func getPathId(path string) (string, error) {
	p := strings.Split(path, "/")
	return (p[len(p)-1]), nil
}

