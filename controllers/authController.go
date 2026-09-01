package controllers

import (
	"linkedveil/initializers"
	"linkedveil/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) { 

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var foundUser models.User
	initializers.DB.Where("username=?", authInput.Username).First(&foundUser)

	if foundUser.ID != 0 { 
		c.JSON(http.StatusBadRequest, gin.H{"error" : "username already exists: " + foundUser.Username })
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return 
	}
	
	user := models.User{
		Username: authInput.Username,
		Password: string(passwordHash),
		Role: models.UserRole,
	}

	if err:= initializers.DB.Create(&user).Error; err != nil { 
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func LoginUser(c *gin.Context) { 

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}

	var userFound models.User
	initializers.DB.Where("username = ?", authInput.Username).First(&userFound)

	if userFound.ID == 0 { 
		c.JSON(http.StatusBadRequest, gin.H{"error" : "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error" : "wrong password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id" : userFound.ID,
		"role" : userFound.Role,
		"exp" : time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate JWT token"})
		return
	}

	c.JSON(200, gin.H{
		"token" : token,
		"user": gin.H{
			"id" : userFound.ID,
			"username" : userFound.Username,
			"role" : userFound.Role,
		},
	})
}

func GetUserProfile(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists { 
		c.JSON(http.StatusUnauthorized, gin.H{"error" : "user not found"})
		return
	}
	u := user.(models.User)
	c.JSON(200, gin.H{
		"id": u.ID,
		"username" : u.Username,
		"role" : u.Role,
	})
}