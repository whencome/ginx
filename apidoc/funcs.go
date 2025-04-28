package apidoc

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// MIMEMaps 定义mime映射关系
var MIMEMaps = map[string]string{
	"text": "text/plain",                                                              // 纯文本（.txt）
	"bin":  "application/octet-stream",                                                // 二进制数据
	"html": "text/html",                                                               // HTML文档（.html）
	"css":  "text/css",                                                                // 层叠样式表（.css）
	"csv":  "text/csv",                                                                // 逗号分隔值（.csv）
	"json": "application/json",                                                        // JSON数据（.json）
	"xml":  "application/xml",                                                         // XML文档（.xml）
	"jpeg": "image/jpeg",                                                              // JPEG图片（.jpg）
	"png":  "image/png",                                                               // PNG图片（.png）
	"webp": "image/webp",                                                              // WebP图片（.webp）
	"svg":  "image/svg+xml",                                                           // 矢量图（.svg）
	"mpeg": "audio/mpeg",                                                              // MP3音频（.mp3）
	"mp4":  "video/mp4",                                                               // MP4视频（.mp4）
	"ogg":  "application/ogg",                                                         // OGG媒体（.ogv/.oga）
	"pdf":  "application/pdf",                                                         // PDF文档（.pdf）
	"doc":  "application/msword",                                                      // Word旧格式（.doc）
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document", // DOCX（.docx）
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",       // （.xlsx）
	"ppt":  "application/vnd.ms-powerpoint",                                           // .ppt
	"zip":  "application/zip",                                                         // ZIP压缩（.zip）
	"rar":  "application/x-rar-compressed",                                            // RAR压缩（.rar）
	"7z":   "application/x-7z-compressed",                                             // 7-Zip压缩（.7z）
}

// Markdown2Html
func Markdown2Html(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

// GetMIMEType 获取mime类型
func GetMIMEType(contentType string) string {
	if mimetype, ok := MIMEMaps[contentType]; ok {
		return mimetype
	}
	return "text/plain"
}
