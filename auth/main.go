package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"time"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req, err := http.Get("http://localhost:8080/customer/05936752070")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return nil, errors.New("customer not found")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	type Customer struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	var customer Customer
	if err = json.Unmarshal(body, &customer); err != nil {
		return nil, err
	}

	expiration := time.Now().Add(time.Minute * 10)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": customer.ID,
		"nbf":  expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))

	return &events.APIGatewayProxyResponse{
		Body:       tokenString,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
