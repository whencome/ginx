package apidoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"
)

// DocParser 文档解析器
type DocParser struct {
	conf *Config
}

func NewDocParser(c *Config) *DocParser {
	return &DocParser{
		conf: c,
	}
}

// 解析文档
func (p *DocParser) Parse(r, f interface{}) ApiDocInfo {
	// 请求
	reqStruct := p.ParseStruct(r)
	// 方法
	methodInfo := p.GetMethodInfo(f)
	// 构建文档
	return p.buildDoc(reqStruct, methodInfo)
}

// 获取方法信息
func (p *DocParser) GetMethodInfo(method interface{}) MethodInfo {
	// 获取方法的反射值
	methodValue := reflect.ValueOf(method)
	if methodValue.Kind() != reflect.Func {
		return MethodInfo{}
	}

	// 获取方法指针
	methodPtr := runtime.FuncForPC(methodValue.Pointer())
	if methodPtr == nil {
		return MethodInfo{}
	}

	// 解析方法名称
	fullName := methodPtr.Name()
	parts := strings.Split(fullName, ".")
	methodName := parts[len(parts)-1]

	// 确定接收者类型
	var receiver string
	if len(parts) > 2 {
		// 方法有接收者
		receiver = parts[len(parts)-2]
		// 去掉前面的(*)
		receiver = strings.TrimPrefix(receiver, "(")
		receiver = strings.TrimPrefix(receiver, "*")
		receiver = strings.TrimSuffix(receiver, ")")
	}

	// 获取方法注释
	comment := p.getMethodComment(methodPtr)

	// 获取参数和返回值类型
	paramTypes, returnTypes := p.getMethodSignature(methodValue.Type())

	return MethodInfo{
		Name:     methodName,
		Receiver: receiver,
		Comment:  comment,
		Params:   paramTypes,
		Returns:  returnTypes,
	}
}

// 获取方法签名信息
func (p *DocParser) getMethodSignature(methodType reflect.Type) (params, returns []string) {
	// 获取参数类型
	for i := 0; i < methodType.NumIn(); i++ {
		params = append(params, methodType.In(i).String())
	}

	// 获取返回值类型
	for i := 0; i < methodType.NumOut(); i++ {
		returns = append(returns, methodType.Out(i).String())
	}

	return params, returns
}

// 获取方法注释
func (p *DocParser) getMethodComment(method *runtime.Func) string {
	filePath, _ := method.FileLine(0)
	if filePath == "" {
		return ""
	}

	// 解析源文件
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return ""
	}

	// 查找方法对应的AST节点
	methodName := method.Name()
	parts := strings.Split(methodName, ".")
	shortName := parts[len(parts)-1]

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == shortName {
				if funcDecl.Doc != nil {
					// 提取注释文本
					var lines []string
					for _, comment := range funcDecl.Doc.List {
						lines = append(lines, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
					}
					return strings.Join(lines, "\n")
				}
				return ""
			}
		}
	}

	return ""
}

// 获取结构体注释
func (p *DocParser) getStructComment(structType reflect.Type) string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	if err != nil {
		return ""
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if typeSpec.Name.Name == structType.Name() {
								if typeSpec.Doc != nil {
									return strings.TrimSpace(typeSpec.Doc.Text())
								}
								return ""
							}
						}
					}
				}
			}
		}
	}
	return ""
}

// 获取字段注释
func (p *DocParser) getFieldComment(structType reflect.Type, fieldName string) string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	if err != nil {
		return ""
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if typeSpec.Name.Name == structType.Name() {
								if structType, ok := typeSpec.Type.(*ast.StructType); ok {
									for _, field := range structType.Fields.List {
										for _, name := range field.Names {
											if name.Name == fieldName {
												if field.Doc != nil {
													return strings.TrimSpace(field.Doc.Text())
												}
												return ""
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return ""
}

// ParseStruct 解析结构体信息
func (p *DocParser) ParseStruct(v interface{}) StructInfo {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return StructInfo{}
	}

	// 结构体信息
	structInf := StructInfo{
		Name:   t.Name(),              // 结构体名称
		Desc:   p.getStructComment(t), //  结构体描述
		Fields: make([]FieldInfo, 0),  // 字段信息
	}

	// 解析字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get(p.conf.FieldTag)
		if tagValue == "" {
			tagValue = field.Name
		}
		comment := p.getFieldComment(t, field.Name)
		// 是否必填
		required := false
		binding := field.Tag.Get("binding")
		if strings.Contains(binding, "required") {
			required = true
		}

		// 字段信息
		fileInf := FieldInfo{
			Name:     field.Name,
			Required: required,
			Type:     field.Type.String(),
			Tag:      tagValue,
			Desc:     comment,
		}

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			childStruct := p.ParseStruct(field.Type)
			childStruct.Name = field.Name
			childStruct.Desc = comment
			fileInf.Struct = childStruct
		}
		structInf.Fields = append(structInf.Fields, fileInf)
	}

	return structInf
}

// buildDoc 构建文档
func (p *DocParser) buildDoc(req StructInfo, method MethodInfo) ApiDocInfo {
	apiDoc := ApiDocInfo{
		Name:    "", // 接口方法名称，这里是注解名称，用于展示，对应注解@Summary
		Path:    "", // 接口路径，对应注解@Router
		Method:  "", // 请求方法，POST,GET,等，对应注解@Router
		Content: "", // 接口文档内容
	}
	// 文档内容
	content := bytes.Buffer{}

	// 先解析请求
	reqDoc := "### 请求参数"
	if req.Name != "" {
		reqDoc += `
|参数名|必选|类型|说明|
|:----|:----|:----|----|`
		for _, field := range req.Fields {
			reqDoc += `
            ` + fmt.Sprintf("|%s|%s|%s|%s|\n", field.Name, field.Required, field.Type, field.Desc)
		}
	} else {
		reqDoc += `
            ` + "- 无"
	}

	// 解析接口文档
	funcComment := strings.TrimSpace(method.Comment)
	lines := strings.Split(funcComment, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@Summary") {
			apiDoc.Name = strings.TrimSpace(strings.TrimPrefix(line, "@Summary"))
		} else if strings.HasPrefix(line, "@Router") {
			router := strings.TrimSpace(strings.TrimPrefix(line, "@Router"))
			mStart := strings.Index(router, "[")
			mEnd := strings.Index(router, "]")
			var path, methods string
			if mStart > 0 {
				path = strings.TrimSpace(router[:mStart])
				methods = strings.TrimSpace(router[mStart+1 : mEnd])
			} else {
				path = strings.TrimSpace(router)
			}
			apiDoc.Path = path
			apiDoc.Method = methods
		} else if strings.HasPrefix(line, "@Tags") {
            apiDoc.Group = strings.TrimSpace(strings.TrimPrefix(line, "@Tags"))
		} else if strings.HasPrefix(line, "@Tags") {
        }
	}

	return apiDoc
}
