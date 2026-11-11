package artisan

import (
	"os"
	"regexp"
	"strings"
)

const defaultAppName = "inertia-go-solid-template"

var unsafeName = regexp.MustCompile(`[^a-z0-9_-]+`)

// AppName returns the service/binary name taken from APP_NAME (.env or the
// process environment), reduced to a lowercase slug such as "my-shop". The
// executable file name is what htop, ps and systemd show, so each project on
// a shared server stays identifiable.
func AppName() string {
	return slug(os.Getenv("APP_NAME"))
}

func slug(value string) string {
	value = unsafeName.ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-")
	value = strings.Trim(value, "-_")
	if value == "" {
		return defaultAppName
	}
	return value
}
