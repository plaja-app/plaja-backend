package controllers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (c *BaseController) GetImage(w http.ResponseWriter, r *http.Request) {
	basePath := "storage"

	filePath := chi.URLParam(r, "*")

	fullPath := filepath.Join(basePath, filePath)

	if !strings.HasPrefix(fullPath, filepath.Clean(basePath)+string(os.PathSeparator)) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(fullPath); err == nil {
		http.ServeFile(w, r, fullPath)
	} else {
		http.NotFound(w, r)
	}
}
