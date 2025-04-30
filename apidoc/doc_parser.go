package apidoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"
)

// DocParser 文档解析器
type DocParser struct{}

func NewDocParser() *DocParser {
	return &DocParser{}
}

// 解析文档
func (p *DocParser) Parse(r, f interface{}) ApiDocInfo {
	// 请求
	reqStruct := p.ParseRequest(r)
	// 方法
	methodInfo := p.ParseMethodInfo(f)
	// 构建文档
	return p.buildDoc(reqStruct, methodInfo)
}

// 获取方法信息
func (p *DocParser) ParseMethodInfo(method interface{}) MethodInfo {
	if IsNil(method) {
		return MethodInfo{}
	}
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

// ParseRequest 解析请求参数结构体信息
// 这是一个定制化的接口，用于gin通过Bind方式绑定参数的请求解析
func (p *DocParser) ParseRequest(v interface{}) RequestInfo {
	if IsNil(v) {
		return RequestInfo{}
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return RequestInfo{}
	}

	// 结构体信息
	structInf := RequestInfo{
		Name:   t.Name(),              // 结构体名称
		Desc:   p.getStructComment(t), // 结构体描述
		Fields: make([]FormField, 0),  // 字段信息
	}

	// 解析字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 显示字段名
		showFieldName := field.Tag.Get("form")
		if showFieldName == "" {
			showFieldName = field.Name
		}
		// 字段描述
		fieldDesc := field.Tag.Get("desc")
		if fieldDesc == "" {
			fieldDesc = p.getFieldComment(t, field.Name)
		}
		// 是否必填
		required := false
		binding := field.Tag.Get("binding")
		if strings.Contains(binding, "required") {
			required = true
		}

		// 字段信息
		fieldInf := FormField{
			Name:     field.Name,
			IsStruct: false,
			Required: required,
			Type:     field.Type.String(),
			Tag:      showFieldName,
			Desc:     fieldDesc,
		}

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			childStruct := p.ParseRequest(field.Type)
			childStruct.Name = field.Name
			childStruct.Desc = fieldDesc
			fieldInf.IsStruct = true
			fieldInf.Struct = childStruct
		}
		structInf.Fields = append(structInf.Fields, fieldInf)
	}

	return structInf
}

// ParseStruct 解析通用的结构体信息
func (p *DocParser) ParseStruct(v interface{}) StructInfo {
	// 如果对象为nil，则不处理
	if IsNil(v) {
		return StructInfo{}
	}

	// 获取结构体反射类型
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
		Desc:   p.getStructComment(t), // 结构体描述
		Fields: make([]FieldInfo, 0),  // 字段信息
	}
	// 解析字段信息
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 解析显示字段
		showFieldName := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			if strings.Contains(jsonTag, ",") {
				jsonTag = jsonTag[:strings.Index(jsonTag, ",")]
			}
			showFieldName = jsonTag
		}
		// 解析注释说明，应当放在desc标签中
		descTag := field.Tag.Get("desc")

		// 字段信息
		fieldInf := FieldInfo{
			Name:     field.Name,
			Tag:      showFieldName,
			IsStruct: false,
			Type:     field.Type.String(),
			Desc:     descTag,
		}

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			childStruct := p.ParseStruct(field.Type)
			childStruct.Name = field.Name
			childStruct.Desc = descTag
			fieldInf.IsStruct = true
			fieldInf.Struct = childStruct
		}
		structInf.Fields = append(structInf.Fields, fieldInf)
	}

	return structInf
}

// buildDoc 构建文档
func (p *DocParser) buildDoc(req RequestInfo, method MethodInfo) ApiDocInfo {
	apiDoc := ApiDocInfo{}
	// 解析接口文档
	funcComment := strings.TrimSpace(method.Comment)
	lines := strings.Split(funcComment, "\n")
	openMarkdown := false
	// 标识是否单独定义了参数，如果是，则不解析结构体
	definedParam := false
	// 响应结果
	respStructName := ""
	markdown := bytes.Buffer{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if openMarkdown {
			if strings.HasPrefix(line, "@Markdown") {
				openMarkdown = false
				continue
			} else {
				markdown.WriteString(line)
				markdown.WriteString("\n")
				continue
			}
		} else {
			if strings.HasPrefix(line, "@Markdown") {
				openMarkdown = !openMarkdown
				continue
			}
			if strings.HasPrefix(line, "@Summary") {
				apiDoc.Name = strings.TrimSpace(strings.TrimPrefix(line, "@Summary"))
				continue
			}
			if strings.HasPrefix(line, "@Description") {
				apiDoc.Description = strings.TrimSpace(strings.TrimPrefix(line, "@Description"))
				continue
			}
			if strings.HasPrefix(line, "@Router") {
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
				continue
			}
			if strings.HasPrefix(line, "@Tags") {
				apiDoc.Group = strings.TrimSpace(strings.TrimPrefix(line, "@Tags"))
				continue
			}
			if strings.HasPrefix(line, "@Produce") {
				produce := strings.TrimSpace(strings.TrimPrefix(line, "@Produce"))
				apiDoc.MIME = GetMIMEType(produce)
				continue
			}
			if strings.HasPrefix(line, "@Param") {
				definedParam = true
				reqParam, ok := p.parseParam(strings.TrimSpace(strings.TrimPrefix(line, "@Param")))
				if ok {
					apiDoc.Params = append(apiDoc.Params, reqParam)
				}
				continue
			}
			if strings.HasPrefix(line, "@Request") {
				if req.Name != "" { // 优先使用注册路由时使用的结构体
					continue
				}
				structName := strings.TrimSpace(strings.TrimPrefix(line, "@Request"))
				if structName == "" {
					continue
				}
				if structVal, ok := registeredTypes[structName]; ok {
					req = p.ParseRequest(structVal)
				}
				continue
			}
			if strings.HasPrefix(line, "@Response") {
				respStructName = strings.TrimSpace(strings.TrimPrefix(line, "@Response"))
				continue
			}
		}
	}
	if definedParam {
		apiDoc.ParamMD = p.buildParamMDByParams(apiDoc.Params)
	} else {
		apiDoc.ParamMD = p.buildParamMDByStruct(req)
	}

	if respStructName != "" {
		apiDoc.RespMD = p.buildRespMDByStruct(respStructName)
	}

	// 解析接口内容为html
	apiDoc.DocMd = markdown.String()

	return apiDoc
}

