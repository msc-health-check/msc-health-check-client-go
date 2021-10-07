package v1

import (
	glog "bitbucket.org/maironmscosta/golang-log/v1"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	ErrorNewRequest          = "erro ao iniciar o request"
	ErrorExecuteRequest      = "erro ao executar o request"
	ErrorFromAPI             = "erro vindo da API"
	ErrorReadingResponseBody = "erro ao ler o corpo do response"
	NotFound                 = "sem resultado"

	ServicePath = "/health-check"
)

var logger *log.Logger

func init() {

	logger = log.New(os.Stdout, "", log.Flags())

}

type Client interface {
	AddService(request ProjectCheckRequest) (ProjectCheck, error)
	AddLiveSignal(ID, appName string) (ProjectCheck, error)
}

func NewClient(httpClient *http.Client, host string) *ClientImpl {

	var logging = glog.NewLogging(logger)
	return &ClientImpl{
		Host:       host,
		HttpClient: httpClient,
		Logger:     logging,
	}
}

type ClientImpl struct {
	Host       string
	HttpClient *http.Client
	Logger     glog.Logging
}

func (client *ClientImpl) executeRequest(methodRequest string, url string, jsonBody []byte) ([]byte, error) {

	var err error
	var request *http.Request
	var payload *strings.Reader
	if jsonBody != nil {
		payload = strings.NewReader(string(jsonBody))
		request, err = http.NewRequest(methodRequest, url, payload)
	} else {
		request, err = http.NewRequest(methodRequest, url, nil)
	}

	if err != nil {
		return nil, errors.New(ErrorNewRequest)
	}

	response, err := client.HttpClient.Do(request)
	if err != nil {
		return nil, errors.New(ErrorExecuteRequest)
	}

	defer response.Body.Close()

	if !is2XX(response) {
		if response.StatusCode == 404 {
			return nil, errors.New(NotFound)
		}
		return nil, errors.Wrap(err, ErrorFromAPI)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(ErrorReadingResponseBody)
	}

	return responseBody, nil
}

func is2XX(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}

func (client *ClientImpl) AddService(service ProjectCheckRequest) (ProjectCheck, error) {
	// http://localhost:8081/ara-api/comment/post/20200812BpLnf/
	url := fmt.Sprintf("%s/%s", client.Host, ServicePath)
	body, err := json.Marshal(service)
	if err != nil {
		return ProjectCheck{}, err
	}

	responseBody, err := client.executeRequest(http.MethodPost, url, body)
	if err != nil {
		return ProjectCheck{}, err
	}

	var project ProjectCheck
	err = json.Unmarshal(responseBody, &project)
	if err != nil {
		return ProjectCheck{}, err
	}

	return project, nil
}

func (client *ClientImpl) AddLiveSignal(ID, appName string) (ProjectCheck, error) {
	// http://localhost:8081/ara-api/comment/post/20200812BpLnf/
	url := fmt.Sprintf("/%s/%s/%s/%s", client.Host, ServicePath, appName, ID)

	responseBody, err := client.executeRequest(http.MethodPut, url, nil)
	if err != nil {
		return ProjectCheck{}, err
	}

	var project ProjectCheck
	err = json.Unmarshal(responseBody, &project)
	if err != nil {
		return ProjectCheck{}, err
	}

	return project, nil
}
