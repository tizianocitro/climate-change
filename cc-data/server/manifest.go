// This file is automatically generated. Do not modify it manually.

package main

import (
	"encoding/json"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "climate-change-data",
  "name": "Climate Change Data",
  "description": "Climate Change Data is a plugin for a climate change collaboration platform.",
  "homepage_url": "https://github.com/tizianocitro/climate-change/cc-data/",
  "support_url": "https://github.com/tizianocitro/climate-change/cc-data/issues",
  "release_notes_url": "https://github.com/tizianocitro/climate-change/cc-data/releases/tag/",
  "icon_path": "assets/plugin_icon.svg",
  "version": "+",
  "min_server_version": "7.6.0",
  "server": {
    "executables": {
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "darwin-arm64": "server/dist/plugin-darwin-arm64",
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "linux-arm64": "server/dist/plugin-linux-arm64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "",
    "footer": "",
    "settings": []
  }
}
`

func init() {
	_ = json.NewDecoder(strings.NewReader(manifestStr)).Decode(&manifest)
}
