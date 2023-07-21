package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
)

// All currently used ENV VARS
// Used in the following form: <ENV_PREFIX>_<ENV_VAR>

// Helper function to read ENV Var with optional default value
func readEnv(key string, def string) string {
	v := os.Getenv(key)

	if v == "" {
		v = def
	}

	return v
}

// Helper function to read INT parameter from the ENV
func readEnvInt(key string, def int) int {
	v := readEnv(key, fmt.Sprintf("%d", def))
	i, err := strconv.Atoi(v)
	if err != nil {
		i = def
	}
	return i
}

// My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param != "" {
			res.Errors[v.Field] = fmt.Sprintf("{%v: %v}", v.Tag, v.Param)
		} else {
			res.Errors[v.Field] = fmt.Sprintf("{key: %v}", v.Tag)
		}

	}
	return res
}

// Warp the error info in a object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

// LivenessCheck
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ReadinessCheck
func ReadinessCheck(c *gin.Context) {
	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("status", err))
		return
	}
	err = db.DB().Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("status", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}
