package main

import (
	"fmt"
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
    
}

func isRegular(path string)(bool, error){
	info,err := os.Stat(path)
	if err != nil {
		return false,err
	}

	return info.Mode().IsRegular() , nil
}