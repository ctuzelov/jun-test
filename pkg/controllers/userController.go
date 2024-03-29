package controllers

import (
	"auth-service/pkg/database/mongodb/models"
	"auth-service/pkg/database/mongodb/repository"
	"auth-service/pkg/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Function that extracts user data and validates it then creates a user with it if valid.
func SignUp(c *gin.Context) (err error) {
	var user models.User

	if err = c.BindJSON(&user); err != nil {
		return
	}

	count, err := repository.CountUsersByEmail(user.Email) // Check if user email already exists
	if err != nil {
		err = errors.New("error occured while checking the email")
		return
	}

	if count > 0 {
		err = errors.New("this email already exists")
		return
	}

	password := utils.HashPassword(user.Password) // Hash password for the user
	user.Password = password
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, err := utils.GenerateAllTokens(user.Email, user.Name, user.User_id)

	if err != nil {
		err = errors.New("error occured generating tokens")
		return
	}
	user.Token = token
	user.Refresh_Token = refreshToken
	err = repository.CreateUser(user)

	return

}

// Function that extracts user data from context and compares it to database if its valid, Tokens are regenrated for the user and are returned
func Login(c *gin.Context) (token string, refreshToken string, err error) {
	var user models.User
	if err = c.BindJSON(&user); err != nil {
		return
	}

	foundUser, err := repository.GetUserByEmail(user.Email)

	if err != nil {
		return
	}
	passwordIsValid, msg := utils.VerifyPassword(user.Password, foundUser.Password)

	if !passwordIsValid {
		err = errors.New(msg)
		return
	}
	token, refreshToken, err = utils.GenerateAllTokens(foundUser.Email, foundUser.Name, foundUser.User_id)

	if err != nil {
		return
	}
	_, _, err = utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	return

}

// Function that accepts context, extracts refreshToken from it and returns new authorization token and refresh token.
func Refresh_Token(c *gin.Context) (token string, refreshToken string, err error) {
	// temporary struct to capture the refresh token from context
	type Input struct {
		Token_ string `json:"refresh_token"`
	}

	var in Input
	if err = c.BindJSON(&in); err != nil {
		return
	}

	refreshToken = in.Token_
	claims, err := utils.ValidateToken(refreshToken)

	if err != nil {
		return
	}
	//Capture user details extracted from the refresh token provided
	name, email, uid := claims["Name"].(string), claims["Email"].(string), claims["Uid"].(string)
	token, refreshToken, err = utils.GenerateAllTokens(name, email, uid)

	if err != nil {
		return
	}
	token, refreshToken, err = utils.UpdateAllTokens(token, refreshToken, uid)

	return

}

func EditPersonalData(c *gin.Context) error {
	var userData models.User
	// Bind the updated user data from the request body
	if err := c.BindJSON(&userData); err != nil {
		return err
	}

	claims, ok := c.Get("user_email")
	if !ok {
		return errors.New("unable to retrieve user email from JWT claims")
	}
	userEmail := claims.(string)

	// Ensure that the user making the request is authorized to edit their own data
	if userData.Email != userEmail {
		return errors.New("unauthorized to edit data for this user")
	}

	// Check if the user with the provided email exists
	existingUser, err := repository.GetUserByEmail(userData.Email)
	if err != nil {
		return err
	}
	if existingUser == (models.User{}) {
		return errors.New("user not found")
	}

	// Update the user's personal data
	if userData.Name != "" {
		existingUser.Name = userData.Name
	}
	if userData.Number != "" {
		existingUser.Number = userData.Number
	}
	if userData.DateOfBirth != "" {
		existingUser.DateOfBirth = userData.DateOfBirth
	}

	// Save the updated user data to the database
	if err := repository.UpdateUser(existingUser); err != nil {
		return err
	}

	return nil
}

func GetAllUsers() ([]models.User, error) {
	users, err := repository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}