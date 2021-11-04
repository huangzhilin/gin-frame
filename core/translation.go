package core

//gin验证器validator国际化
import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

//Trans 定义一个全局翻译器
var Trans ut.Translator

type Validators struct {
	Func        func(fl validator.FieldLevel) bool
	Translation string
}

// InitTrans 初始化翻译器  locale：zh 、en
func InitTrans(locale string, Validators map[string]Validators) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 自定义错误提示信息的字段名
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			//name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			name := strings.SplitN(fld.Tag.Get("comment"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		uni := ut.New(enT, zhT, enT)

		var ok bool
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		}

		for key, val := range Validators {
			// 注册自定义的校验方法
			if err := v.RegisterValidation(key, val.Func); err != nil {
				return err
			}
			// 注册自定义翻译方法
			if err := v.RegisterTranslation(
				key,
				Trans,
				registerTranslator(key, val.Translation),
				translate,
			); err != nil {
				return err
			}
		}

		return
	}
	return
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

//RemoveTopStruct 去除字段名中的结构体名称标识
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
