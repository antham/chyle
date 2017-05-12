package chyle

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ErrResponse is triggered when status code is greater or equal to 400
type ErrResponse struct {
	request  *http.Request
	response *http.Response
	body     []byte
}

// Error output error as string
func (e ErrResponse) Error() string {
	return fmt.Sprintf("an error occurred when contacting remote api through %s, status code %d, body %s", e.request.URL, e.response.StatusCode, e.body)
}

// setHeaders setup headers on request from a map header key -> header value
func setHeaders(request *http.Request, headers map[string]string) {
	for k, v := range headers {
		request.Header.Set(k, v)
	}
}

// sendRequest picks a request and send it with given client
func sendRequest(client *http.Client, request *http.Request) (int, []byte, error) {
	response, err := client.Do(request)

	if err != nil {
		return 0, []byte{}, err
	}

	defer func() {
		err = response.Body.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return response.StatusCode, b, ErrResponse{request, response, b}
	}

	if response.StatusCode >= 400 {
		return response.StatusCode, b, ErrResponse{request, response, b}
	}

	return response.StatusCode, b, nil
}
