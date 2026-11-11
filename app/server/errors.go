package server

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

const stackKey = "app.error.stack"

// recoverMiddleware turns panics into errors for ErrorHandler and remembers the
// stack so the development page can show where it happened.
func recoverMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, _ any) {
			c.Locals(stackKey, string(debug.Stack()))
		},
	})
}

// errorHandler renders every unhandled error. Production responses are generic
// and never expose messages or stacks; every other environment gets a detailed
// page (or JSON for API clients). Replace it with your own fiber.Config.ErrorHandler
// to customize the output.
func errorHandler(environment string) fiber.ErrorHandler {
	detailed := environment != "production"
	return func(c fiber.Ctx, err error) error {
		status := http.StatusInternalServerError
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
		}
		title := fmt.Sprintf("%d %s", status, http.StatusText(status))
		message, stack := "", ""
		if detailed {
			message = err.Error()
			stack, _ = c.Locals(stackKey).(string)
		}
		c.Status(status)
		if strings.Contains(c.Get("Accept"), "application/json") && !strings.Contains(c.Get("Accept"), "text/html") {
			body := fiber.Map{"error": http.StatusText(status)}
			if detailed {
				body["message"], body["stack"] = message, stack
			}
			return c.JSON(body)
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		if !detailed {
			return c.SendString(fmt.Sprintf(errorShell, html.EscapeString(title), `<main><h1>`+html.EscapeString(title)+`</h1><p>Something went wrong on our side. Please try again later.</p></main>`))
		}
		details := fmt.Sprintf(`<main><p class="tag">Development only</p><h1>%s</h1><p class="msg">%s</p><p class="req">%s %s</p>`,
			html.EscapeString(title), html.EscapeString(message), html.EscapeString(c.Method()), html.EscapeString(c.OriginalURL()))
		if stack != "" {
			details += `<h2>Stack trace</h2><pre>` + html.EscapeString(stack) + `</pre>`
		}
		return c.SendString(fmt.Sprintf(errorShell, html.EscapeString(title), details+`<p class="hint">Production shows a generic error page instead of this one.</p></main>`))
	}
}

const errorShell = `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>%s</title>
<style>body{margin:0;font:16px/1.6 system-ui,sans-serif;background:#f5f6f2;color:#1d302c}main{max-width:60rem;margin:3rem auto;padding:0 1.5rem}h1{margin:.2rem 0;font-size:1.8rem}h2{font-size:1rem;margin-top:2rem}
.tag{display:inline-block;margin:0;padding:.1rem .6rem;border-radius:99px;background:#a52a25;color:#fff;font-size:.75rem;font-weight:700;text-transform:uppercase}.msg{font-size:1.15rem;font-weight:600;color:#a52a25;word-break:break-word}.req,.hint{color:#53605a;font-family:ui-monospace,monospace;font-size:.85rem}
pre{overflow:auto;padding:1rem;border-radius:8px;background:#1d302c;color:#e8eee7;font-size:.8rem;line-height:1.5}</style></head><body>%s</body></html>`
