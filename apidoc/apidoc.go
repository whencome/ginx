package apidoc

import (
	"fmt"
	"strings"
)

// 全局配置
var config *Config

// 文档解析器
var docParser *DocParser

// 是否启用文档
var docEnabled bool = true

// 保存全局文档信息
var apiDocs *DocGroup

// 初始化全局对象
func init() {
	config = DefaultConfig()
	apiDocs = &DocGroup{
		Name:        "",                     // 分组名称
		Description: "",                     // 分组描述
		Sort:        100,                    // 用于控制文档排序
		Groups:      make([]*DocGroup, 0),   // 子分组
		Docs:        make([]*ApiDocInfo, 0), // 文档列表
	}
	docParser = NewDocParser(config)
}

// DefaultConfig 生成一个默认配置
func DefaultConfig() *Config {
	return &Config{
		// Title, default `API Doc`
		Title: "Ginx文档",
		// Version, default `1.0.0`
		Version: "1.0",
		// Description
		Description: "",
		// Custom url prefix, default `/docs/api`
		UrlPrefix: "/ginx/docs",
		// No document text, default `No documentation found for this API`
		NoDocText: "<no documents>",

		// 是否启用文档
		EnableDoc: true,
		// 解析的字段标签名称，默认json
		FieldTag: "form",

		// SHA256 encrypted authorization password, e.g. here is admin
		// echo -n admin | shasum -a 256
		// `8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918`
		PasswordSha2: "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
	}
}

// Init 初始化配置
func Init(c *Config) {
	if c == nil {
		c = DefaultConfig()
	}
	// 初始化全局对象
	if !c.EnableDoc { // 不启用文档
		return
	}
	// 初始化全局对象
	config = c
	docEnabled = c.EnableDoc
	docParser = NewDocParser(c)
}

// Parse 解析文档
// r 请求结构体
// f 接口方法
func Parse(r, f interface{}) {
	fmt.Println("--- 解析文档")
	if !docEnabled { // 文档未启用
		fmt.Println("--- 解析文档: 文档未启用")
		return
	}
	apiDoc := docParser.Parse(r, f)
	// 忽略标题为空的文档
	if apiDoc.Name == "" {
		fmt.Println("--- 解析文档: 文档名称为空")
		return
	}
	// 添加文档
	AddDoc(&apiDoc)
}

// AddDoc 添加文档
func AddDoc(doc *ApiDocInfo) {
	if doc.Name == "" {
		return
	}
	groupName := strings.TrimSpace(doc.Group)
	if groupName == "" {
		apiDocs.Docs = append(apiDocs.Docs, doc)
		return
	}
	found := false
	for _, g := range apiDocs.Groups {
		if g.Name == groupName {
			g.Docs = append(g.Docs, doc)
			found = true
			break
		}
	}
	if !found {
		g := &DocGroup{
			Name:   groupName,
			Sort:   100,
			Docs:   make([]*ApiDocInfo, 0),
			Groups: make([]*DocGroup, 0),
		}
		g.Docs = append(g.Docs, doc)
		apiDocs.Groups = append(apiDocs.Groups, g)
	}

	/*
	   groupPaths := strings.Split(groupName, "/")
	   var group *DocGroup = apiDocs
	   for _, groupPath := range groupPaths {
	       if group.Groups == nil {
	           group.Groups = make([]*DocGroup, 0)
	       } else {
	           found := false
	           for _, g := range group.Groups {
	               if g.Name == groupPath {
	                   group = g
	                   found = true
	                   break
	               }
	           }
	           if !found {
	               g := &DocGroup{
	                   Name:   groupPath,
	                   Sort:   100,
	                   Docs:   make([]*ApiDocInfo, 0),
	                   Groups: make([]*DocGroup, 0),
	               }
	               group = g
	           }
	       }
	   }
	   group.Docs = append(group.Docs, doc)
	*/

}
