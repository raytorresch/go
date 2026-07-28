package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func main() {
	// PATHFILE
	fmt.Println("current runtime.GOOS:", runtime.GOOS)
	fmt.Println("filepath.Separator:", string(filepath.Separator))
	fmt.Println()

	fmt.Println("path.Join (always '/', ignore SO):    ", path.Join("config", "settings.json"))
	fmt.Println("filepath.Join (use SO separator):  ", filepath.Join("config", "settings.json"))

	fmt.Println()
	rutaMezclada := filepath.Join("config", "sub") + "/" + "settings.json"
	fmt.Println("mixed path (filepath.Join + '/' hardcodeado):", rutaMezclada)
	fmt.Println("-> Linux mix is not detected")
	fmt.Println("   Windows this is a problem: config\\sub/settings.json -- mixed")

	// PERMISSIONS
	f, _ := os.Create("/tmp/test_permisos.txt")
	f.Close()

	// Unix permissions style: owner/group/other, read/write/execute
	os.Chmod("/tmp/test_permisos.txt", 0640) // rw-r-----

	// Windows permissions problem
	info, _ := os.Stat("/tmp/test_permisos.txt")
	fmt.Println("Permissions:", info.Mode().Perm())
	fmt.Printf("Current permissions: %o\n", info.Mode().Perm())

	os.Remove("/tmp/test_permisos.txt")
}
