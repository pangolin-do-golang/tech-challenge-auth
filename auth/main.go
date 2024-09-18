package main

import (
	"encoding/base64"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Statement struct {
	Action   string `json:"Action"`
	Effect   string `json:"Effect"`
	Resource string `json:"Resource"`
}

type Response struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type EventHeaders struct {
	XAMZDate                 string `json:"X-AMZ-Date"`
	Accept                   string `json:"Accept"`
	Authorization            string `json:"Authorization"`
	CloudFrontViewerCountry  string `json:"CloudFront-Viewer-Country"`
	CloudFrontForwardedProto string `json:"CloudFront-Forwarded-Proto"`
	CloudFrontIsTabletViewer string `json:"CloudFront-Is-Tablet-Viewer"`
	CloudFrontIsMobileViewer string `json:"CloudFront-Is-Mobile-Viewer"`
	UserAgent                string `json:"User-Agent"`
}

type Event struct {
	MethodArn  string       `json:"methodArn"`
	Resource   string       `json:"resource"`
	Path       string       `json:"path"`
	HttpMethod string       `json:"httpMethod"`
	Headers    EventHeaders `json:"headers"`
}

func getUser(id string) error {
	baseUrl := os.Getenv("API_URL")
	log.Println("Fazendo req para " + baseUrl)
	req, err := http.Get(baseUrl + "/customer/" + id)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	log.Println("Resposta recebida com status " + strconv.Itoa(req.StatusCode))

	if req.StatusCode != 200 {
		return errors.New("customer not found")
	}

	return nil
}

func parseRequest(header string) error {
	if len(header) == 0 {
		return errors.New("authorization header not found")
	}

	decodedIdentifier, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return err
	}

	log.Println("Processando evento para identificador " + string(decodedIdentifier))

	return getUser(string(decodedIdentifier))
}

func handler(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (response *events.APIGatewayCustomAuthorizerResponse, err error) {
	log.Println("Evento recebido: ", request)

	effect := "Allow"
	err = parseRequest(request.Headers["Authorization"])
	if err != nil {
		effect = "Deny"
	}

	log.Println(err)

	return &events.APIGatewayCustomAuthorizerResponse{
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{request.MethodArn},
				},
			},
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
