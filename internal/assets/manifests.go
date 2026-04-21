package assets

import "embed"

// ManifestsFS embeds all Kubernetes manifests for the Claw instance
//
//go:embed manifests
var ManifestsFS embed.FS
