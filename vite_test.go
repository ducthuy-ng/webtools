package webtools_test

import (
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/ducthuy-ng/webtools"
	"github.com/labstack/echo/v4"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

func bunInstall() {
	viteCmd := exec.Command("bun", "i", "-D")
	viteCmd.Dir = "./ui-test"
	viteCmd.Run()
}

func TestRenderDev(t *testing.T) {
	e := echo.New()
	viteIntegrationConfigs := webtools.NewViteIntegrationConfigs("./ui-test/").SetIsDevEnvironment(true)
	err := webtools.ApplyViteIntegration(e, viteIntegrationConfigs)
	assert.Nil(t, err)

	bunInstall()
	viteCmd := exec.Command("bun", "run", "vite")
	viteCmd.Dir = "./ui-test"
	viteCmd.Start()

	go func() {
		e.Start(":3000")
	}()

	t.Cleanup(func() {
		viteCmd.Process.Signal(os.Kill)
		e.Close()
	})

	pw, err := playwright.Run()
	assert.Nil(t, err)
	browser, err := pw.Chromium.Launch()
	assert.Nil(t, err)
	page, err := browser.NewPage()
	assert.Nil(t, err)
	resp, err := page.Goto("http://localhost:3000/")
	assert.Nil(t, err)
	if resp.Status() != http.StatusOK {
		t.Fatalf("expected status code %d, received: %d", http.StatusOK, resp.Status())
	}

	/* Replicate reactivity by clicking */
	button := page.GetByRole("button").First()
	button.Click()
	button.Click()
	button.Click()
	content, err := button.TextContent()
	assert.Nil(t, err)
	assert.Equal(t, content, "count is 3")
}

func TestRenderProd(t *testing.T) {
	e := echo.New()

	err := webtools.ApplyViteIntegration(
		e,
		webtools.NewViteIntegrationConfigs("./ui-test/dist").SetIsDevEnvironment(false),
	)
	assert.Nil(t, err)

	/* Build vite */
	bunInstall()
	buildCmd := exec.Command("bun", "run", "build")
	buildCmd.Dir = "./ui-test/"
	buildCmd.Run()

	go func() {
		e.Start(":3000")
	}()

	defer func() {
		e.Close()
	}()

	pw, err := playwright.Run()
	assert.Nil(t, err)
	browser, err := pw.Chromium.Launch()
	assert.Nil(t, err)
	page, err := browser.NewPage()
	assert.Nil(t, err)
	resp, err := page.Goto("http://localhost:3000/")
	assert.Nil(t, err)
	if resp.Status() != http.StatusOK {
		t.Fatalf("expected status code %d, received: %d", http.StatusOK, resp.Status())
	}

	/* Replicate reactivity by clicking */
	button := page.GetByRole("button").First()
	button.Click()
	button.Click()
	button.Click()
	content, err := button.TextContent()
	assert.Nil(t, err)
	assert.Equal(t, content, "count is 3")
}
