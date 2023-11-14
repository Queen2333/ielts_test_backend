package utils

import (
	"bytes"
	"html/template"
)

// GenerateHTMLPageWithPIN 生成包含PIN码的HTML页面
func GenerateHTMLPageWithPIN(pinCode string, htmlTemplate string) (string, error) {
	// 解析HTML模板
	t, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	// 创建一个包含PIN码的数据结构
	data := struct {
		PINCode string
	}{
		PINCode: pinCode,
	}

	// 创建一个缓冲区来存储生成的HTML
	var buf bytes.Buffer

	// 生成HTML到缓冲区
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	// 将缓冲区的内容转换为字符串并返回
	return buf.String(), nil
}
