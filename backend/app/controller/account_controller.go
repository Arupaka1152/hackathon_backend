package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode/utf8"
)

type SignupReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context) {
	req := new(SignupReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if utf8.RuneCountInString(req.Password) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "length of password must be more than 9"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	password := string(hash)

	accountId := utils.GenerateId()
	newAccount := model.Account{
		Id:       accountId,
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	if err := dao.CreateAccount(&newAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authentication": token})
}

func Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetAccount := model.Account{}
	if err := dao.FindAccountByEmail(&targetAccount, req.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "account not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(targetAccount.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid password"})
		return
	}

	token, err := auth.GenerateToken(targetAccount.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authentication": token})
}
