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
	header := func() string {
		if request.Headers["Authorization"] != "" {
			return request.Headers["Authorization"]
		}
		return request.Headers["authorization"]
	}()
	err = parseRequest(header)
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
