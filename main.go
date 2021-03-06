package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CallDetails struct {
	SessionID string `json:"sessionId"`
	Token     string `json:"token"`
	Sip       struct {
		URI     string `json:"uri"`
		From    string `json:"from"`
		Headers struct {
		} `json:"headers"`
		Auth struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"auth"`
		Secure bool `json:"secure"`
	} `json:"sip"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
	Body       struct {
		ID           string `json:"id"`
		ConnectionID string `json:"connectionId"`
		StreamID     string `json:"streamId"`
	} `json:"body"`
}

type Body struct {
	ID           string `json:"id"`
	ConnectionID string `json:"connectionId"`
	StreamID     string `json:"streamId"`
}

func HandleRequest(ctx context.Context, call CallDetails) (events.APIGatewayProxyResponse, error) {
	callJSON, err := json.Marshal(call)
	token, _ := creatToken()
	log.Println("Recived request to make a call to: " + call.Sip.URI)
	if err != nil {
		log.Fatal(err)
	}
	url := "https://api.opentok.com/v2/project/" + os.Getenv("API_KEY") + "/dial"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(callJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OPENTOK-AUTH", token)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func creatToken() (string, error) {
	id := uuid.New()
	claims := jwt.MapClaims{}
	claims["iss"] = os.Getenv("API_KEY")
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["jti"] = id.String()
	claims["ist"] = "project"
	claims["iat"] = time.Now()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	log.Println("Token used: " + token)
	return token, nil
}

func main() {
	lambda.Start(HandleRequest)
}
