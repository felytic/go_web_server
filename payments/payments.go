package payments

import (
	"fmt"
	"log"
	"net/http"
)

var paymentURLs = map[string]string{
	"applePayUrl":  "http://0.0.0.0:8090",
	"googlePayUrl": "http://0.0.0.0:8091",
	"payPalUrl":    "http://0.0.0.0:8092",
	"stripeUrl":    "http://0.0.0.0:8093",
}

type PaymentResponse struct {
	Name string
	Res  http.Response
	Err  error
}

func getProviderResponse(name string, planId int) *PaymentResponse {
	url := paymentURLs[name]
	resp, err := http.Get(fmt.Sprintf("%s?planId=%v", url, planId))

	if err != nil {
		log.Print("[ERROR] ", err)
	}

	return &PaymentResponse{name, *resp, err}
}

func GetPaymentResponses(planId int) []PaymentResponse {
	resultsChan := make(chan *PaymentResponse)

	// make sure we close these channels when we're done with them
	defer func() {
		close(resultsChan)
	}()

	for name := range paymentURLs {
		go func(name string) {
			resp := getProviderResponse(name, planId)
			resultsChan <- resp
		}(name)
	}

	var results []PaymentResponse

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
