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

func TestRenderDev(t *testing.T) {
	e := echo.New()
	viteIntegrationConfigs := webtools.NewViteIntegrationConfigs("./ui-test/").SetIsDevEnvironment(true)
	err := webtools.ApplyViteIntegration(e, viteIntegrationConfigs)
	assert.Nil(t, err)

	exec.Command("rm", "-rf", "./ui-test/dist").Run()
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

func TestRenderDevRoute(t *testing.T) {
	e := echo.New()
	viteIntegrationConfigs := webtools.NewViteIntegrationConfigs("./ui-test/").SetIsDevEnvironment(true)
	err := webtools.ApplyViteIntegration(e, viteIntegrationConfigs)
	assert.Nil(t, err)

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
	resp, err := page.Goto("http://localhost:3000/blog")
	assert.Nil(t, err)
	if resp.Status() != http.StatusOK {
		t.Fatalf("expected status code %d, received: %d", http.StatusOK, resp.Status())
	}

	/* Replicate reactivity by clicking */
	heading := page.GetByText("My Blog", playwright.PageGetByTextOptions{Exact: playwright.Bool(true)})
	assertion := playwright.NewPlaywrightAssertions()
	err = assertion.Locator(heading).ToHaveCount(1)
	assert.Nil(t, err)
}

func TestRenderProd(t *testing.T) {
	e := echo.New()
	err := webtools.ApplyViteIntegration(
		e,
		webtools.NewViteIntegrationConfigs("./ui-test/dist").SetIsDevEnvironment(false),
	)
	assert.Nil(t, err)

	/* Build vite */
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

func TestRenderProdRoute(t *testing.T) {
	e := echo.New()
	err := webtools.ApplyViteIntegration(
		e,
		webtools.NewViteIntegrationConfigs("./ui-test/dist").SetIsDevEnvironment(false),
	)
	assert.Nil(t, err)

	/* Build vite */
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
	resp, err := page.Goto("http://localhost:3000/blog")
	assert.Nil(t, err)
	if resp.Status() != http.StatusOK {
		t.Fatalf("expected status code %d, received: %d", http.StatusOK, resp.Status())
	}

	/* Replicate reactivity by clicking */
	heading := page.GetByText("My Blog", playwright.PageGetByTextOptions{Exact: playwright.Bool(true)})
	assertion := playwright.NewPlaywrightAssertions()
	err = assertion.Locator(heading).ToHaveCount(1)
	assert.Nil(t, err)
}
