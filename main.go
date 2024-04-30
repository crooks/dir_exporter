package main

import (
   "fmt"
   "os"
   "path/filepath"
   "time"
   "log"
)

// timeTrack can be used to time the processing duration of a function.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func dirSize(dirPath string) (int64, error) {
  defer timeTrack(time.Now(), dirPath)
  var totalSize int64
  err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    if !info.IsDir() {
      totalSize += info.Size()
    }
    return err
  })
  return totalSize, err
}

func dirsInfo(dirs []string) {
  for _, d := range dirs {
    du, err := dirSize(d)
    if err != nil {
      panic(err)
    }
    fmt.Println(du)
  }
}

func main() {
  dirs := []string{"/opt", "/home"}
  dirsInfo(dirs)
}
