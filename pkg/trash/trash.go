package trash

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"time"

	"gopkg.in/ini.v1"
)

type TrashInfo struct {
	DeletionDate string `ini:"DeletionDate"`
	OriginPath   string `ini:"Path"`
}

const (
	timeFormat = "2006-01-02T15:04:05"
)

var (
	TrashDirectory      string
	TrashInfoDirectory  string
	TrashFilesDirectory string

	ErrRecursiveNotSet = fmt.Errorf("recursive not set when deleting a directory")
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	TrashDirectory = usr.HomeDir + "/.local/share/Trash"
	TrashInfoDirectory = TrashDirectory + "/info"
	TrashFilesDirectory = TrashDirectory + "/files"
}

func List() ([]TrashInfo, error) {
	infos := []TrashInfo{}
	files, err := os.ReadDir(TrashInfoDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		info := TrashInfo{}
		cfg, err := ini.Load(TrashInfoDirectory + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		err = cfg.Section("Trash Info").MapTo(&info)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func Put(path string, recursive bool) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}

	if info.IsDir() && !recursive {
		return "", ErrRecursiveNotSet
	}

	name := info.Name()
	outputDir := TrashFilesDirectory + "/" + name

	err = os.Rename(path, outputDir)
	if err != nil {
		return "", err
	}

	trashInfo := TrashInfo{
		DeletionDate: time.Now().Format(timeFormat),
		OriginPath:   path,
	}

	cfg := ini.Empty()
	err = cfg.Section("Trash Info").ReflectFrom(&trashInfo)
	if err != nil {
		return "", err
	}

	err = cfg.SaveTo(TrashInfoDirectory + "/" + name + ".trashinfo")
	if err != nil {
		return "", err
	}

	return name, nil
}

func Empty() ([]string, error) {
	files, err := os.ReadDir(TrashFilesDirectory)
	removed := []string{}
	if err != nil {
		return removed, err
	}

	for _, file := range files {
		err = os.RemoveAll(TrashFilesDirectory + "/" + file.Name())
		if err != nil {
			return removed, err
		}

		removed = append(removed, file.Name())
	}

	files, err = os.ReadDir(TrashInfoDirectory)
	if err != nil {
		return removed, err
	}

	for _, file := range files {
		err = os.RemoveAll(TrashInfoDirectory + "/" + file.Name())
		if err != nil {
			return removed, err
		}
	}

	return removed, nil
}

func Restore(filename string, overwrite bool) (string, error) {
	trashInfoPath := TrashInfoDirectory + "/" + filename + ".trashinfo"
	cfg, err := ini.Load(trashInfoPath)
	if err != nil {
		return "", err
	}

	info := TrashInfo{}
	err = cfg.Section("Trash Info").MapTo(&info)
	if err != nil {
		return "", err
	}

	if !overwrite {
		_, err := os.Stat(info.OriginPath)
		if err == nil {
			return "", errors.New("file already exists in original location")
		}
	}

	trashFilePath := TrashFilesDirectory + "/" + filename
	err = os.Rename(trashFilePath, info.OriginPath)
	if err != nil {
		return "", err
	}

	err = os.Remove(trashInfoPath)
	if err != nil {
		return "", err
	}

	return info.OriginPath, nil
}

// Delete permanently deletes a file from the Trash directory.
func Delete(filename string, recursive bool) error {
	trashFilePath := TrashFilesDirectory + "/" + filename
	_, err := os.Stat(trashFilePath)
	if os.IsNotExist(err) {
		return errors.New("file not found in trash")
	} else if err != nil {
		return err
	}

	if recursive {
		err = os.RemoveAll(trashFilePath)
	} else {
		err = os.Remove(trashFilePath)
	}

	if err != nil {
		return err
	}

	trashInfoPath := TrashInfoDirectory + "/" + filename + ".trashinfo"
	err = os.Remove(trashInfoPath)
	if err != nil {
		return err
	}

	return nil
}
