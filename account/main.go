package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/myndshft/nemgo"
)

type Account struct {
	Address string  `json:"address"`
	Balance float32 `json:"balance"`
}

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	address := request.QueryStringParameters["address"]
	env := request.QueryStringParameters["env"]

	// envパラメータによって環境を変える
	host := "23.228.67.85:7890"
	network := nemgo.Testnet
	if env == "mainnet" {
		host = "209.126.98.204:7890"
		network = nemgo.Mainnet
	}

	// AccountData
	client := nemgo.New(nemgo.WithNIS(host, network))
	// Obtain account data using address
	accountData, err := client.AccountData(nemgo.Address(address))
	if err != nil {
		log.Fatal(err)
	}

	balance := float32(accountData.Account.Balance) / float32(math.Pow(10, 6.0))
	account := Account{Address: accountData.Account.Address, Balance: balance}
	bs, _ := json.Marshal(account)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", string(bs)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
