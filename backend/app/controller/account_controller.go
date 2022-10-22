package controller

import (
	"backend/app/dao"
	"backend/app/model"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
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
	req := new(SignupReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	accountId := ulid.Make().String()
	newAccount := model.Account{
		Id:       accountId,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, //パスワードはハッシュ化する！！
	}

	if err := dao.CreateAccount(&newAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	//トークンを生成、トークンをヘッダーに入れて送信
}

func Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetAccount := model.Account{}
	if err := dao.FetchAccountByEmail(&targetAccount, req.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
	}

	//パスワード（ハッシュ化したもの）を比較して一致しなければhttp.StatusUnauthorizedを返す

	//トークンを生成、トークンをヘッダーに入れて送信
}