// parseParam 解析参数
func (p *DocParser) parseParam(param string) (ApiReqParam, bool) {
	reqParam := ApiReqParam{}
	if strings.ToLower(param) == "none" {
		return reqParam, false
	}
	chars := []rune(param)
	pos := 0
	openQuote := false
	data := make([]rune, 0)
	writeData := false
	for _, char := range chars {
		if char == '"' {
			openQuote = !openQuote
			writeData = true
		} else if char == ' ' && !openQuote {
			writeData = true
		} else {
			data = append(data, char)
		}
		if writeData {
			writeData = false
			switch pos {
			case 0:
				reqParam.Name = string(data)
			case 1:
				reqParam.Type = string(data)
			case 2:
				reqParam.Required = strings.ToLower(string(data)) == "true"
			case 3:
				reqParam.Description = string(data)
			}
			data = make([]rune, 0)
			pos++
		}
	}
	if len(data) > 0 {
		switch pos {
		case 0:
			reqParam.Name = string(data)
		case 1:
			reqParam.Type = string(data)
		case 2:
			reqParam.Required = strings.ToLower(string(data)) == "true"
		case 3:
			reqParam.Description = string(data)
		}
	}
	return reqParam, true
}

// buildParamMDByParams 根据@Param定义的参数解析请求参数markdown内容
func (p *DocParser) buildParamMDByParams(params []ApiReqParam) string {
	reqParamMD := `
|参数名|必选|类型|说明|
|:----|:----|:----|----|
`
	for _, param := range params {
		reqParamMD += fmt.Sprintf("|%s|%v|%s|%s|\n", param.Name, param.Required, param.Type, param.Description)
	}
	return reqParamMD
}

// buildParamMDByStruct 根据注册路由时使用的结构体或者通过@Request定义的结构体解析markdown内容
func (p *DocParser) buildParamMDByStruct(req RequestInfo) string {
	reqParamMD := ""
	if req.Name != "" {
		reqParamMD += `
|参数名|必选|类型|说明|
|:----|:----|:----|----|
`
		for _, field := range req.Fields {
			reqParamMD += fmt.Sprintf("|%s|%v|%s|%s|\n", field.Tag, field.Required, field.Type, field.Desc)
		}
	}
	return reqParamMD
}

// buildRespMDByStruct 根据注册的结构体（structName）解析响应结果内容
func (p *DocParser) buildRespMDByStruct(structName string) string {
	// 注册名称为空或者未找到结构体，不予解析
	if structName == "" {
		return ""
	}
	structVal, ok := registeredTypes[structName]
	if !ok {
		return ""
	}
	// 解析响应结构体对象
	obj := p.ParseStruct(structVal)
	// 解析为markdown内容
	respMD := ""
	if obj.Name == "" {
		return respMD
	}
	respMD += `
|参数名|类型|说明|
|:----|:----|----|
`
	respMD += p.buildStructMD(obj, "")

	// 添加相应结果示例
	jsonDemo, err := json.MarshalIndent(structVal, "", "    ")
	if err == nil {
		respMD += fmt.Sprintf("\n\n**示例**\n\n```json\n%s\n```\n", string(jsonDemo))
	}

	return respMD
}

// buildStructMD 构造结构体markdown内容，主要用于响应结果解析
func (p *DocParser) buildStructMD(obj StructInfo, fieldPrefix string) string {
	reqParamMD := ""
	for _, field := range obj.Fields {
		reqParamMD += fmt.Sprintf("|%s|%s|%s|\n", fieldPrefix+field.Tag, field.Type, field.Desc)
		if field.IsStruct {
			reqParamMD += p.buildStructMD(field.Struct, fieldPrefix+field.Tag+".")
		}
	}
	return reqParamMD
}
