package goon

import "fmt"

type Runtime struct {
  ns map[string]*Value
}

func New() *Runtime {
  runtime := &Runtime{}
  runtime.ns = make(map[string]*Value, 1024)

  return runtime
}

func (r *Runtime) Interperet(input string) *Value {
  root, err := Parse(input)
  if err != nil {
    fmt.Printf("Error! %s\n", err)
    return nil
  }

  //root.Describe(0)
  return root.Evaluate(r)
}