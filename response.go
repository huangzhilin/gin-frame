package gin_frame

import (
	"errors"
	"net/http"
	"strings"

	"github.com/huangzhilin/gin-frame/core"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//REST API 统一响应格式规范，0表示失败，1表示成功
const (
	ERROR   = 0
	SUCCESS = 1
)

type Response struct {
	Code int         `json:"code"` //0错误 1成功
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(code int, data interface{}, msg string) *Response {
	return &Response{
		code,
		data,
		msg,
	}
}

func Ok() *Response {
	return Result(SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithMessage(message string) *Response {
	return Result(SUCCESS, map[string]interface{}{}, message)
}

func OkWithData(data interface{}) *Response {
	return Result(SUCCESS, data, "操作成功")
}

func OkWithDetailed(data interface{}, message string) *Response {
	return Result(SUCCESS, data, message)
}

func Fail() *Response {
	return Result(ERROR, map[string]interface{}{}, "操作失败")
}

func FailWithMessage(message string) *Response {
	return Result(ERROR, map[string]interface{}{}, message)
}

func FailWithDetailed(data interface{}, message string) *Response {
	return Result(ERROR, data, message)
}

//Convert 在路由中使用
func Convert(f func(*gin.Context) *Response) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := f(c)
		c.JSON(http.StatusOK, resp)
	}
}

//BindValidParam  绑定and验证请求数据
//type AdminLoginRequest struct {
//	UserName string `form:"username" json:"username" comment:"用户名" binding:"required"`
//	Password string `form:"password" json:"password" comment:"密码"  binding:"required"`
//}
func BindValidParam(c *gin.Context, param interface{}) error {
	if err := c.ShouldBind(param); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok { // validator.ValidationErrors类型错误则进行翻译
			mapErrs := core.RemoveTopStruct(errs.Translate(core.Trans))
			sliceErrs := []string{}
			for _, e := range mapErrs {
				sliceErrs = append(sliceErrs, e)
			}
			return errors.New(strings.Join(sliceErrs, ";"))
		} else { // 非validator.ValidationErrors类型错误直接返回
			return errors.New("shouldBind err:" + err.Error())
		}
	}
	return nil
}
