// +build dev

package migration

import "net/http"

var Assets http.FileSystem = http.Dir("./scripts")
