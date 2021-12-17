package marketplace

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/RedHatInsights/sources-api-go/config"
)

// GetHttpClient variable that holds the function which returns an HttpClient. This allows us to set up in runtime
// which http client we want for the "GetToken" function, and allows us to mock it easily.
var GetHttpClient func() HttpClient

// HttpClient abstracts away the client to be used in the GetToken function, and allows mocking it easily for the
// tests.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetHttpClientStdlib returns a "http.Client" with a timeout of 10 seconds.
func GetHttpClientStdlib() HttpClient {
	return &http.Client{Timeout: 10}
}

// GetToken sends a request to the marketplace to request a bearer token.
func GetToken(apiKey string) (*BearerToken, error) {
	// Reference docs for the request: https://marketplace.redhat.com/en-us/documentation/api-authentication
	data := url.Values{}
	data.Set("apikey", apiKey)
	data.Set("grant_type", "urn:ibm:params:oauth:grant-type:apikey")

	request, err := http.NewRequest(
		"POST",
		config.Get().MarketplaceHost+"/api-security/om-auth/cloud/token",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, fmt.Errorf("could not create the request object: %s", err)
	}

	// Set the proper headers to accept JSON, and let the server know we're sending urlencoded data.
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := GetHttpClient()
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not perform the request to the marketplace: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code received from the marketplace: %d", response.StatusCode)
	}

	return DecodeMarketplaceTokenFromResponse(response)
}
