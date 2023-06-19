# Trash: A Trash Management CLI

Trash is a command line application written in Go that manages files and directories by implementing the [FreeDesktop.org Trash specification](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html). It allows you to move files and directories to the trash, restore them, and even permanently delete them if needed.

### Table of Contents

- [Features](#features)
- [Installation](#installation)
    - [Installing from Source](#installing-from-source)
    - [Installing Using `go install`](#installing-using-go-install)
- [Usage](#usage)
- [Examples](#examples)
- [Contribution](#contribution)
- [License](#license)

### Features

- Move files and directories to the trash.
- List contents of the trash.
- Restore files and directories from the trash.
- Permanently delete files and directories from the trash.
- Empty the trash.

### Installation

Before you start, make sure you have Go installed on your system. You can download it from [here](https://golang.org/dl/).

#### Installing from Source

1. Clone the repository:

   ```sh
   git clone https://github.com/warrenulrich/trash.git
   ```

2. Change directory to the cloned repository:

   ```sh
   cd trash
   ```

3. Build the application:

   ```sh
   go build
   ```

4. Optionally, move the built binary to a directory in your `PATH` for easier access:

   ```sh
   sudo mv trash /usr/local/bin/
   ```

#### Installing Using `go install`

If you prefer not to clone the repository and build from source, you can use `go install` to download and install GoTrash directly:

```sh
go install github.com/warrenulrich/trash@latest
```

This command downloads, builds, and installs the package in your `GOPATH`. If the `GOPATH` is in your system `PATH`, you should be able to run the `trash` command from anywhere.

### Usage

Trash can be used through the command line with the following syntax:

```
usage: trash [options] <file/directory>

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
        -v, --verbose      Output filenames while deleting them.
```

### Examples

1. Move a file to trash:

   ```sh
   trash put myfile.txt
   ```

2. Move a directory and its contents to trash:

   ```sh
   trash put -r mydirectory
   ```

3. List contents of the trash:

   ```sh
   trash list
   ```

4. Restore a file from the trash:

   ```sh
   trash restore myfile.txt
   ```

5. Permanently delete a directory from the trash:

   ```sh
   trash delete -r mydirectory
   ```

6. Empty the trash with confirmation:

   ```sh
   trash empty -c
   ```

### Contribution

Contributions to Trash are always welcome. Whether it's feature requests, bug fixes, or documentation improvements, feel free to open an issue or submit a pull request.

### License

Trash is released under the [MIT License](LICENSE).