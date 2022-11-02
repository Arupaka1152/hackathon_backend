package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
	r := new(SignupReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	password := string(hash)

	accountId := utils.GenerateId()
	newAccount := model.Account{
		Id:       accountId,
		Name:     r.Name,
		Email:    r.Email,
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
	r := new(LoginReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetAccount := model.Account{}
	if err := dao.FindAccountByEmail(&targetAccount, r.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "account not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(targetAccount.Password), []byte(r.Password)); err != nil {
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
