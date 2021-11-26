package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"time"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			ResponseError(c, InternalErrorCode, errors.New("请求未携带token，无权限访问"))
			c.Abort()
			return
		}
		j := NewJWT()
		clamis, err := j.ParserToken(token)
		if err != nil {
			if err == TokenExpired {
				ResponseError(c, InternalErrorCode, errors.New("授权已过期"))
				c.Abort()
				return
			}
			ResponseError(c, InternalErrorCode, errors.New(err.Error()))
			c.Abort()
			return
		}
		c.Set("claims", clamis)
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "newtrekWang"
)


func NewJWT() *JWT{
	return &JWT{
		[]byte(GetSignKey()),
	}
}

type CustomClaims struct {
	Id int `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error)  {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}


func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error){
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok{
			if ve.Errors & jwt.ValidationErrorMalformed != 0{
				return nil, fmt.Errorf("token 不可用")
			}else if ve.Errors & jwt.ValidationErrorExpired != 0{
				return nil, fmt.Errorf("token 过期")
			}else if ve.Errors & jwt.ValidationErrorNotValidYet != 0{
				return nil, fmt.Errorf("无效的token")
			}else{
				return nil, fmt.Errorf("token不可用")
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid{
		return claims, nil
	}
	return nil, fmt.Errorf("token 无效")
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}