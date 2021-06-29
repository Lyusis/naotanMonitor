package fetcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Lyusis/NaotanMonitor/logger"
)

func GetFetcher(url string) ([]byte, error) {

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	response, responseError := client.Get(url)
	if responseError != nil {
		logger.Sugar.Warn(logger.FormatMsg("A no responded request or a failed response"), logger.FormatTitle("WRONG"), responseError)
		return nil, responseError
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Sugar.Warn(logger.FormatMsg("Failed to close request"), logger.FormatTitle("WRONG"), err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received HTTP request exception, Code: %d", response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}