package api

import (
	"fmt"
	"net/http"
)

func JSONError(w http.ResponseWriter, err string, code int) {
	http.Error(w, fmt.Sprintf(`{"error":%q}`, err), code)
}
