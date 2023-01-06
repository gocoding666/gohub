// Package jwt处理JWT认证
package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"strings"
	"time"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvaild           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

// JWT定义一个jwt对象
type JWT struct {
	//密钥，用以加密JWT,读取配置信息app.key
	SignKey []byte
	//刷新Token的最大过期时间
	MaxRefresh time.Duration
}

// JWTCustomClaims 自定义载荷
type JWTCustomClaims struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	ExpireAtTime int64  `json:"expire_time"`
	// StandardClaims 结构体实现了Claims接口继承了Valid()方法
	//JWT 规定了7个官方字段，提示使用：
	//- iss(issusr):发布者
	// - sub(subject):主题
	//- iat(Issued At):生成签名的时间
	// - exp(expiration time):签名过期时间
	// - aud(audience):观众，相当于接收者
	// - nbf (Not Before):生效时间
	// - jti(JWT ID):编号
	jwtpkg.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("app.key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

func (jwt *JWT) ParserToken(c *gin.Context) (*JWTCustomClaims, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}
	//1.调用jwt库解析用户传参的Token
	token, err := jwt.parseTokenString(tokenString)
	//2.解析出错
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}

		}
		return nil, ErrTokenInvaild
	}
	//3.将token中的cliams信息解析出来和JWTCustiomClaims数据结构进行校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvaild
}

// RefreshToken更新token,用以提供refresh token接口
func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {
	// 1.从Header里获取token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}
	//2.调用jwt库解析用户传参的token
	token, err := jwt.parseTokenString(tokenString)
	//3.解析出错，未报错证明是合法的Token(甚至未到过期时间)
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		//满足refresh的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}
	//4.解析JWTCustonCliams的数据
	claims := token.Claims.(*JWTCustomClaims)
	//5.检查是否过了【最大允许刷新的时间】
	x := app.TimenowInTimezone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		//修改过期时间
		claims.StandardClaims.ExpiresAt = jwt.expireAtTime()
		return jwt.createToken(*claims)
	}
	return "", ErrTokenExpiredMaxRefresh
}

// IssueToken 生成Token,在登录成功时调用
func (jwt *JWT) IssueToken(userID string, userName string) string {
	//1.构造用户claims信息（负荷）
	expireAtTime := jwt.expireAtTime()
	claims := JWTCustomClaims{
		userID,
		userName,
		expireAtTime,
		jwtpkg.StandardClaims{
			NotBefore: app.TimenowInTimezone().Unix(), //签名生效时间
			IssuedAt:  app.TimenowInTimezone().Unix(), //首次签名时间（后续刷新Token不会更新）
			ExpiresAt: expireAtTime,                   //签名过期时间（后续刷新Token不会更新）
			Issuer:    config.GetString("app.name"),   //签名颁发者
		},
	}
	//2. 根据claims生成token对象
	token, err := jwt.createToken(claims)
	if err != nil {
		logger.LogIf(err)
		return ""
	}
	return token
}

// createToken创建Token,内部使用，外部请调用IssueToken
func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {
	//使用HS256算法进行token生成
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}

// expireAtTime过期时间
func (jwt *JWT) expireAtTime() int64 {
	timenow := app.TimenowInTimezone()
	var expireTime int64
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
	} else {
		expireTime = config.GetInt64("jwt.expire_time")
	}
	expire := time.Duration(expireTime) * time.Minute
	return timenow.Add(expire).Unix()
}
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	//按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}
