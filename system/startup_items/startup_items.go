package startup_items

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type StartupItem struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Args     string `json:"args"`
	Type     string `json:"type"`
	Source   string `json:"source"`
	Status   string `json:"status"`
	Username string `json:"username"`
}

not done yet