package auth

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateNewUser(registerDTO *RegisterDTO) (int, error) {
	mySql := mysql.MysqlDB()
	defer mySql.Close()

	stmtAuth, errAuth := mySql.Prepare("INSERT INTO auth(username,email,pwd) VALUES(?,?,?)")

	if errAuth != nil {
		return 0, errAuth
	}
	defer stmtAuth.Close()

	hash, _ := HashPassword(registerDTO.Password)

	resAuth, errExecAuth := stmtAuth.Exec(registerDTO.Username, registerDTO.Email, hash)

	if errExecAuth != nil {
		return 0, errExecAuth
	}

	lastInsetedIdAuth, _ := resAuth.LastInsertId()

	return int(lastInsetedIdAuth), nil

}

func AuthRegister(c *gin.Context) {
	var registerDTO RegisterDTO

	if err := c.ShouldBindWith(&registerDTO, binding.JSON); err == nil {

		userId, err := CreateNewUser(&registerDTO)

		if err != nil {
			handler.ErrorHandler(c, http.StatusBadRequest, "error while creating account. please try again", err)
			return
		}

		accessToken, err := CreateJwtToken(userId, registerDTO.Username)

		if err != nil {
			handler.ErrorHandler(c, http.StatusInternalServerError, "some errro occured on server side", err)
			return
		}

		handler.SuccessHandler(c, http.StatusCreated, "account created", accessToken)

	} else {
		c.JSON(400, gin.H{"error": err.Error()})
	}

}

func AuthLogin(c *gin.Context) {
	var loginDTO LoginDTO
	if err := c.ShouldBindWith(&loginDTO, binding.JSON); err == nil {

		mySql := mysql.MysqlDB()
		defer mySql.Close()

		var auth Auth
		errorAuthLogin := mySql.QueryRow("SELECT userId,username,email,pwd FROM auth WHERE email=? ", loginDTO.Email).Scan(&auth.UserId, &auth.Username, &auth.Email, &auth.Password)

		if errorAuthLogin != nil {
			handler.ErrorHandler(c, http.StatusInternalServerError, "Invalid Credentails. Account Not Found.", errorAuthLogin)
			return
		}

		if CheckPasswordHash(loginDTO.Password, auth.Password) {

			accessToken, err := CreateJwtToken(auth.UserId, auth.Username)

			if err != nil {
				handler.ErrorHandler(c, http.StatusInternalServerError, "some errro occured on server side", err)
				return
			}

			handler.SuccessHandler(c, http.StatusOK, "account login successfully", accessToken)
			return

		} else {
			handler.ErrorHandler(c, http.StatusBadRequest, "Invalid Credentails", fmt.Errorf("Invalid Passowrd"))
			return
		}

	}

	handler.ErrorHandler(c, http.StatusBadRequest, "error in login account", fmt.Errorf(""))

}
