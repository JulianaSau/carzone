package login

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/JulianaSau/carzone/driver"
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

	// valid := (credentials.UserName == "admin" && credentials.Password == "admin123")

	// if !valid {
	// 	http.Error(w, "Incorrect Username or Password", http.StatusUnauthorized)
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// Get the user from the database by username
	user, err := getUserByUsername(credentials.UserName)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error fetching user: ", err)
		return
	}

	// Check if the password is correct
	if err := user.CheckPassword(credentials.Password); err != nil {
		http.Error(w, "Incorrect Username or Password", http.StatusUnauthorized)
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

// getUserByUsername retrieves a user from the database based on the username.
func getUserByUsername(username string) (*models.User, error) {
	// Prepare the SQL query to find a user by username
	query := `
		SELECT username, password, first_name, last_name, email, phone_number, role, id, active, created_at, updated_at
		FROM "user"
		WHERE username = $1
	`
	var user models.User

	// Use the database connection to query the database
	err := driver.GetDB().QueryRow(query, username).Scan(
		&user.UserName,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNumber,
		&user.Role,
		&user.ID,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// If no rows are found, return an error indicating that the user does not exist
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	// Log any other error that occurs during the query
	if err != nil {
		log.Println("Error retrieving user:", err)
		return nil, err
	}

	// Return the user object
	return &user, nil
}

// getUserByUsername retrieves a user from the database based on the username.
