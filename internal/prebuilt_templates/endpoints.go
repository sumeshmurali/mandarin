package prebuilttemplates

import (
	"net/http"

	"errors"
)


var ErrTemplateNotFound = errors.New("template not found")

func GetTemplate(name string) (http.HandlerFunc, error) {
	switch name {
	case "echo":
		return EchoTemplate, nil
	default:
		return nil, ErrTemplateNotFound
	}
}


func EchoTemplate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}