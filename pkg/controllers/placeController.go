package controllers

import (
	logic "location_store/pkg/businessLogic"
	log "location_store/pkg/infrastructure/logger"
	"location_store/pkg/utils"
	"net/http"
	"strconv"
)

func UploadCSV(w http.ResponseWriter, r *http.Request) {
	payload := r.Body

	places, err := logic.ParseCsvIntoPlaceModel(payload)
	if err != nil {
		log.Error().Err(err).Msg("Error on parsing CSV into Place model")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = logic.CreatePlaces(places)
	if err != nil {
		log.Error().Err(err).Msg("Error on creating places")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := logic.GetAndSaveNearbyPlacesByListConcurently(places)
	if err != nil {
		log.Error().Err(err).Msg("Error on getting and saving nearby places")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	totalCount, err := logic.CountPlaces()
	if err != nil {
		log.Error().Err(err).Msg("Error on counting places")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response.ActualCount = totalCount

	utils.WriteHttpJsonResponse(w, http.StatusOK, response)
}

func ListPlaces(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "0"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Error().Err(err).Msg("Error on converting offset to int")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Error().Err(err).Msg("Error on converting limit to int")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	places, err := logic.ListPlaces(offsetInt, limitInt)
	if err != nil {
		log.Error().Err(err).Msg("Error on listing places")
		utils.WriteHttpJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteHttpJsonResponse(w, http.StatusOK, places)
}
