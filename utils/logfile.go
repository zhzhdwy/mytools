package utils

import (
	"log"
	"os"
)

func MyLoger(filepath string, in chan string, print bool) error {
	err := MakeDirForFile(filepath)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		return err
	}
	go func(file *os.File) {
		for {
			l := <-in
			if print {
				log.Print(l)
			}
			file.WriteString(l)
		}
		defer file.Close()
	}(file)
	return nil
}
