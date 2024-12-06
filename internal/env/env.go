package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/moonwalker/comet/internal/secrets"
)

var (
	defenv = ".env"
	usrenv = ".env.local"
)

func Load() {
	// .env (default)
	godotenv.Load(defenv)

	// .env.local # local user specific (usually git ignored)
	godotenv.Overload(usrenv)

	// iterate over all env vars and replace with secrets if found
	for _, e := range os.Environ() {
		k, v, ok := parseEnvLine(e)
		if ok {
			decodeEnvSecret(k, v)
		}
	}
}

func parseEnvLine(e string) (string, string, bool) {
	parts := strings.Split(e, "=")
	if len(parts) == 2 {
		k, v := parts[0], parts[1]
		if len(k) > 0 && len(v) > 0 {
			return k, v, true
		}
	}
	return "", "", false
}

func decodeEnvSecret(k, v string) {
	s, err := secrets.Get(v)
	if len(s) > 0 && err == nil {
		os.Setenv(k, s)
	}
}
