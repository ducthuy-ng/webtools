package webtools_test

import (
	"embed"
	"strings"
	"testing"

	"github.com/ducthuy-ng/webtools"
)

func TestParseUsingPath(t *testing.T) {
	templates, err := webtools.ParseTemplates("./templates")
	if err != nil {
		t.Errorf("expect load templates not failed, %v:", err)
	}

	var resultBuffer strings.Builder
	err = templates.ExecuteTemplate(&resultBuffer, "index.html", nil)
	if err != nil {
		t.Errorf("expect parse templates/index.html not failed, %v:", err)
	}

	parsedIndexHtml := resultBuffer.String()
	if parsedIndexHtml != "<h1>Index</h1>" {
		t.Errorf("parsed index.html does not correct, received: %s", parsedIndexHtml)
	}
}

//go:embed templates/*.html
var templatesFS embed.FS

func TestParseFSInEmbed(t *testing.T) {
	templates, err := webtools.ParseFS(templatesFS, "templates")
	if err != nil {
		t.Errorf("expect load templates not failed, %v:", err)
	}

	var resultBuffer strings.Builder
	err = templates.ExecuteTemplate(&resultBuffer, "index.html", nil)
	if err != nil {
		t.Errorf("expect parse templates/index.html not failed, %v:", err)
	}

	parsedIndexHtml := resultBuffer.String()
	if parsedIndexHtml != "<h1>Index</h1>" {
		t.Errorf("parsed templates/index.html does not correct, received: %s", parsedIndexHtml)
	}
}
