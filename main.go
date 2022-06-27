package main

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"reflect"
	"time"
)

const (
    base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

//生成token的参数
type UserClaims struct {
	SN    string `json:"sn"`
	//jwt-go提供的标准claim
	jwt.StandardClaims
}

var (
	//自定义的token秘钥
	secret = []byte("9N3fq7yclx6MvsgSK5uB4")
	//该路由下不校验token
	noVerify = []interface{}{"/getToken"}
	//token有效时间（纳秒）
	effectTime = 2 * time.Hour
)

// 判断obj是否在target中，target支持的类型arrary,slice,map
func IsContain(obj interface{}, target interface{}) (bool) {
    targetValue := reflect.ValueOf(target)
    switch reflect.TypeOf(target).Kind() {
    case reflect.Slice, reflect.Array:
        for i := 0; i < targetValue.Len(); i++ {
            if targetValue.Index(i).Interface() == obj {
                return true
            }
        }
    case reflect.Map:
        if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
            return true
        }
    }
  
    return false
}

// 生成token
func GenerateToken(claims *UserClaims) string {
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		panic(err)
	}
	return sign
}

//验证token
func JwtVerify(c *gin.Context) {
	//过滤是否验证token
	if IsContain(c.Request.URL.Path, noVerify) {
		return
	}
	token := c.GetHeader("token")
	if token == "" {
		panic("token not exist !")
	}
	//验证token，并存储在请求中
	c.Set("user", parseToken(token))
}

// 解析Token
func parseToken(tokenString string) *UserClaims {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("token is valid")
	}
	return claims
}

func main() {
	readYaml()
	err := initDB() // 调用输出化数据库的函数
    if err != nil {
        fmt.Printf("初始化失败！,err:%v\n", err)
        return
    }else{
        fmt.Printf("初始化成功.")
    }

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(JwtVerify)

	r.GET("/getToken", func(c *gin.Context) {
		sn := c.Request.FormValue("sn")
		secret := cfg.GetString("sn")
		if sn == secret {
			c.JSON(http.StatusOK, gin.H{
				"token": GenerateToken(&UserClaims{
				 SN:             "ABC123abc",
				 StandardClaims: jwt.StandardClaims{},
				}),
			})
		} else {
			c.String(http.StatusBadRequest, "SN不正确")
		}
	})

	r.GET("/v1/get/test", func(c *gin.Context) {
		lastUpdate := c.Request.FormValue("lastUpdate")
		if lastUpdate == "" {
			c.String(http.StatusBadRequest, "参数lastUpdate不能为空！")
			return
		}
		data, err := getJSON("select * from test where UNIX_TIMESTAMP(updateTime) > " + lastUpdate)
		if err != nil {
			fmt.Printf("获取数据失败！,err:%v\n", err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}else{
			fmt.Printf("获取数据成功.")
		}

		c.JSON(http.StatusOK, data)
	})

	r.GET("/v1/custom/sql", func(c *gin.Context) {
		secret := c.Request.FormValue("sn")
		if secret == "" {
			return
		}
		if secret != cfg.GetString("sn") {
			return
		}
		sql := c.Request.FormValue("sql")
		if sql == "" {
			c.String(http.StatusBadRequest, "参数sql不能为空！")
			return
		}
		data, err := getJSON(sql)
		if err != nil {
			fmt.Printf("获取数据失败！,err:%v\n", err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}else{
			fmt.Printf("获取数据成功.")
		}

		c.JSON(http.StatusOK, data)
	})

	r.Run(":10000")
}