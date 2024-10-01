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

func handler(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (response bool, err error) {
	log.Println("Evento recebido: ", request)

	err = parseRequest(request.Headers["authorization"])
	if err != nil {
		return false, nil
	}
	return true, nil
}

func main() {
	lambda.Start(handler)
}
