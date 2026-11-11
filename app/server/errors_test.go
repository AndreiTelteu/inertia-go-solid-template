package server

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func errorApp(environment string) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler(environment)})
	app.Use(recoverMiddleware())
	app.Get("/panic", func(fiber.Ctx) error { panic("secret-boom") })
	app.Get("/fail", func(fiber.Ctx) error { return errors.New("secret-failure") })
	return app
}

func fetch(t *testing.T, app *fiber.App, path, accept string) (int, string) {
	t.Helper()
	req := httptest.NewRequest("GET", path, nil)
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(body)
}

func TestProductionErrorsAreGeneric(t *testing.T) {
	app := errorApp("production")
	for _, path := range []string{"/panic", "/fail"} {
		status, body := fetch(t, app, path, "text/html")
		if status != 500 || !strings.Contains(body, "500 Internal Server Error") || strings.Contains(body, "secret") || strings.Contains(body, "goroutine") {
			t.Fatalf("%s leaked or wrong status %d: %s", path, status, body)
		}
	}
	if _, body := fetch(t, app, "/fail", "application/json"); strings.Contains(body, "secret") {
		t.Fatalf("json leaked: %s", body)
	}
	if status, _ := fetch(t, app, "/missing", ""); status != 404 {
		t.Fatalf("status %d", status)
	}
}

func TestDevelopmentErrorsShowDetails(t *testing.T) {
	app := errorApp("development")
	status, body := fetch(t, app, "/panic", "text/html")
	if status != 500 || !strings.Contains(body, "secret-boom") || !strings.Contains(body, "Stack trace") || !strings.Contains(body, "goroutine") {
		t.Fatalf("missing details (%d): %s", status, body)
	}
	if _, body = fetch(t, app, "/fail", "application/json"); !strings.Contains(body, "secret-failure") {
		t.Fatalf("json lacks message: %s", body)
	}
}
