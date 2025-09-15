package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wispo/compression"
	"github.com/wispo/schedule"
)

type Backup interface {
	AddFile(sourcePath, destionationPath, compressionMethod, scheduling, retainPolicy string) error
	RemoveFile(destinationPath string) error
	PrintFiles()
	Sync() error
}

type Storage struct {
	Files      []File
	FreeVolume int
	UsedVolume int
	Lock       bool
}

type File struct {
	sourcePath, destionationPath string
	compressionMethod            string
	retainPolicy                 string
	scheduling                   string
	synced                       bool
}

func NewStorage(Volume int) *Storage {
	return &Storage{
		FreeVolume: Volume,
		UsedVolume: 0,
		Lock:       false,
	}
}

func (s *Storage) AddFile(sourcePath, destionationPath, compressionMethod, scheduling, retainPolicy string) error {
	file := File{sourcePath: sourcePath,
		destionationPath:  "",
		compressionMethod: compressionMethod,
		scheduling:        scheduling,
		retainPolicy:      retainPolicy,
		synced:            false,
	}

	switch compressionMethod {
	case "gzip":
		path, err := compression.GzStore(sourcePath, destionationPath+".gz")
		if err != nil {
			log.Fatal("Cant compress the file :", err)
		}
		file.destionationPath = path
		s.Files = append(s.Files, file)

	case "zlib":
		path, err := compression.GzStore(sourcePath, destionationPath+".gz")
		if err != nil {
			log.Fatal("Cant compress the file :", err)
		}
		file.destionationPath = path
		s.Files = append(s.Files, file)
	default:
		panic("Not defined compression method")
	}

	return nil
}

func (s *Storage) RemoveFile(destinationPath string) error {
	for i := range s.Files {
		if s.Files[i].destionationPath == destinationPath {
			err := os.Remove(destinationPath)
			if err != nil {
				fmt.Println("Error deleting file:", err)
			}
			s.Files = append(s.Files[:i], s.Files[i+1:]...)
			return nil
		}
	}
	panic("Cannot Find the file")
}

func (s *Storage) PrintFiles() {
	for i := range s.Files {
		fmt.Printf("%s : Synced: %t \n", s.Files[i].destionationPath, s.Files[i].synced)
	}
}

func (s *Storage) Sync() error {
	return nil
}

func main() {
	s := NewStorage(10000)
	c := schedule.NewCron()

	s.AddFile("sample.txt", "backup/sample.txt", "gzip", "*/1 * * * *", "1y")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	cmdChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			cmdChan <- text
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}()

	for {
		select {
		case t := <-ticker.C:
			fmt.Println(t)
		case cmd := <-cmdChan:
			switch cmd {
			case "AddFile":
				c.AddNewJob("*/1 * * * *")
			case "RemoveFile":

			case "PrintFiles":
				s.PrintFiles()
			}
		}

	}

}
