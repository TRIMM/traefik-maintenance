// Package traefik_maintenance a maintenance page plugin.
package traefik_maintenance

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	Enabled          bool   `json:"enabled"`
	Filename         string `json:"filename"`
	TriggerFilename  string `json:"triggerFilename"`
	HttpResponseCode int    `json:"httpResponseCode"`
	HttpContentType  string `json:"httpContentType"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Enabled:          false,
		Filename:         "",
		TriggerFilename:  "",
		HttpResponseCode: http.StatusServiceUnavailable,
		HttpContentType:  "text/html; charset=utf-8",
	}
}

// MaintenancePage a maintenance page plugin.
type MaintenancePage struct {
	next             http.Handler
	enabled          bool
	filename         string
	triggerFilename  string
	httpResponseCode int
	HttpContentType  string
	name             string
	template         *template.Template
}

// New created a new MaintenancePage plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Filename) == 0 {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	return &MaintenancePage{
		enabled:          config.Enabled,
		filename:         config.Filename,
		triggerFilename:  config.TriggerFilename,
		httpResponseCode: config.HttpResponseCode,
		HttpContentType:  config.HttpContentType,
		next:             next,
		name:             name,
		template:         template.New("MaintenancePage").Delims("[[", "]]"),
	}, nil
}

func (a *MaintenancePage) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.maintenanceEnabled() {
		// Return the maintenance page
		bytes, err := os.ReadFile(a.filename)
		if err == nil {
			rw.Header().Add("Content-Type", a.HttpContentType)
			rw.WriteHeader(a.httpResponseCode)
			_, err = rw.Write(bytes)
			if err == nil {
				return
			} else {
				log.Printf("Could not serve maintenance template %s: %s", a.filename, err)
			}
		} else {
			log.Printf("Could not read maintenance template %s: %s", a.filename, err)
		}
	}

	a.next.ServeHTTP(rw, req)
}

// Indicates if maintenance mode has been enabled
func (a *MaintenancePage) maintenanceEnabled() bool {
	if !a.enabled {
		return false
	}

	if a.enabled && len(a.triggerFilename) == 0 {
		return true
	}

	// Check if the trigger exists
	if _, err := os.Stat(a.triggerFilename); err == nil {
		return true
	}

	return false
}
