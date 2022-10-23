package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SignupReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	r := new(SignupReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	password := string(hash)

	accountId := ulid.Make().String()
	newAccount := model.Account{
		Id:       accountId,
		Name:     r.Name,
		Email:    r.Email,
		Password: password,
	}

	if err := dao.CreateAccount(&newAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	token, err := auth.GenerateToken(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"authentication": token})
}

func Login(c *gin.Context) {
	r := new(LoginReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetAccount := model.Account{}
	if err := dao.FetchAccountByEmail(&targetAccount, r.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "account not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(targetAccount.Password), []byte(r.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid password"})
	}

	token, err := auth.GenerateToken(targetAccount.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"authentication": token})
}
