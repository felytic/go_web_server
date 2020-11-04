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

func getProviderResponse(name string, planId int) *ProviderResponse {
	url := paymentURLs[name]
	addr := fmt.Sprintf("%s?planId=%v", url, planId)
	response, err := http.Get(addr)

	if err != nil {
		log.Print("[ERROR] GET ", addr, " returned ", err)
		return &ProviderResponse{name, "", err}
	}

	if response.StatusCode != 200 {
		log.Print("[ERROR] GET ", addr, " returned ", response.Status)
		return &ProviderResponse{name, "", errors.New(response.Status)}
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Print("[ERROR] GET ", addr, " : ", err)
		return &ProviderResponse{name, "", err}
	}

	return &ProviderResponse{name, string(responseData), err}
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
