package login

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JulianaSau/carzone/models"
	"github.com/golang-jwt/jwt/v4"
)

// LoginHandler godoc
// @Summary Authenticate user and generate a JWT token
// @Description Validates user credentials and returns a JWT token on success
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid credentials"
// @Router /api/v1/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error decoding credentials: ", err)
		return
	}

	valid := (credentials.UserName == "admin" && credentials.Password == "admin123")

	if !valid {
		http.Error(w, "Incorrect Username or Password", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// generate token
	tokenString, err := GenerateToken(credentials.UserName)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error generating token: ", err)
		return
	}

	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GenerateToken(username string) (string, error) {
	// implement token generation logic here
	expiration := time.Now().Add(24 * time.Hour)

	// claims := &jwt.RegisteredClaims{
	// 	ExpiresAt: jwt.NewNumericDate(expiration),
	// 	Subject:   username,
	// 	IssuedAt:  jwt.NewNumericDate(time.Now()),
	// }
	claims := &jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
