package cmd

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

func RunHook(dir string, hook string) error {
	hooks := path.Join(dir, "hooks")

	files, err := os.ReadDir(hooks)
	if err != nil {
		if os.IsNotExist(err) {
			// Don't do anything if the hooks directory doesn't exist
			return nil
		}
		return err
	}
	for i := 0; i < len(files); i++ {
		name := files[i].Name()
		if !strings.HasPrefix(name, hook) {
			continue
		}
		cmd := exec.Command(path.Join(hooks, name))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			if os.IsPermission(err) {
				// We silently ignore files that are not executable
				continue
			}
			return err
		}
	}
	return nil
}
