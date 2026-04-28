package ginx

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/whencome/ginx/log"
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
	// tplDir register template file path
	tplDir string // "view"
	// tplFiles register common template file list
	tplFiles []string
	// tplExtension define template file extension
	tplExtension string // ".html"
	// funcMaps define custom function list
	funcMaps template.FuncMap
	// template cache for better performance
	templateCache map[string]*template.Template
	cacheMutex    sync.RWMutex
}

// NewView create a new view
func NewView(options ...ViewOption) *View {
	view := &View{
		tplDir:        "view",
		tplFiles:      make([]string, 0),
		tplExtension:  ".html",
		funcMaps:      template.FuncMap{},
		templateCache: make(map[string]*template.Template),
	}
	if len(options) > 0 {
		for _, o := range options {
			o(view)
		}
	}
	return view
}

// SetTplDir set template directory
func (view *View) SetTplDir(d string) {
	if d != "" {
		view.tplDir = d
	}
}

// ContainsString check if string slice contains a value
func (view *View) ContainsString(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

// AddTplFiles add common template files
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

// ResetTplFiles reset template file list
func (view *View) ResetTplFiles() {
	view.tplFiles = make([]string, 0)
}

// SetTplExtension set template file extension
func (view *View) SetTplExtension(ext string) {
	view.tplExtension = ext
}

// SetFuncMap set custom function list
func (view *View) SetFuncMap(m template.FuncMap) {
	view.funcMaps = m
}

// calcTplFiles calculate all template files to load
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

// renderHtml render file with caching support
func (view *View) renderHtml(w http.ResponseWriter, name string, files []string, v interface{}) error {
	// set header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// generate cache key
	cacheKey := strings.Join(files, "|")

	// try to get from cache
	view.cacheMutex.RLock()
	if tmpl, ok := view.templateCache[cacheKey]; ok {
		view.cacheMutex.RUnlock()
		// execute cached template
		err := tmpl.Execute(w, v)
		if err != nil {
			log.Errorf("template execute failed: %s", err)
			return err
		}
		return nil
	}
	view.cacheMutex.RUnlock()

	// parse template files
	view.cacheMutex.Lock()
	defer view.cacheMutex.Unlock()

	// double check after acquiring write lock
	if tmpl, ok := view.templateCache[cacheKey]; ok {
		err := tmpl.Execute(w, v)
		if err != nil {
			log.Errorf("template execute failed: %s", err)
			return err
		}
		return nil
	}

	t := template.New(name)
	if len(view.funcMaps) > 0 {
		t.Funcs(view.funcMaps)
	}
	t, err := t.ParseFiles(files...)
	if err != nil {
		log.Errorf("parse template files failed: %s", err)
		return err
	}

	// cache the template
	view.templateCache[cacheKey] = t

	// output content
	err = t.Execute(w, v)
	if err != nil {
		log.Errorf("template execute failed: %s", err)
		return err
	}
	return nil
}

// Render render file
func (view *View) Render(w http.ResponseWriter, f string, v interface{}) error {
	tmpTplFiles := view.calcTplFiles(f)
	// render html
	return view.renderHtml(w, f, tmpTplFiles, v)
}

// RenderDirect directly render specified file
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

// RenderPage render page based on Page info
func (view *View) RenderPage(w http.ResponseWriter, p *Page) error {
	tmpTplFiles := view.calcTplFiles(p.Tpl)
	// render html
	return view.renderHtml(w, p.Tpl, tmpTplFiles, p)
}

// Show render page based on Page info
func (view *View) Show(w http.ResponseWriter, p *Page) {
	err := view.RenderPage(w, p)
	if err != nil {
		p.AddError(err)
		err = view.RenderPage(w, p)
		if err != nil {
			log.Errorf("render page failed: %s", err)
		}
	}
}

// ShowDirect render page based on Page info directly
func (view *View) ShowDirect(w http.ResponseWriter, p *Page) error {
	err := view.RenderDirect(w, p.Tpl, []string{p.Tpl}, p)
	if err != nil {
		log.Errorf("render page failed: %s", err)
	}
	return err
}
