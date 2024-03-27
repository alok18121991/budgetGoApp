package controllers

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/models"
	"alok/web-service-budget/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

func GetAuthUserDetails(c echo.Context) error {
	code := c.QueryParam("code")

	// Exchange code for access token
	tokenResp, err := exchangeCodeForToken(code)
	// fmt.Println("token response.....", tokenResp)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to exchange code for token")
	}

	// Get user info
	userInfo, err := getUserInfo(tokenResp.AccessToken)
	// fmt.Println("userInfo response.....", userInfo)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get user info")
	}

	dbUser, err := services.GetUserByEmail(c, userInfo.Email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error occured while fetching user from db")
	}

	// // Create session for user
	// session, err := configs.Store.Get(c.Request(), configs.SessionName)
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, "Failed to create session")
	// }

	// // Store user info in session
	// session.Values["userID"] = dbUser.Id.Hex()
	// // session.Values["userName"] = dbUser.FirstName
	// // session.Values["userEmail"] = dbUser.Email
	// err = session.Save(c.Request(), c.Response())
	// if err != nil {
	// 	fmt.Println("session error.....", err)
	// 	return c.String(http.StatusInternalServerError, "Failed to save session")
	// }

	// // Respond with user info
	// return c.JSON(http.StatusOK, dbUser)

	// Generate JWT token
	// Token expires in 24 hours
	// tokenString, err := generateToken(dbUser.Email)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate JWT token")
	// }

	// // Return JWT token in response
	// return c.JSON(http.StatusOK, echo.Map{"token": tokenString, "userData": dbUser})

	// Generate JWT token and refresh token for the user
	expirationTime := time.Now().Add(5 * time.Minute) // Access token expires in 5 minutes
	accessToken, err := generateToken(dbUser, expirationTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := generateRefreshToken(dbUser.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token")
	}

	// Return JWT tokens in response
	return c.JSON(http.StatusOK, echo.Map{"accessToken": accessToken, "refreshToken": refreshToken, "userData": dbUser})

}

func generateRefreshToken(userEmail string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
	claims := &configs.Claims{
		User: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configs.JwtRefreshSecret)
}

func generateToken(user *models.User, expirationTime time.Time) (string, error) {
	claims := &configs.TokenClaims{
		Id:   user.Id.Hex(),
		User: user.Email,
		Name: user.FirstName + user.LastName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configs.JwtSecret)
}

// Refresh token handler
func RefreshTokenHandler(c echo.Context) error {
	// Parse refresh token from request body
	type RefreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	var refreshRequest RefreshRequest
	if err := c.Bind(&refreshRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Parse and validate refresh token
	token, err := jwt.Parse(refreshRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.JwtRefreshSecret, nil
	})
	if err != nil {
		fmt.Println("unexpected signing method:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	// Check if the token is valid and not expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract username from claims
		user := claims["user"].(string)
		// Generate a new access token for the user
		dbUser, err := services.GetUserByEmail(c, user)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error occured while fetching user from db")
		}
		expirationTime := time.Now().Add(5 * time.Minute)
		newAccessToken, err := generateToken(dbUser, expirationTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate JWT token")
		}
		// Return new access token in response
		return c.JSON(http.StatusOK, echo.Map{"accessToken": newAccessToken})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
}

func exchangeCodeForToken(code string) (*GoogleTokenResponse, error) {

	// fmt.Println("token response.....", configs.GetEnvValueFor("clientID"))

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", configs.GetEnvValueFor("clientID"))
	data.Set("client_secret", configs.GetEnvValueFor("clientSecret"))
	data.Set("redirect_uri", configs.GetEnvValueFor("redirectURI"))
	data.Set("grant_type", "authorization_code")

	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	// fmt.Println("exchangeCodeForToken response.....", resp)
	// fmt.Println("exchangeCodeForToken response.....22", strings.NewReader(data.Encode()))

	defer resp.Body.Close()

	var tokenResp GoogleTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
