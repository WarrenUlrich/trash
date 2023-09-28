package main

import (
	"flag"
	"fmt"

	"github.com/warrenulrich/trash/pkg/trash"
)

const (
	helpMessage string = `usage: trash [options] <file/directory>
	
Options:
	-h, --help          Show help message and exit.
	--version           Show program's version number and exit.

Commands:
	put <file/directory>   Move specified file or directory to the trash.
		-r, --recursive    Required if deleting directories. Move directories and their contents to the trash.
		-v, --verbose      Output filenames while moving them to the trash.

	list                   List all files and directories currently in the trash.

	restore <file>         Restore the specified file or directory from the trash to its original location.
		-v, --verbose      Output filenames while restoring them.
		-o, --overwrite    Overwrite an existing file in the original location with the restored file.

	delete <file>          Permanently delete the specified file or directory from the trash.
		-r, --recursive    Required if deleting directories. Permanently delete directories and their contents from the trash.
		-v, --verbose      Output filenames while deleting them.

	empty                  Permanently delete all files and directories from the trash.
		-c, --confirm      Ask for confirmation before emptying the trash.
		-v, --verbose      Output filenames while deleting them.`

	versionMessage = "trash version 0.0.1"
)

func put(args []string) error {
	flags := flag.NewFlagSet("put", flag.ExitOnError)
	verbose := flags.Bool("v", false, "Output filenames while moving them to the trash")
	recursive := flags.Bool("r", false, "Move directories and their contents to the trash")

	if err := flags.Parse(args); err != nil {
		return err
	}

	dir := flags.Arg(0)
	if dir == "" {
		return fmt.Errorf("no file or directory specified")
	}

	name, err := trash.Put(dir, *recursive)
	if err != nil {
		return err
	}

	if *verbose {
		fmt.Printf("Moved %s to trash\n", name)
	}

	return nil
}

func list(args []string) error {
	infos, err := trash.List()
	if err != nil {
		return err
	}

	for _, info := range infos {
		fmt.Printf("%s %s\n", info.DeletionDate, info.OriginPath)
	}

	return nil
}

func empty(args []string) error {
	flags := flag.NewFlagSet("empty", flag.ExitOnError)
	verbose := flags.Bool("v", false, "Output filenames while deleting them")
	confirm := flags.Bool("c", false, "Ask for confirmation before emptying the trash")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *confirm {
		fmt.Print("Are you sure you want to empty the trash? [y/N] ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return nil
		}
	}

	removed, err := trash.Empty()
	if err != nil {
		return err
	}

	if *verbose {
		for _, name := range removed {
			fmt.Printf("Deleted %s\n", name)
		}
	}

	return nil
}

func restore(args []string) error {
	flags := flag.NewFlagSet("restore", flag.ExitOnError)
	verbose := flags.Bool("v", false, "Output filenames while restoring them")
	overwrite := flags.Bool("o", false, "Overwrite an existing file in the original location with the restored file")

	if err := flags.Parse(args); err != nil {
		return err
	}

	dir := flags.Arg(0)
	if dir == "" {
		return fmt.Errorf("no file or directory specified")
	}

	name, err := trash.Restore(dir, *overwrite)
	if err != nil {
		return err
	}

	if *verbose {
		fmt.Printf("Restored %s from trash\n", name)
	}

	return nil
}

func del(args []string) error {
	flags := flag.NewFlagSet("delete", flag.ExitOnError)
	verbose := flags.Bool("v", false, "Output filenames while deleting them")
	recursive := flags.Bool("r", false, "Permanently delete directories and their contents from the trash")

	if err := flags.Parse(args); err != nil {
		return err
	}

	dir := flags.Arg(0)
	if dir == "" {
		return fmt.Errorf("no file or directory specified")
	}

	err := trash.Delete(dir, *recursive)
	if err != nil {
		return err
	}

	if *verbose {
		fmt.Printf("Deleted %s from trash\n", dir)
	}

	return nil
}

func runCommand(command string, args []string) error {
	switch command {
	case "help":
		fmt.Println(helpMessage)
	case "put":
		return put(args)
	case "list":
		return list(args)
	case "restore":
		return restore(args)
	case "delete":
		return del(args)
	case "empty":
		return empty(args)
	}

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Println(helpMessage)
	}

	help := flag.Bool("h", false, "Show help message and exit")
	version := flag.Bool("version", false, "Show program's version number and exit")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *version {
		fmt.Println(versionMessage)
		return
	}

	remainingArgs := flag.Args()

	if len(remainingArgs) < 1 {
		flag.Usage()
		return
	}

	err := runCommand(remainingArgs[0], remainingArgs[1:])
	if err != nil {
		panic(err)
	}
}
