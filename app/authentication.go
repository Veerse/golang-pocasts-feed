package app

import (
	"database/sql"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"strconv"
)

// See documentation for below functions : https://github.com/appleboy/gin-jwt
// Login workflow : LoginHandler (default) -> Authenticator -> PayloadFunc -> LoginResponse (default)
// Request on secured endpoint workflow : MiddlewareFunc (default) -> IdentityHandler -> Authorizator

// Authenticator verify the user credentials given the gin context and returns a struct or a map that contains the user
// data that will be embedded in a JWT token. If an error is returned, it triggers the function Unauthorized.
func Authenticator(a *App) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		type login struct {
			Email    string `form:"email" json:"email" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}

		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		if u, err := GetUserByEmailAndPassword(loginVals.Email, loginVals.Password, &a.DB); err == nil {
			return &u, nil
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				LogError.Printf("authentificator authentification : %s", err.Error())
				return "", err
			}
		}

		return nil, jwt.ErrFailedAuthentication
	}
}

// PayloadFunc is called after a successful authentication by the Authenticator. It takes whatever was returned by the
// Authenticator and convert it into a MapClaims. The only mandatory field to put inside MapClaims is the IdentityKey
// which should correspond to the userId. We can put additional fields such as privilege, etc.
func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			"id":        v.Id,
			"privilege": v.Privilege,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler fetches the user claims within the token and pass a struct to the Authorizator. It assumes that
// MiddlewareFunc, which has been called before (see workflow) has checked that the token exists and is valid.
func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		Id:        int(claims["id"].(float64)),
		Privilege: int(claims["privilege"].(float64)),
	}
}

// Authorizator should check if the user is authorized to be reaching this endpoint (on the endpoints where the
// MiddlewareFunc applies) given the user identity value (data parameter) and the gin context. If it returns false,
// Unauthorized is triggered.
func Authorizator(a *App) func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if u, ok := data.(*User); ok {
			if u.Privilege == unverified {
				return false
			}

			if u.Privilege == admin {
				return true
			}

			if id := c.Param("podcastId"); id != "" {
				podcastId, _ := strconv.Atoi(id)
				if f, exists := a.AppCache.Podcasts[podcastId]; exists {
					if f.UserId == u.Id && (u.Privilege == poster) {
						return true
					}
				}
			}
		}
		return false
	}
}

// Unauthorized is triggered when trying to access a non-authorized resource
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
