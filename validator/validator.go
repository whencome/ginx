package validator

import (
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zt "github.com/go-playground/validator/v10/translations/zh"
    "log"
    "reflect"
)

var (
    trans         ut.Translator
    valid         *validator.Validate
    showAllErrors = false
    errSeparator  = "\n"
)

// 初始化时自动注册中文解释器
func init() {
    translator := zh.New()
    trans, _ = ut.New(translator, translator).GetTranslator("zh")
    valid = binding.Validator.Engine().(*validator.Validate)
    zt.RegisterDefaultTranslations(valid, trans)
    // 注册一个函数，获取struct tag里自定义的label作为字段名
    valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := fld.Tag.Get("label")
        return name
    })
    valid.RegisterTranslation("json", trans, func(ut ut.Translator) error {
        return ut.Add("json", "{0}不是一个有效的json字符串", true)
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("json", fe.Field(), fe.Field())
        return t
    })
}

// ErrorTranslator 错误解释器，用于多语言环境的错误处理
type ErrorTranslator interface {
    RegisterTranslations(v *validator.Validate) (ut.Translator, error)
}

// UseTranslator 注册一个自定义的错误解释器，主要用于非汉语环境
func UseTranslator(et ErrorTranslator) {
    defer func() {
        if r := recover(); r != nil {
            log.Println(r)
        }
    }()
    t, err := et.RegisterTranslations(valid)
    if err != nil {
        log.Printf("register custom translations failed: %s", err)
        return
    }
    if trans != nil {
        trans = t
    }
}

// ShowFullError 当有多个错误时，是否全部显示，默认只显示一个错误
func ShowFullError(b bool) {
    showAllErrors = b
}

// SetErrSeparator 设置错误分割符，当showAllErrors为true时，将以此符号分割
func SetErrSeparator(s string) {
    if s == "" {
        s = "\n"
    }
    errSeparator = s
}

// Translate 翻译错误消息
func Translate(err error) string {
    var result string
    errs := err.(validator.ValidationErrors)
    for _, err := range errs {
        result += err.Translate(trans) + ";"
    }
    return result
}

// Error 获取并翻译错误信息
func Error(err error) string {
    defer func() {
        if r := recover(); r != nil {
            log.Println(r)
        }
    }()
    errs, ok := err.(validator.ValidationErrors)
    if !ok {
        return err.Error()
    }
    msg := ""
    for _, err := range errs {
        if msg != "" {
            msg += errSeparator
        }
        msg += err.Translate(trans)
        if !showAllErrors {
            break
        }
    }
    return msg
}
