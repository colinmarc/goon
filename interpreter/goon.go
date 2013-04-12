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
        fmt.Printf("err: %s\n", err)
      }
      break
    }

    if len(raw_line) > 1 {
      res, err := goon.Parse(raw_line)

      if err != nil {
        fmt.Printf("err: %s\n", err)
      } else {
        fmt.Printf("%s\n", res.ToString())
      }
    }
  }
}