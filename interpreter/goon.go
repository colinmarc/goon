package main

import (
  "fmt"
  "os"
  "io"
  "bufio"
  "goon"
)

func main() {
  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Printf(">> ")
    raw_line, err := reader.ReadBytes('\n')

    if err != nil {
      if err == io.EOF {
        fmt.Printf("quitting...\n")
        break
      } else {
        fmt.Printf("err: %s", err)
      }
      break
    }

    if len(raw_line) > 1 {
      res, _ := goon.Parse(raw_line)
      fmt.Printf("%s\n", res.ToString())
    }
  }
}