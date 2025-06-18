// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package repome

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/Trae-AI/stream-to-river/apiservice/resource"
)

var (
	htmlTemplate string
	cssContent   string
	jsContent    string
)

func init() {
	// Read embedded files using resources package
	htmlTemplateBytes, err := resource.GetRepoMeHTMLTemplate()
	if err != nil {
		hlog.Fatalf("Failed to read embedded HTML template file: %v", err)
	}
	htmlTemplate = string(htmlTemplateBytes)

	cssContentBytes, err := resource.GetRepoMeCSSContent()
	if err != nil {
		hlog.Fatalf("Failed to read embedded CSS file: %v", err)
	}
	cssContent = string(cssContentBytes)

	jsContentBytes, err := resource.GetRepoMeJSContent()
	if err != nil {
		hlog.Fatalf("Failed to read embedded JS file: %v", err)
	}
	jsContent = string(jsContentBytes)
}

// TemplateData holds the data for HTML template rendering
type TemplateData struct {
	CSS        template.CSS
	Navigation template.HTML
	Content    template.HTML
	JavaScript template.JS
}

func RepoMeHandler(ctx context.Context, c *app.RequestContext) {
	// Set Content-Type to text/html
	c.Header("Content-Type", "text/html; charset=utf-8")

	// Read markdown file content using resources package
	markdownContent, err := resource.GetRepoMeMarkdownContent()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read documentation file: %v", err)
		return
	}

	// Generate Feishu-style HTML with proper markdown rendering
	htmlContent := generateRepoMeDocHTML(string(markdownContent))

	// Return HTML content
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

// StaticHandler serves static files like images
func StaticHandler(ctx context.Context, c *app.RequestContext) {
	// Get the file path from URL parameter
	filename := c.Param("filename")
	if filename == "" {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// Try to read file content using resources package
	fileContent, err := resource.GetRepoMeStaticFile(filename)
	if err != nil {
		c.String(http.StatusNotFound, "File not found: %s", filename)
		return
	}

	// Set appropriate content type based on file extension
	ext := filepath.Ext(filename)
	var contentType string
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	default:
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, fileContent)
}

func generateRepoMeDocHTML(markdownContent string) string {
	// Extract navigation from markdown content
	navigation := extractNavigationItems(markdownContent)

	// Convert markdown to HTML properly
	htmlContent := properMarkdownToHTML(markdownContent)

	// Parse template
	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		return fmt.Sprintf("Template parse error: %v", err)
	}

	// Prepare template data
	data := TemplateData{
		CSS:        template.CSS(getCSS()),
		Navigation: template.HTML(navigation),
		Content:    htmlContent,
		JavaScript: template.JS(getJavaScript()),
	}

	// Execute template
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Sprintf("Template execution error: %v", err)
	}

	return result.String()
}

// getCSS returns the CSS styles for the page
func getCSS() string {
	return cssContent
}

// getJavaScript returns the JavaScript code for the page
func getJavaScript() string {
	return jsContent
}

func extractNavigationItems(content string) string {
	lines := strings.Split(content, "\n")
	var nav strings.Builder
	var inCodeBlock bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 检查代码块标记
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		// 跳过代码块内的内容
		if inCodeBlock {
			continue
		}

		// 处理标题（只有在非代码块内才处理）
		if strings.HasPrefix(line, "#") {
			level := 0
			for _, c := range line {
				if c == '#' {
					level++
				} else {
					break
				}
			}

			title := strings.TrimSpace(line[level:])
			id := generateID(title)

			levelClass := ""
			switch level {
			case 1:
				levelClass = "level-1"
			case 2:
				levelClass = "level-2"
			case 3:
				levelClass = "level-3"
			case 4:
				levelClass = "level-4"
			default:
				levelClass = "level-4"
			}

			nav.WriteString(`<div class="nav-item ` + levelClass + `" data-href="#` + id + `">`)

			if level == 1 {
				nav.WriteString(`<span class="nav-icon"> 🚀 </span>`)
			}
			nav.WriteString(`<span class="nav-text">` + html.EscapeString(title) + `</span>`)
			nav.WriteString(`</div>`)
		}
	}

	return nav.String()
}

func generateID(title string) string {
	// 转换为 URL 友好的 ID
	id := strings.ToLower(title)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, ".", "")
	id = strings.ReplaceAll(id, "/", "-")
	id = strings.ReplaceAll(id, "(", "")
	id = strings.ReplaceAll(id, ")", "")
	return id
}

