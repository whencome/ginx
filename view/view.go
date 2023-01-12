package view

import (
    "html/template"
    "log"
    "net/http"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/contrib/sessions"
)

// tplDir 注册模板文件路径
var tplDir string = "view"

// tplFiles 注册通用模板文件列表
var tplFiles []string = make([]string, 0)

// tplExtension 定义模板文件扩展名
var tplExtension = ".html"

// funcMaps 定义自定义方法列表
var funcMaps template.FuncMap = nil

// 设置全局变量store
var store = sessions.NewCookieStore([]byte("ginx_sess"))

// Init 初始化View
// tplDir 为模板的根目录
// tplExt 为模板文件的统一后缀，包含“.”号，如“.html”
// publicTplFiles 为通用模板文件，如果没有，则置空
func Init(tplDir, tplExt string, publicTplFiles ...string) {
    SetTplDir(tplDir)
    SetTplExtension(tplExt)
    AddTplFiles(publicTplFiles...)
}

// SetTplDir 设置模板目录
func SetTplDir(d string) {
    if d != "" {
        tplDir = d
    }
}

// ContainsString 检查字符串列表是否包含某个值
func ContainsString(arr []string, s string) bool {
    for _, v := range arr {
        if v == s {
            return true
        }
    }
    return false
}

// AddTplFiles 添加通用模板文件
func AddTplFiles(files ...string) {
    if len(files) <= 0 {
        return
    }
    for _, f := range files {
        if ContainsString(tplFiles, f) {
            continue
        }
        tplFiles = append(tplFiles, f)
    }
}

// ResetTplFiles 情况模板文件列表
func ResetTplFiles() {
    tplFiles = make([]string, 0)
}

// SetTplExtension 设置模板文件扩展名
func SetTplExtension(ext string) {
    tplExtension = ext
}

// SetFuncMap 设置自定义方法列表
func SetFuncMap(m template.FuncMap) {
    funcMaps = m
}

// calcTplFiles 计算需要加载的全部模板文件
func calcTplFiles(tpl string) []string {
    tmpTplFiles := make([]string, 0)
    tplFile := filepath.Join(tplDir, tpl)
    if !strings.HasSuffix(tplFile, tplExtension) {
        tplFile += tplExtension
    }
    tmpTplFiles = append(tmpTplFiles, tplFile)
    for _, tplFile := range tplFiles {
        tplFile = filepath.Join(tplDir, tplFile)
        if !strings.HasSuffix(tplFile, tplExtension) {
            tplFile += tplExtension
        }
        tmpTplFiles = append(tmpTplFiles, tplFile)
    }
    return tmpTplFiles
}

// renderHtml 渲染文件
func renderHtml(w http.ResponseWriter, name string, files []string, v interface{}) error {
    // 设置header
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // 解析模板文件
    t, e := template.New(name).Funcs(funcMaps).ParseFiles(files...)
    if e != nil {
        return e
    }
    // 输出内容
    e = t.Execute(w, v)
    if e != nil {
        return e
    }
    return nil
}

// Render 渲染文件
func Render(w http.ResponseWriter, f string, v interface{}) error {
    tmpTplFiles := calcTplFiles(f)
    // 渲染html
    return renderHtml(w, f, tmpTplFiles, v)
}

// RenderDirect 直接渲染指定的文件
func RenderDirect(w http.ResponseWriter, name string, files []string, v interface{}) error {
    tmpTplFiles := make([]string, 0)
    for _, tplFile := range files {
        tplFile = filepath.Join(tplDir, tplFile)
        if !strings.HasSuffix(tplFile, tplExtension) {
            tplFile += tplExtension
        }
        tmpTplFiles = append(tmpTplFiles, tplFile)
    }
    return renderHtml(w, name, tmpTplFiles, v)
}

// RenderPage 根据Page信息渲染页面
func RenderPage(w http.ResponseWriter, p *Page) error {
    tmpTplFiles := calcTplFiles(p.Tpl)
    // 渲染html
    return renderHtml(w, p.Tpl, tmpTplFiles, p)
}

// ShowPage 根据Page信息渲染页面
func ShowPage(w http.ResponseWriter, p *Page) {
    err := RenderPage(w, p)
    if err != nil {
        p.AddError(err)
        err = RenderPage(w, p)
        if err != nil {
            log.Printf("render page failed: %s", err)
        }
    }
}

// ShowPageDirect 根据Page信息渲染页面
func ShowPageDirect(w http.ResponseWriter, p *Page) {
    err := RenderDirect(w, p.Tpl, []string{p.Tpl}, p)
    if err != nil {
        log.Printf("render page failed: %s", err)
    }
}
