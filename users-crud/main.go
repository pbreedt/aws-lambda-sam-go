package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestPayload struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type ResponsePayload struct {
	ID   string `json:"id"`
	Job  string `json:"job,omitempty"`
	Name string `json:"name,omitempty"`
	Data struct {
		Email   string `json:"email,omitempty"`
		Name    string `json:"first_name,omitempty"`
		Surname string `json:"last_name,omitempty"`
		Avatar  string `json:"avatar,omitempty"`
	} `json:"data,omitempty"`
}

const (
	apiKey = "fake_api_key"
	apiUrl = "https://reqres.in/api/users"
)

var ErrInvalidRequest = fmt.Errorf("invalid request")
var ErrNon200Response = fmt.Errorf("non 200 response")

func apiGatewayProxyHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Received request: ", request.Body)

	switch request.HTTPMethod {
	case "GET":
		return doGet(&request)
	case "POST":
		return doPost(&request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "Method Not Allowed",
			StatusCode: 405,
		}, nil
	}
}

func main() {
	lambda.Start(apiGatewayProxyHandler)
}

func doGet(request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("url: %s\n", apiUrl)
	// Part 2: get data from backend service
	resp, err := http.Get(apiUrl + "/" + request.QueryStringParameters["id"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Printf("received backend response: %s/%v\n", resp.Status, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Printf("Received response: %s\n", respBody)

	// Part 3: parse response
	var respPayload ResponsePayload
	err = json.Unmarshal(respBody, &respPayload)
	if err != nil {
		return events.APIGatewayProxyResponse{}, ErrInvalidRequest
	}
	fmt.Printf("user.name: %s\n", respPayload.Data.Name)
	fmt.Printf("user.surname: %s\n", respPayload.Data.Surname)

	// Part 4: return APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		Body:       string(respBody),
		StatusCode: 200,
	}, nil
}

func doPost(request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Part 1: get data from request
	var payload RequestPayload
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{}, ErrInvalidRequest
	}

	log.Printf("url: %s\n", apiUrl)
	// Part 2: get data from backend service
	resp, err := http.Post(apiUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"name":"%s","job":"%s"}`, payload.Name, payload.Job)))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Printf("received backend response: %s/%s\n", resp.Status, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Printf("Received response: %s\n", respBody)

	// Part 3: parse response
	var respPayload ResponsePayload
	err = json.Unmarshal(respBody, &respPayload)
	if err != nil {
		return events.APIGatewayProxyResponse{}, ErrInvalidRequest
	}
	fmt.Printf("user.id: %s\n", respPayload.ID)

	// Part 4: return APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		Body:       string(respBody),
		StatusCode: 200,
	}, nil
}

func (r ResponsePayload) String() string {
	return fmt.Sprintf("{\"id\":\"%s\",\"first_name\":\"%s\",\"last_name\":\"%s\",\"email\":\"%s\",\"avatar\":\"%s\"}", r.ID, r.Name, r.Data.Surname, r.Data.Email, r.Data.Avatar)
}
