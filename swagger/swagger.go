package swagger

import (
	_ "embed"
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

//go:embed api.swagger.json
var spec []byte

func Mux(path string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(path, swaggerui.Handler(spec))
	return mux
}