// 改进的Markdown到HTML转换函数
func properMarkdownToHTML(text string) template.HTML {
	lines := strings.Split(text, "\n")
	var result strings.Builder
	var inCodeBlock bool
	var inTable bool
	var inBlockquote bool

	for i, line := range lines {
		line = strings.TrimRight(line, " \t")

		// 处理代码块
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "```") {
			if inCodeBlock {
				result.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				lang := strings.TrimPrefix(trimmedLine, "```")
				if lang == "" {
					lang = "text"
				}
				result.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">", html.EscapeString(lang)))
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			result.WriteString(html.EscapeString(line) + "\n")
			continue
		}

		// 处理表格
		if strings.Contains(line, "|") && !strings.HasPrefix(line, "#") {
			if !inTable {
				result.WriteString("<table>\n")
				inTable = true
			}

			cells := strings.Split(line, "|")
			if len(cells) > 2 { // 过滤掉空的分割结果
				cells = cells[1 : len(cells)-1] // 移除首尾空元素

				// 检查是否是表头分隔行
				if i > 0 && strings.Contains(line, "---") {
					continue
				}

				// 判断是否是表头
				isHeader := i == 0 || (i > 0 && !strings.Contains(lines[i-1], "|"))
				if i > 0 {
					// 查找前一个非空行
					for j := i - 1; j >= 0; j-- {
						if strings.TrimSpace(lines[j]) != "" {
							isHeader = !strings.Contains(lines[j], "|")
							break
						}
					}
				}

				tag := "td"
				if isHeader {
					tag = "th"
					result.WriteString("<thead>\n")
				} else if inTable && i > 0 && !strings.Contains(lines[i-1], "|") {
					result.WriteString("<tbody>\n")
				}

				result.WriteString("<tr>")
				for _, cell := range cells {
					cell = strings.TrimSpace(cell)
					result.WriteString(fmt.Sprintf("<%s>%s</%s>", tag, processInlineElements(cell), tag))
				}
				result.WriteString("</tr>\n")

				if isHeader {
					result.WriteString("</thead>\n")
				}
			}
			continue
		} else if inTable {
			result.WriteString("</tbody></table>\n")
			inTable = false
		}

		// 处理引用块
		if strings.HasPrefix(line, ">") {
			if !inBlockquote {
				result.WriteString("<blockquote>\n")
				inBlockquote = true
			}
			content := strings.TrimSpace(strings.TrimPrefix(line, ">"))
			if content != "" {
				result.WriteString(fmt.Sprintf("<p>%s</p>\n", processInlineElements(content)))
			}
			continue
		} else if inBlockquote {
			result.WriteString("</blockquote>\n")
			inBlockquote = false
		}

		// 处理标题（只有在非代码块内才处理）
		if !inCodeBlock && strings.HasPrefix(line, "#") {
			level := 0
			for _, c := range line {
				if c == '#' {
					level++
				} else {
					break
				}
			}
			if level <= 6 {
				title := strings.TrimSpace(line[level:])
				result.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", level, processInlineElements(title), level))
				continue
			}
		}

		// 处理列表
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			content := strings.TrimSpace(line[2:])
			result.WriteString(fmt.Sprintf("<ul><li>%s</li></ul>\n", processInlineElements(content)))
			continue
		}

		// 处理有序列表
		if matched, _ := regexp.MatchString(`^\d+\.\s+`, line); matched {
			re := regexp.MustCompile(`^\d+\.\s+(.*)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				content := matches[1]
				result.WriteString(fmt.Sprintf("<ol><li>%s</li></ol>\n", processInlineElements(content)))
				continue
			}
		}

		// 处理段落
		if strings.TrimSpace(line) != "" {
			result.WriteString(fmt.Sprintf("<p>%s</p>\n", processInlineElements(line)))
		} else {
			result.WriteString("<br/>\n")
		}
	}

	// 关闭未闭合的标签
	if inCodeBlock {
		result.WriteString("</code></pre>\n")
	}
	if inTable {
		result.WriteString("</tbody></table>\n")
	}
	if inBlockquote {
		result.WriteString("</blockquote>\n")
	}

	return template.HTML(result.String())
}

// 处理行内元素
func processInlineElements(text string) string {
	// 转义HTML
	text = html.EscapeString(text)

	// 处理粗体
	re := regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = re.ReplaceAllString(text, "<strong>$1</strong>")

	// 处理斜体
	re = regexp.MustCompile(`\*(.*?)\*`)
	text = re.ReplaceAllString(text, "<em>$1</em>")

	// 处理行内代码
	re = regexp.MustCompile("`([^`]+)`")
	text = re.ReplaceAllString(text, "<code>$1</code>")

	// 处理图片 - 先处理图片再处理链接，避免冲突
	re = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		imgRe := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
		matches := imgRe.FindStringSubmatch(match)

		if len(matches) >= 3 {
			alt := matches[1]
			src := matches[2]

			// 如果是本地图片文件，转换为静态服务路径
			if !strings.HasPrefix(src, "http") && !strings.HasPrefix(src, "//") {
				src = "/api/static/" + filepath.Base(src)
			}

			return fmt.Sprintf(`<img src="%s" alt="%s" style="max-width: 100%%; height: auto; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); margin: 16px 0;">`, src, alt)
		}
		return match
	})

	// 处理链接
	re = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = re.ReplaceAllString(text, `<a href="$2">$1</a>`)

	return text
}
