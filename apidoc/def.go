package apidoc

import (
	"fmt"
	"sort"
)

// 版本
const Version = "1.0"

// Config 配置信息
type Config struct {
	// Title, default `API Doc`
	Title string
	// Version, default `1.0.0`
	Version string
	// Description
	Description string
	// Custom url prefix, default `/docs/api`
	UrlPrefix string
	// No document text, default `No documentation found for this API`
	NoDocText string

	// 是否启用文档
	EnableDoc bool `json:"enable_doc"`
	// 解析的字段标签名称，默认json
	FieldTag string `json:"field_tag"`

	// SHA256 encrypted authorization password, e.g. here is admin
	// echo -n admin | shasum -a 256
	// `8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918`
	PasswordSha2 string
}

// StructInfo 结构体信息
type StructInfo struct {
	Name   string      `json:"name"`   // 结构体名称
	Desc   string      `json:"desc"`   //  结构体描述
	Fields []FieldInfo `json:"fields"` // 字段信息
}

// FieldInfo 字段信息
type FieldInfo struct {
	Name     string     `json:"name"`     // 字段名称
	Required bool       `json:"required"` // 是否必填
	Type     string     `json:"type"`     // 字段类型
	Desc     string     `json:"desc"`     // 字段描述
	Tag      string     `json:"tag"`      // 字段标签
	Struct   StructInfo `json:"fields"`   // 如果是结构体，则包含字段信息
}

// 方法信息结构
type MethodInfo struct {
	Name     string   // 方法名
	Receiver string   // 接收者类型
	Comment  string   // 方法注释
	Params   []string // 参数列表
	Returns  []string // 返回值列表
}

// ApiReqParam api请求参数
type ApiReqParam struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

// ApiDocInfo 接口方法信息
type ApiDocInfo struct {
	Hash        string        `json:"hash"`        // 接口hash值，用于防重
	Name        string        `json:"name"`        // 接口方法名称
	Description string        `json:"description"` // 接口描
	MIME        string        `json:"mime"`        // 响应的MIME类型
	Path        string        `json:"path"`        // 接口路径
	Method      string        `json:"method"`      // 请求方法，POST,GET,等
	Group       string        `json:"group"`       // 文档分组
	Params      []ApiReqParam `params`             // 请求参数列表
	ParamMD     string        `json:"param_md"`    // 请求参数, markdown内容
	DocMd       string        `json:"content_md"`  // 接口文档扩展内容，markdown内容
}

func (doc *ApiDocInfo) ApiMap() KVMap {
	return KVMap{
		"api_type":    "api",
		"doc":         "",
		"description": doc.Description,
		"param_md":    doc.ParamMD,
		"mime":        doc.MIME,
		"doc_md":      doc.DocMd,
		"method":      doc.Method,
		"name":        doc.Name,
		"name_extra":  "",
		"router":      doc.Group,
		"url":         fmt.Sprintf("%s\t[%s]", doc.Path, doc.Method),
	}
}

// DocGroup 文档分组
type DocGroup struct {
	Name        string        `json:"name"`        // 分组名称
	Description string        `json:"description"` // 分组描述
	Sort        int           `json:"sort"`        // 用于控制文档排序
	Groups      []*DocGroup   `json:"groups"`      // 子分组
	Docs        []*ApiDocInfo `json:"docs"`        // 文档列表
}

// ToApiData 将文档分组转换为api数据
func (dg *DocGroup) ToApiData() DataMap {
	docMap := make(DataMap)
	// 默认分组
	defaulGroup := &DocGroup{
		Name:   "默认",
		Sort:   1,
		Docs:   make([]*ApiDocInfo, 0),
		Groups: make([]*DocGroup, 0),
	}
	if len(dg.Docs) > 0 {
		for _, doc := range dg.Docs {
			if doc.Group == "" {
				doc.Group = "默认"
			}
		}
		defaulGroup.Docs = append(defaulGroup.Docs, dg.Docs...)
	}
	docGroups := make([]*DocGroup, 0)
	docGroups = append(docGroups, defaulGroup)
	docGroups = append(docGroups, dg.Groups...)
	// 对分组进行排序
	sort.Slice(docGroups, func(i, j int) bool {
		return docGroups[i].Sort < docGroups[j].Sort
	})
	// 生成需要的结构体数据
	for _, docGroup := range docGroups {
		docMap[docGroup.Name] = docGroup.GetDocMaps()
	}
	return docMap
}

func (dg *DocGroup) GetDocMaps() RouterMap {
	routerMap := make(RouterMap)
	routerMap["children"] = make([]KVMap, 0)
	// 添加文档
	for _, doc := range dg.Docs {
		routerMap["children"] = append(routerMap["children"], doc.ApiMap())
	}
	return routerMap
}
