package main

import (
	"fmt"
	"os"
	"path/filepath"
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