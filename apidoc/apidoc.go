package apidoc

// 全局配置
var config *Config

// 结构体解析器
var structParser *StructParser

// 文档解析器
var docParser *DocParser

// 是否启用文档
var docEnabled bool = false

// 保存全局文档信息
var apiDocs []*DocGroup

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
	docParser = NewDocParser(c)
	structParser = NewStructParser(c)
	apiDocs = make([]*DocGroup, 0)
}

// AddDoc 添加文档
func AddDoc(doc *ApiDocInfo) {
    
}
