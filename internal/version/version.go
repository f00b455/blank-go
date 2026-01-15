// Package version provides centralized version management for the application.
// The version can be set at build time using ldflags:
//
//	go build -ldflags "-X github.com/f00b455/blank-go/internal/version.Version=1.2.3"
package version

// Version is the current application version.
// Default is "dev" for development builds.
// Set via ldflags during release builds.
var Version = "dev"
