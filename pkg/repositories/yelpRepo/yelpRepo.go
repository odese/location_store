package yelpRepo

import (
	"encoding/json"
	"io"
	"location_store/pkg/domain/dto"
	"location_store/pkg/infrastructure/config"
	log "location_store/pkg/infrastructure/logger"
	"location_store/pkg/utils"
	"net/http"
)

// TODO: In case of requests are too fast like exceeding rate limiting to the Yelp API, or any other error,
// we should implement some retry mechanism until it get success.
// FIXED: For rate limiting error recursion is applied.
// SearchBusinesses returns list of businesses from Yelp
func SearchBusinesses(term, latitude, longitude, offset string) (response dto.YelpSearchBusinessesResponse, err error) {
	url := prepareURLForSearchingBusinesses(term, latitude, longitude, offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Str("URL", url).Msg("Error on creating request to yelp")
		return response, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", config.Call().GetString("yelp.token"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("URL", url).Msg("Error on sending request to yelp")
		return response, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Str("URL", url).Msg("Error on reading response from yelp")
		return response, err
	}
	defer res.Body.Close()

	response, err = parseSearchBusinessesResponse(body)
	if err != nil {
		log.Error().Err(err).Msg("Error on parsing response from yelp")
		return response, err
	}

	if len(response.Businesses) == 0 {

		errResponse, err := parseYelpResponseForError(body)
		if err != nil {
			log.Error().Err(err).Msg("Error on parsing response from yelp")
			return response, err
		}

		if utils.YelpTooManyRequestsErr(errResponse) {
			log.Debug().Strs("Request Params", []string{term, latitude, longitude, offset}).Msg("Too many requests to Yelp API, retrying...")
			response, err = SearchBusinesses(term, latitude, longitude, offset)
			return response, err
		}

	}

	return response, err
}

// XXX: Yelp says that Limit + Offset must be less than or equal to 1000.
func prepareURLForSearchingBusinesses(term, latitude, longitude, offset string) string {
	url := utils.YelpSearchBusinessesURL + "?limit=" + utils.LimitForYelpStr + "&radius=" + utils.NearbyRadiusForYelp
	if term != "" {
		url += "&term=" + term
	}
	if latitude != "" {
		url += "&latitude=" + latitude
	}
	if longitude != "" {
		url += "&longitude=" + longitude
	}
	if offset != "" {
		url += "&offset=" + offset
	}
	return url
}

func parseSearchBusinessesResponse(body []byte) (response dto.YelpSearchBusinessesResponse, err error) {
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error().Err(err).Msg("Error on parsing response")
		return response, err
	}
	return response, err
}

func parseYelpResponseForError(body []byte) (errorResponse dto.YelpErrorResponse, err error) {
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		log.Error().Err(err).Msg("Error on parsing response")
		return
	}
	return errorResponse, err
}
