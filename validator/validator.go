package validator

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zt "github.com/go-playground/validator/v10/translations/zh"
	"github.com/whencome/ginx/log"
)

var (
	trans         ut.Translator
	valid         *validator.Validate
	showAllErrors = false
	errSeparator  = "\n"
	mu            sync.RWMutex // protect concurrent access to trans
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

// UseTranslator register a custom error translator, mainly for non-Chinese environments
func UseTranslator(et ErrorTranslator) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered from: %v", r)
		}
	}()
	t, err := et.RegisterTranslations(valid)
	if err != nil {
		log.Errorf("register custom translations failed: %s", err)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if trans != nil {
		trans = t
	}
}

// ShowFullError when there are multiple errors, whether to display all, default shows only one error
func ShowFullError(b bool) {
	showAllErrors = b
}

// SetErrSeparator set error separator, used when showAllErrors is true
func SetErrSeparator(s string) {
	if s == "" {
		s = "\n"
	}
	errSeparator = s
}

// Translate translate error message (deprecated, use Error instead)
// Deprecated: Use Error() instead
func Translate(err error) string {
	var result string
	errs := err.(validator.ValidationErrors)
	for _, err := range errs {
		result += err.Translate(trans) + ";"
	}
	return result
}

// Error get and translate error information
func Error(err error) string {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered from: %v", r)
		}
	}()
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	mu.RLock()
	defer mu.RUnlock()
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
