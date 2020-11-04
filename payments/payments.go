package payments

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var paymentURLs = map[string]string{
	"applePayUrl":  "http://0.0.0.0:8090",
	"googlePayUrl": "http://0.0.0.0:8091",
	"payPalUrl":    "http://0.0.0.0:8092",
	"stripeUrl":    "http://0.0.0.0:8093",
}

type ProviderResponse struct {
	Name string
	Url  string
	Err  error
}

func IsAllResponsesOk(responses []ProviderResponse) bool {
	for _, response := range responses {
		if response.Err != nil || response.Url == "" {
			return false
		}
	}

	return true
}

func ConstructProviderResponse(name string, response *http.Response, err error) *ProviderResponse {
	if err != nil {
		log.Print("[ERROR] GET ", name, " URL returned ", err)
		return &ProviderResponse{name, "", err}
	}

	if response.StatusCode != 200 {
		log.Print("[ERROR] GET ", response.Request.URL, " returned ", response.Status)
		return &ProviderResponse{name, "", errors.New(response.Status)}
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Print("[ERROR] GET ", response.Request.URL, " : ", err)
		return &ProviderResponse{name, "", err}
	}

	return &ProviderResponse{name, string(responseData), err}
}

func getProviderResponse(name string, planId int) *ProviderResponse {
	url := paymentURLs[name]
	response, err := http.Get(fmt.Sprintf("%s?planId=%v", url, planId))

	return ConstructProviderResponse(name, response, err)
}

func GetProvidersURLs(planId int) []ProviderResponse {
	// Chanel for storing all async responses
	resultsChan := make(chan *ProviderResponse)

	// make sure we close these channels when we're done with them
	defer func() {
		close(resultsChan)
	}()

	for name := range paymentURLs {
		go func(name string) {
			resultsChan <- getProviderResponse(name, planId)
		}(name)
	}

	var results []ProviderResponse

	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(paymentURLs) {
			break
		}
	}

	return results
}
