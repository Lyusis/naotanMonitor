package engine

import (
	"logger"
	"math/rand"
	"monitor/fetcher"
	"time"
)

func Worker(request Request) (ResultItems, error) {

	if request.Name != "" {
		logger.Sugar.Info(logger.FormatMsg("Fetching"), logger.FormatTitle("URL"), request.Url, logger.FormatTitle("Name"), request.Name)
	} else {
		logger.Sugar.Info(logger.FormatMsg("Fetching"), logger.FormatTitle("URL"), request.Url)
	}
	body, bodyErr := fetcher.GetFetcher(request.Url)
	if bodyErr != nil {
		logger.Sugar.Error(logger.FormatMsg("Failed to receive request body"), bodyErr)
		logger.WriteFile(logger.FormatMsg("Writing failure infomation"),
			time.Now().Format(logger.TimeFormatDate)+"_fail-request-body"+string(rune(rand.Intn(19960730)))+".log", body)
		return NilResult(), bodyErr
	}

	result := request.PrimaryParser(body)

	return result, nil
}
