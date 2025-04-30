# ginx

## change logs

* 2025.04.30 增加接口文档支持

## 接口文档支持说明

* 本项目接口文档借鉴了[https://github.com/kwkwc/gin-docs](https://github.com/kwkwc/gin-docs)项目，对后端代码进行了部分重写以实现和ginx项目的完美兼容

**文档注释示例**

```go
// 显示当前时间
// @Summary Show Time
// @Description 显示当前时间
// @Produce json
// @Param timezone string false 时区
// @Response TimeResponse
// @Markdown
// ### 返回内容
//
// ```json
//
//	{
//	    "message": "2005-01-02"
//	}
//
// ```
//
// **返回值说明**
//
// | 字段 | 类型 | 说明 |
// | --- | --- | --- |
// | message | string | 当前时间 |
//
// @Markdown
// @Router	/time [get]
func TimeLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
	return "2005-01-02", nil
}
```

在上面的示例中，各注解的说明如下：

* **@Summary** 接口的简单说明，可以理解为接口的名称
* **@Description** 接口的文本说明，可以添加较为详细的介绍，此内容为纯文本信息
* **@Produce** 响应的内容类型，如json、xml等，最终将转换为MIME类型
* **@Param** 定义请求参数信息，格式为：@Param 字段名 类型 是否必填 参数说明
* **@Response** 响应内容，这里是注册的结构体名称，应当在注册路由之前调用apidoc.RegisterStructs或apidoc.RegisterStruct进行注册，具体参考示例example/api_example，在此示例中，结构体注册放在response包中的init方法中执行。需要说明的是，这里的结构体名称不是定义的名称，而是注册时指定的字符串名称
* **Request** 请求的结构体信息，在注册路由之前需要先注册结构体，与@Response相同
* **@Markdown** 此标签成对出现，中间的内容视为markdown文档附加到接口文档的末尾
* **@Router** 注册的路由以及请求方式

**关于请求参数的特别说明**

* 默认情况下，在调用ginx.NewApiHandler()方法时，会自动解析请求参数信息，如果不指定@Param或者@Request，则使用此处解析的结构体作为请求参数文档，否则使用@Param或者@Request处的定义，其中，@Param优先级高于@Request

**关于字段说明**

* apidoc会默认尝试解析字段的注释，如果失败，则解析字段的“desc”标签中的内容作为字段说明；
* 显示字段暂时为写死内容，如果时请求参数，则解析“form”tag内容作为显示字段，其它则取"json" tag内容作为显示字段
