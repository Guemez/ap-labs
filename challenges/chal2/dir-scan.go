package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// scanDir stands for the directory scanning implementation
func scanDir(dir string) error {
	
	scan_result := make(map[string]int)	

    	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		
		switch mode := info.Mode(); {
		case mode.IsRegular():
			scan_result["Other"] = scan_result["Other"] + 1
			//fmt.Println("regular file")
		case mode.IsDir():
			scan_result["Directory"] = scan_result["Directory"] + 1
			//fmt.Println("directory")
		case mode&os.ModeSymlink != 0:
			scan_result["Symbolic link"] = scan_result["Symbolic link"] + 1
			//fmt.Println("symbolic link")
		case mode&os.ModeNamedPipe != 0:
			scan_result["Other"] = scan_result["Other"] + 1
			//fmt.Println("pipe")
		}
		if err != nil{
			scan_result["Other"] = scan_result["Other"] + 1
		}

        	return nil
    	})
    	if err != nil {
        	panic(err)
   	}
	for k, v :=range scan_result{
		fmt.Println(k, v)
	}
	//fmt.Println(scan_result)
	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dir-scan <directory>")
		os.Exit(1)
	}

	scanDir(os.Args[1])
}


