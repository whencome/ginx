package apidoc

import "strings"

// 全局配置
var config *Config

// 文档解析器
var docParser *DocParser

// 是否启用文档
var docEnabled bool = false

// 保存全局文档信息
var apiDocs = &DocGroup{
	Name:        "ginx接口文档",             // 分组名称
	Description: "",                     // 分组描述
	Sort:        100,                    // 用于控制文档排序
	Groups:      make([]*DocGroup, 0),   // 子分组
	Docs:        make([]*ApiDocInfo, 0), // 文档列表
}

// DefaultConfig 生成一个默认配置
func DefaultConfig() *Config {
	return &Config{
		EnableDoc: false, // 默认不启用文档
		FieldTag:  "form",
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
	if !docEnabled { // 文档未启用
		return
	}
	apiDoc := docParser.Parse(r, f)
	// 忽略标题为空的文档
	if apiDoc.Name == "" {
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
}
