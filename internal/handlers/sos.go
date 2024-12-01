package handlers

import (
	"encoding/json"
	"net/http"
	"sos-crud/internal/database"
	"sos-crud/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateSos(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "Invalid userID", http.StatusInternalServerError)
		return
	}

	var user models.User

	var sos models.Sos

	if err := json.NewDecoder(r.Body).Decode(&sos); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	sos.UserID = userID

	if err := database.DB.Create(&sos).Error; err != nil {
		http.Error(w, "Unable to crate sos", http.StatusInternalServerError)
		return
	}

	// Can't using preload, so I after load
	sos.User = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sos)
}

func GetSosByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var sos models.Sos

	if err := database.DB.Preload("User").First(&sos, id).Error; err != nil {
		http.Error(w, "Sos not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sos)
}

func UpdateSosByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "Invalid userID", http.StatusInternalServerError)
		return
	}

	var sos models.Sos

	if err := database.DB.First(&sos, id).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&sos); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if sos.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := database.DB.Save(&sos).Error; err != nil {
		http.Error(w, "Problem when edit sos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sos)
}

func DeleteSosByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "Invalid userID", http.StatusInternalServerError)
		return
	}

	var sos models.Sos

	if err := database.DB.First(&sos, id).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if sos.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := database.DB.Delete(&sos, id).Error; err != nil {
		http.Error(w, "Problem when deleted sos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Sos deleted successful"})
}

func GetAllSos(w http.ResponseWriter, r *http.Request) {

	pageNumberStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("size")

	if pageNumberStr == "" {
		pageNumberStr = "1"
	}

	if pageSizeStr == "" {
		pageSizeStr = "2"
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)

	if err != nil || pageNumber < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)

	if err != nil || pageSize < 1 {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	offset := (pageNumber - 1) * pageSize

	var soss []models.Sos

	if err := database.DB.Preload("User").Offset(offset).Limit(pageSize).Find(&soss).Error; err != nil {
		http.Error(w, "Unable to fetch Soss", http.StatusInternalServerError)
		return
	}

	var totalSosCount int64

	err = database.DB.Model(&models.Sos{}).Count(&totalSosCount).Error

	if err != nil {
		http.Error(w, "Error when fetch count total sos", http.StatusInternalServerError)
		return
	}

	totalPages := (totalSosCount + int64(pageSize) - 1) / int64(pageSize)

	response := map[string]interface{}{
		"sosPosts": soss,
		"pagination": map[string]interface{}{
			"currentPage":  pageNumber,
			"pageSize":     pageSize,
			"totalPages":   totalPages,
			"totalRecords": totalSosCount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
