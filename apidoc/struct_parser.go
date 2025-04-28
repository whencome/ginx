package apidoc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// StructParser 结构体解析器
type StructParser struct {
	conf *Config
}

func NewStructParser(c *Config) *StructParser {
	return StructParser{
		conf: c,
	}
}

// 获取结构体注释
func (p StructParser) getStructComment(structType reflect.Type) string {
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
func (p StructParser) getFieldComment(structType reflect.Type, fieldName string) string {
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

// Parse 解析结构体信息
func (p StructParser) Parse(v interface{}) StructInfo {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
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

		// 字段信息
		fileInf := FieldInfo{
			Name: field.Name,
			Type: field.Type.String(),
			Tag:  tagValue,
			Desc: comment,
		}

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			childStruct := p.Parse(field.Type)
			childStruct.Name = field.Name
			childStruct.Desc = comment
			fileInf.Struct = childStruct
		}
		structInf.Fields = append(structInf.Fields, fileInf)
	}

	return structInf
}
