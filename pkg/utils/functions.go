package utils

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"location_store/pkg/domain/dto"
)

// GetEnvValue ...
func GetEnvValue( /*key string*/ ) string {
	// return os.Getenv(key)
	return "LOCAL_DEV"
}

func ReadCsvFromPayload(payload io.ReadCloser) ([][]string, error) {
	file := csv.NewReader(payload)
	records, err := file.ReadAll()
	return records, err
}

func WriteHttpJsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func PostgreDuplicateKeyErr(err error) (isDublicateKey bool) {
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return true
		}
	}
	return isDublicateKey
}

func PostgreConnBusyErr(err error) (isDublicateKey bool) {
	if err != nil {
		if strings.Contains(err.Error(), "conn busy") {
			return true
		}
	}
	return isDublicateKey
}

func YelpTooManyRequestsErr(errResponse dto.YelpErrorResponse) (isTooManyRequests bool) {
	if errResponse.Error.Code == "TOO_MANY_REQUESTS_PER_SECOND" {
		return true
	}
	return isTooManyRequests
}
