package handlers

import (
	"encoding/json"
	"net/http"
	"sos-crud/internal/database"
	"sos-crud/internal/models"
	"sos-crud/internal/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid inut", http.StatusBadRequest)
		return
	}

	var existingUser models.User

	if err := database.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		http.Error(w, "username already taken", http.StatusBadRequest)
		return
	}

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: input.Username,
		Password: hashPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)

}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var foundUser models.User

	if err := database.DB.Where("username = ?", input.Username).First(&foundUser).Error; err != nil {
		http.Error(w, "User name or password not correct", http.StatusBadRequest)
		return
	}

	isSame := utils.CheckPasswordHash(input.Password, foundUser.Password)

	if !isSame {
		http.Error(w, "User name or password not correct", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateJwT(foundUser)

	if err != nil {
		http.Error(w, "Something when wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}

	var user models.User

	if err := database.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
