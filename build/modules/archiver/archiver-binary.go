package archiver

import (
	"os/exec"
)

func Archive() {
	exec.Command("go build -o ./l ./${module_path}")
}
