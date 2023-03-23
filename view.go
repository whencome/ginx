package ginx

import (
    "html/template"
    "net/http"
    "path/filepath"
    "strings"
)

// ViewOption option for view
type ViewOption func(*View)

func WithTplDir(d string) ViewOption {
    return func(view *View) {
        view.tplDir = d
    }
}

func WithTplFiles(f ...string) ViewOption {
    return func(view *View) {
        view.tplFiles = append(view.tplFiles, f...)
    }
}

func WithTplExtension(ext string) ViewOption {
    return func(view *View) {
        view.tplExtension = ext
    }
}

type View struct {
    // tplDir 注册模板文件路径
    tplDir string // "view"
    // tplFiles 注册通用模板文件列表
    tplFiles []string
    // tplExtension 定义模板文件扩展名
    tplExtension string // ".html"
    // funcMaps 定义自定义方法列表
    funcMaps template.FuncMap
}

// NewView create a new view
func NewView(options ...ViewOption) *View {
    view := &View{
        tplDir:       "view",
        tplFiles:     make([]string, 0),
        tplExtension: ".html",
        funcMaps:     template.FuncMap{},
    }
    if len(options) > 0 {
        for _, o := range options {
            o(view)
        }
    }
    return view
}

// SetTplDir 设置模板目录
func (view *View) SetTplDir(d string) {
    if d != "" {
        view.tplDir = d
    }
}

// ContainsString 检查字符串列表是否包含某个值
func (view *View) ContainsString(arr []string, s string) bool {
    for _, v := range arr {
        if v == s {
            return true
        }
    }
    return false
}

// AddTplFiles 添加通用模板文件
func (view *View) AddTplFiles(files ...string) {
    if len(files) <= 0 {
        return
    }
    for _, f := range files {
        if view.ContainsString(view.tplFiles, f) {
            continue
        }
        view.tplFiles = append(view.tplFiles, f)
    }
}

// ResetTplFiles 情况模板文件列表
func (view *View) ResetTplFiles() {
    view.tplFiles = make([]string, 0)
}

// SetTplExtension 设置模板文件扩展名
func (view *View) SetTplExtension(ext string) {
    view.tplExtension = ext
}

// SetFuncMap 设置自定义方法列表
func (view *View) SetFuncMap(m template.FuncMap) {
    view.funcMaps = m
}

// calcTplFiles 计算需要加载的全部模板文件
func (view *View) calcTplFiles(tpl string) []string {
    tmpTplFiles := make([]string, 0)
    tplFile := filepath.Join(view.tplDir, tpl)
    if !strings.HasSuffix(tplFile, view.tplExtension) {
        tplFile += view.tplExtension
    }
    tmpTplFiles = append(tmpTplFiles, tplFile)
    for _, tplFile := range view.tplFiles {
        tplFile = filepath.Join(view.tplDir, tplFile)
        if !strings.HasSuffix(tplFile, view.tplExtension) {
            tplFile += view.tplExtension
        }
        tmpTplFiles = append(tmpTplFiles, tplFile)
    }
    return tmpTplFiles
}

// renderHtml 渲染文件
func (view *View) renderHtml(w http.ResponseWriter, name string, files []string, v interface{}) error {
    // 设置header
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // 解析模板文件
    t := template.New(name)
    if len(view.funcMaps) > 0 {
        t.Funcs(view.funcMaps)
    }
    t, e := t.ParseFiles(files...)
    if e != nil {
        logger.Errorf("parse template files failed: %s", e)
        return e
    }
    // 输出内容
    e = t.Execute(w, v)
    if e != nil {
        logger.Errorf("template execute failed: %s", e)
        return e
    }
    return nil
}

// Render 渲染文件
func (view *View) Render(w http.ResponseWriter, f string, v interface{}) error {
    tmpTplFiles := view.calcTplFiles(f)
    // 渲染html
    return view.renderHtml(w, f, tmpTplFiles, v)
}

// RenderDirect 直接渲染指定的文件
func (view *View) RenderDirect(w http.ResponseWriter, name string, files []string, v interface{}) error {
    tmpTplFiles := make([]string, 0)
    for _, tplFile := range files {
        tplFile = filepath.Join(view.tplDir, tplFile)
        if !strings.HasSuffix(tplFile, view.tplExtension) {
            tplFile += view.tplExtension
        }
        tmpTplFiles = append(tmpTplFiles, tplFile)
    }
    return view.renderHtml(w, name, tmpTplFiles, v)
}

// RenderPage 根据Page信息渲染页面
func (view *View) RenderPage(w http.ResponseWriter, p *Page) error {
    tmpTplFiles := view.calcTplFiles(p.Tpl)
    // 渲染html
    return view.renderHtml(w, p.Tpl, tmpTplFiles, p)
}

// Show 根据Page信息渲染页面
func (view *View) Show(w http.ResponseWriter, p *Page) {
    err := view.RenderPage(w, p)
    if err != nil {
        p.AddError(err)
        err = view.RenderPage(w, p)
        if err != nil {
            logger.Errorf("render page failed: %s", err)
        }
    }
}

// ShowDirect 根据Page信息渲染页面
func (view *View) ShowDirect(w http.ResponseWriter, p *Page) error {
    err := view.RenderDirect(w, p.Tpl, []string{p.Tpl}, p)
    if err != nil {
        logger.Errorf("render page failed: %s", err)
    }
    return err
}
