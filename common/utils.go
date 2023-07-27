package common

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Context keys
const (
	VAULT_TOKEN_HEADER = "X-Vault-Token"
	VAULT_TOKEN        = "vaultToken"
	VAULT_TOKEN_MODEL  = "vaultTokenModel"
	SESSION_POLICIES   = "sessionPolicies"
	PATH_CAPABILITIES  = "pathCapabilities"
	IS_ROOT            = "isRoot"
)

// All currently used ENV VARS
// Used in the following form: <ENV_PREFIX>_<ENV_VAR>

func TrimPrefix(s string, pref string) string {
	return strings.TrimPrefix(s, pref)
}

func GetRequestPath(c *gin.Context) string {
	path := c.FullPath()
	return TrimPrefix(path, "/v1/")
}

// Helper function to check if file exists
func fileExists(path string) bool {
	res := false

	if path == "" {
		return res
	}

	_, err := os.Stat(path)
	if err == nil {
		res = true
	}

	return res
}

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

// Helper function to read BOOL parameter from the ENV
func readEnvBool(key string, def bool) bool {
	v := readEnv(key, fmt.Sprintf("%t", def))
	b, err := strconv.ParseBool(v)
	if err != nil {
		b = def
	}
	return b
}

// My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors []string `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
// func NewValidatorError(err error) CommonError {
// 	res := CommonError{}
// 	res.Errors = make(map[string]interface{})
// 	errs := err.(validator.ValidationErrors)
// 	for _, v := range errs {
// 		// can translate each error one at a time.
// 		//fmt.Println("gg",v.NameNamespace)
// 		if v.Param != "" {
// 			res.Errors[v.Field] = fmt.Sprintf("{%v: %v}", v.Tag, v.Param)
// 		} else {
// 			res.Errors[v.Field] = fmt.Sprintf("{key: %v}", v.Tag)
// 		}
//
// 	}
// 	return res
// }

// Warp the error info in a object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = []string{
		fmt.Sprintf("%s: %s", key, err.Error()),
	}
	return res
}

func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

func EncToB64(str string) string {
	src := []byte(str)
	res := base64.StdEncoding.EncodeToString(src)
	return res
}

func DecFromB64(str string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	res := string(b)
	return res, nil
}

// Verify create acces on individual path
func VerifyCreateAccess(c *gin.Context) bool {
	l, _ := GetLogger()
	isRoot := c.MustGet(IS_ROOT).(bool)
	l.Debug("VerifyCreateAccess", "isRoot", isRoot)
	if isRoot {
		return true
	}
	capabilities := c.MustGet(PATH_CAPABILITIES).(map[string]bool)
	return capabilities["create"]
}
