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
  runtime := goon.New()

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
      line := string(raw_line)
      val := runtime.Interperet(line)

      if val != nil {
        fmt.Printf("%s\n", val.ToString())
      }
    }
  }
}