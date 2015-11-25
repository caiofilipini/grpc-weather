package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	errUnexpectedRespnse = "Unexpected response: %s"
)

type HTTPClient struct{}

var (
	httpClient = HTTPClient{}
)

func (c HTTPClient) get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	c.info(fmt.Sprintf("GET %s -> %d", url, resp.StatusCode))

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf(errUnexpectedRespnse, resp.Status)
		c.info(fmt.Sprintf("Request failed: %v", respErr))
		return nil, respErr
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c HTTPClient) info(msg string) {
	log.Printf("[JSONClient] %s\n", msg)
}
