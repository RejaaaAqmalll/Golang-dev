package auth

import (
	"net/http"
	"nyoba/configg"
	"nyoba/helper"
	"nyoba/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	// Mendapatkan email&passowrd
	var body struct {
		Name     string
		Email    string
		Password string
	}
	// Binding / mengikat
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot Binding data",
		})
		return
	}
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Cannot hash password",
		})
		return
	}
	// Insert User
	user := models.Users{Name: body.Name, Email: body.Email, Password: string(hash)}

	result := configg.KoneksiData().Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Cannot Create User",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Data Succes Insert",
	})
}

func Login(c *gin.Context) {
	// Ambil email dan password dari body

	// var untuk menangkap input email&password
	var body struct {
		Email    string
		Password string
	}
	// Binding / mengikat
	if c.ShouldBindJSON(&body) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot Binding data",
			// "status":  false,
		})
		return
	}
	// cek request
	var user models.Users
	err1 := configg.KoneksiData().Debug().First(&user, "email = ?", body.Email).Error

	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid User or Password",
			// "status":  false,
			// "error":   err1.Error(),
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "User Not found",
			// "status":  false,
		})
		return
	}
	// Membandigkan password user dengan password pada database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid User or Password",
			// "status":  false,
			// "error":   err.Error(),
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_user":    user.ID,
		"email_user": user.Email,
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Cannot get jwt token",
		})
		return
	}

	// Respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorizion", tokenString, 3600*24*30, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "login success",
		"token":   tokenString,
	})
}

func Validate(c *gin.Context) {
	// Validasi telah melakukan login

	// Memanggil data user yang telah melakukan login
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"email":    user.(models.Users).Email,
		"password": user.(models.Users).Password,
	})
}

func Logout(c *gin.Context) {

	// name of cookie
	c.SetCookie("Authorizion", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "you are Logout",
	})
}

func ForgotPassword(c *gin.Context) {

	var Forgotpw struct {
		Email string
	}
	//  Binding data
	err := c.ShouldBindJSON(&Forgotpw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Can't Bind data",
		})
		return
	}

	// check email yang masuk
	var user models.Users
	if err := configg.KoneksiData().Debug().Where("email = ?", Forgotpw.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid email",
		})
		return
	}

	// Insert a code verify
	token := models.Users{Code: helper.GenerateCode()}
	if err1 := configg.KoneksiData().Debug().Model(&user).Where("email = ?", Forgotpw.Email).Update(&token).Error; err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Failed update code",
		})
		return
	}

	//  Send code to email
	if helper.Sendemail(Forgotpw.Email, token.Code) != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Failed Send email",
		})
	}

	// Respon Final
	c.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"message":   "Please Check your Spam Email",
		"your_code": token.Code,
	})

}

func Authyverify(c *gin.Context) {
	var codeInput struct {
		Code string
	}
	if err := c.ShouldBindJSON(&codeInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Can't Binding data code",
		})
		return
	}

	var user models.Users

	// Check input code
	if codeInput.Code == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "The code is null",
			"status":  false,
		})
		return
	}

	if err := configg.KoneksiData().Debug().First(&user, "code = ?", codeInput.Code).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Invalid Code",
			"status":  false,
		})
		return
	}

	// Check Input Code
	if codeInput.Code != user.Code {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Invalid Code",
			"status":  false,
		})
		return
	}

	// Respond Final
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Code Passed",
		"status":  true,
	})

}

func ResetPassword(c *gin.Context) {
	var inputPass struct {
		Email           string
		NewPassword     string
		ConfirmPassword string
	}

	// Binding Input
	if err := c.ShouldBindJSON(&inputPass); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Can't Bind data",
			"status":  false,
		})
		return
	}

	// Compare new pass with confirm pass
	if inputPass.NewPassword != inputPass.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password must be the same",
			"status":  false,
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(inputPass.ConfirmPassword), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "failed hash password",
			"status":  false,
		})
		return
	}

	var user models.Users
	newpass := models.Users{Password: string(hash)}
	if err := configg.KoneksiData().Model(&user).Where("email = ?", inputPass.Email).Update(&newpass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Invalid Password",
			"status":  false,
		})
		return
	}

	//  Respond Final
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Password Succes Change",
		"status":  true,
	})
}
