package cmd

import (
	"os"
	"os/exec"
	"path"

	"github.com/AlecAivazis/survey/v2"
)

const DirPerms = 0700

func InstallCodegen(dir string) error {
	prompt := &survey.Confirm{
		Message: "Do you want to install the TypeScript SDK and code generator?",
		Default: true,
	}
	var install bool
	err := survey.AskOne(prompt, &install)
	if err != nil {
		return err
	}

	if !install {
		return nil
	}

	npm, err := exec.LookPath("npm")
	if err != nil {
		return err
	}

	err = execNpm(npm, "install", "@xata.io/client")
	if err != nil {
		return err
	}

	err = execNpm(npm, "install", "@xata.io/codegen", "-D")
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(dir, "hooks"), DirPerms)
	if err != nil {
		if !os.IsExist((err)) {
			return err
		}
	}

	err = os.WriteFile(path.Join(dir, "hooks", "build"), []byte("#!/bin/bash\n./node_modules/.bin/xata-codegen generate -o client.ts"), 0744)
	if err != nil {
		return err
	}

	return nil
}

func execNpm(npm string, arg ...string) error {
	cmd := exec.Command(npm, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
