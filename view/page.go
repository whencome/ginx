package view

import (
    "fmt"
    "github.com/whencome/ginx/types"
    "runtime/debug"

    "github.com/gin-gonic/contrib/sessions"
    "github.com/gin-gonic/gin"
)

// 定义全局页面初始化方法变量，用于在每次创建Page时进行初始化
var initPageFunc PageInitFunc = nil

// PageInitFunc 定义一个页面初始化方法，返回map[string]string，返回的数据将放到页面的会话数据（Page.Session）中
type PageInitFunc func() map[string]interface{}

func RegisterPageInitFunc(f PageInitFunc) {
    initPageFunc = f
}

// PageError 定义一个页面错误, 用于保存错误以及堆栈信息
type PageError struct {
    Message string
    Trace   string
}

func NewPageError(e error) *PageError {
    return &PageError{
        Message: e.Error(),
        Trace:   string(debug.Stack()),
    }
}

func (pe *PageError) Error() string {
    return pe.Message
}

// Page 定义一个页面数据
type Page struct {
    Ctx     *gin.Context
    Request types.Request          // 页面请求数据
    Tpl     string                 `json:"tpl"`    // 定义模板
    Title   string                 `json:"title"`  // 页面标题
    Data    map[string]interface{} `json:"data"`   // 页面数据，可能会向用户展示
    Sess    map[string]interface{} `json:"Sess"`   // 保存会话数据，用于服务端业务处理，不对用户展示
    Errors  []*PageError           `json:"errors"` // 错误列表
}

// NewPage 创建一个Page对象
func NewPage(c *gin.Context, tpl string) *Page {
    p := &Page{
        Ctx:  c,
        Tpl:  tpl,
        Data: make(map[string]interface{}),
        Sess: make(map[string]interface{}),
    }
    return p
}

// NewPageWithData 创建一个Page对象，并初始化数据
func NewPageWithData(c *gin.Context, tpl string, data map[string]interface{}) *Page {
    p := &Page{
        Ctx:    c,
        Tpl:    tpl,
        Data:   data,
        Sess:   make(map[string]interface{}),
        Errors: make([]*PageError, 0),
    }
    return p
}

// init 页面初始化
func (p *Page) init() {
    if initPageFunc == nil {
        return
    }
    sess := initPageFunc()
    if sess == nil || len(sess) == 0 {
        return
    }
    for k, v := range sess {
        p.Sess[k] = v
    }
}

// Session 获取会话信息
func (p *Page) Session() sessions.Session {
    return sessions.Default(p.Ctx)
}

// SetTitle 设置页面标题
func (p *Page) SetTitle(t string) {
    p.Title = t
}

// SetData 设置页面数据
func (p *Page) SetData(d map[string]interface{}) {
    p.Data = d
}

// AddData 添加数据
func (p *Page) AddData(k string, d interface{}) {
    if p.Data == nil {
        p.Data = make(map[string]interface{})
    }
    p.Data[k] = d
}

// BatchAddData 批量添加数据
func (p *Page) BatchAddData(d map[string]interface{}) {
    if p.Data == nil {
        p.Data = make(map[string]interface{})
    }
    for k, v := range d {
        p.Data[k] = v
    }
}

// AddError 添加错误信息
func (p *Page) AddError(e error) {
    if p.Errors == nil {
        p.Errors = make([]*PageError, 0)
    }
    p.Errors = append(p.Errors, NewPageError(e))
}

// HasError 判断页面是否有错误
func (p *Page) HasError() bool {
    if len(p.Errors) > 0 {
        return true
    }
    return false
}

// Show 显示页面内容
func (p *Page) Show() {
    ShowPage(p.Ctx.Writer, p)
}

// ShowWithError 将错误添加进页面并显示
func (p *Page) ShowWithError(e interface{}) {
    p.AddError(fmt.Errorf("%s", e))
    p.Show()
}

// ShowDirect 显示页面内容
func (p *Page) ShowDirect() {
    ShowPageDirect(p.Ctx.Writer, p)
}

// ShowDirectWithError 显示页面内容
func (p *Page) ShowDirectWithError(e interface{}) {
    p.AddError(fmt.Errorf("%s", e))
    p.ShowDirect()
}
