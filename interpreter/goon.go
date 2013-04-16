package main

import (
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "bufio"
  "goon"
)

func main() {
  if len(os.Args) > 1 {
    file(os.Args[1])
  } else {
    repl()
  }
}

func repl() {
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
      val := runtime.Interperet(string(raw_line))

      if val != nil {
        fmt.Printf("%s\n", val)
      }
    }
  }
}

func file(filename string) {
  input, err := ioutil.ReadFile(filename)
  if err != nil {
    fmt.Printf("error reading file: %s\n", err)
    return
  }

  runtime := goon.New()
  runtime.Interperet(string(input))
}