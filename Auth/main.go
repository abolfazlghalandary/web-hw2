package main

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/fasthash/fnv1a"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type user_account struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Email         string `form:"email" validate:"required,email"`
	Phone_number  string `form:"phone" validate:"required,numeric"`
	Gender        string `form:"gender" validate:"required"`
	First_name    string `form:"first_name" validate:"required,alpha"`
	Last_name     string `form:"last_name" validate:"required,alpha"`
	Password_hash string `form:"password" validate:"required"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type unauthorized_token struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Token      string
	Expiration time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var identityKey = "id"
var port = "8080"

func hash(password string) string {
	return strconv.FormatUint(fnv1a.HashString64(password), 10)
}

type User struct {
	Email     string
	FirstName string
	LastName  string
}

func main() {
	r := gin.Default()
	dsn := "host=localhost user=postgres password=dev dbname=dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	var validate *validator.Validate = validator.New()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("la'anat be web"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password

			var user user_account
			db.First(&user, "email = ?", email)

			if user.Password_hash == hash(password) {
				return &User{
					Email:     user.Email,
					LastName:  user.Last_name,
					FirstName: user.First_name,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")

	auth.POST("/signin", authMiddleware.LoginHandler)

	auth.POST("/signup", func(c *gin.Context) {
		var user user_account
		c.Bind(&user)
		fmt.Println(user)
		err := validate.Struct(user)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatusJSON(400, gin.H{"field_error": err.Error()})
		} else {
			if !(user.Gender == "M" || user.Gender == "F") {
				c.AbortWithStatusJSON(400, gin.H{"field_error": "Invalid Gender, must be M or F"})
			} else {
				user_record := user_account{
					Email:         user.Email,
					Phone_number:  user.Phone_number,
					First_name:    user.First_name,
					Last_name:     user.Last_name,
					Gender:        user.Gender,
					Password_hash: hash(user.Password_hash),
				}
				var error = db.Create(&user_record).Error
				fmt.Println(error)
				if error != nil {
					c.AbortWithStatusJSON(400, gin.H{"register_error": error})
				} else {
					authMiddleware.LoginHandler(c)
				}
			}
		}
	})

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/user_info", func(c *gin.Context) {
			user, _ := c.Get(identityKey)
			var user_info user_account
			db.First(&user_info, "email = ?", user.(*User).Email)
			user_info.Password_hash = ""
			c.JSON(200, gin.H{
				"user": user_info,
			})
		})

		auth.POST("/sign_out", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			fmt.Println(claims)
			token := unauthorized_token{
				Token:      fmt.Sprintf("%s", claims),
				Expiration: time.Now(),
			}
			fmt.Println(token)
			var error = db.Create(&token)
			fmt.Println(error)
			c.JSON(200, gin.H{
				"result": "signed out",
			})
		})
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
