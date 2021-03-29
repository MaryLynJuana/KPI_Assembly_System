package archive_bin

import (
	"os/exec"
	"path"
)

func Archive() {
	module_path := path.Join("stuff", "bin", "the name from build.bood")
	exec.Command("go build -o ./l " + module_path)
}
