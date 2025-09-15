# Wispo Backup Service
Implementation of an scalable, Modoular and light-weight service to maintain and schedule backups. 

# Initial Setup

```go
go mod init github.com/wispodeveloper-lgtm/wispo
go get github.com/robfig/cron/v3
go run .
```

# Data Types 

There are 2 data types to ensure the object-orientation architecture, which is inherently not supported in golang. To cleaning and making our code maintainable I used the following data types : 

**File**: To store data which related to the file and keep track of scheduling and retention policy of the file . 

> **Sync** : I used this flag to check if the all backuped data sync with the schedule and not behind it. 

**Storage**: To store the list of files and to ease of tracking files and storage volume and free space. 

> **Lock** : to prevent concurrently changes on a single file while backup process, I defined lock to turn it on/off to prevent data curroption and dirty flags. 

# Interface 

There are 4 Major commands, users can use each of these commands to add, remove, print or sync the changes in backup folder. 

```go
type Backup interface {
	AddFile(sourcePath, destionationPath string, compressionMethod string, scheduling string) error
	RemoveFile(destinationPath string) error
	PrintFiles()
	Sync() error
}
```

## AddFile

command to add a new file to the backup storage, each file contains the following inputs : 

**sourcePath** : The absolute path of the file<br/>
**destionationPath**: By default is `/backup/<file-name>.<compression-format>`<br/> 
**compressionMethod**: Compression method (gzip, zlip, etc.)<br/> 
**scheduling**: Schedule policy in the cron format : `* * * * *`<br/> 
**retainPolicy**: Retaining Policy in the following format : `1m` or `1h` or `1d`<br/>

### Compression 

Currently the code uses two methods for compressioning, Gzip and zlib. Both are Implemented in `compression` package and can be accessed in main code. You can make your own compression methods for specific data like theses two implemented method. 

```go
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
```

each case uses an specific method to compress data, and it's easily scalable.

## RemoveFile 

Useful when we want to stop backing up an specific file/folder. to free up space or security issues. 

```go
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
```

The function removes the file from the storage files to keep the storage and real data consistent. 


## PrintFiles   

```go
func (s *Storage) PrintFiles() {
	for i := range s.Files {
		fmt.Printf("%s : Synced: %t \n", s.Files[i].destionationPath, s.Files[i].synced)
	}
}
```

Printing files and showing the sync status of each. 

## Sync

```go
func (s *Storage) Sync() error {
	return nil
}
```
The syncing process which I didnt had enough time to complete it. 


# Main code 

In the main code I used a channel called `chan` to keep track of user inputs while the kubernetes operator is already running. 

Also I've defined a Ticker to keep track of time and check the scheduling of the backup process for each file. 

```go
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
```

In the bottom of the code there is a while loop to run endlessly to process both user commands and also ticker pulses. 

# CronJob

To process crons and reduce redundancy in the code I used an third-party package to handle crons and also I've modified it to maintain some other features that are required in this project. 

```go
type CronJob struct {
	cron *cron.Cron
}

func NewCron() *CronJob {
	c := CronJob{
		cron: cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow))),
	}
	return &c
}

func (c *CronJob) AddNewJob(schedule string) {
	_, err := c.cron.AddFunc(schedule, func() {
		fmt.Println("Job executed at:", time.Now().Format("15:04:05"))
	})
	if err != nil {
		panic(err)
	}
}
```


# Kubernetes Operatores

I didn't implement this part of project. I used chatGPT. 


Thanks for your attention.
