package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"os"
	
)

func main() {
	path := "filepath.txt"
	rg,err := isRegular(path)
	if err != nil {
		fmt.Println(err)
	} 
	if (rg == true){
		
		fmt.Printf("file checked")
	}
	err = CreateApplicationsDir()
	if err != nil {
		
	fmt.Println(err)
	}
	
	fmt.Println("Applications directory created")
	return

}

//is the the file a regular or normal file
func isRegular(path string)(bool, error){
	info,err := os.Stat(path)
	if err != nil {
		return false,err
	}

	return info.Mode().IsRegular() , nil
}


func CreateApplicationsDir()error{
	home,err := os.UserHomeDir()
	if err != nil {
		return err
	}

	Newpath := filepath.Join(home,"Applications")
    err =	os.MkdirAll(Newpath,0755)
	if err != nil {
		return err
	}
	
	return nil 
}

func MoveAppImage(path,path2 string)error{
	err := os.Rename(path,path2)
	if err != nil {
		return err
	}
	return nil
}

func ExtractAppImage(path string) error {
    cmd := exec.Command(path, "--appimage-extract")

    err := cmd.Run()
    if err != nil {
        return err
    }

    return nil
}

func FindDesktopFile(root string) (string, error) {
	var found string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".desktop" {
			found = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return found, nil
}

func FindIcon(root string) (string, error) {
	var found string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if ext == ".png" || ext == ".svg" || ext == ".ico" {
			found = path
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if found == "" {
		return "", fmt.Errorf("no icon found")
	}

	return found, nil
}

