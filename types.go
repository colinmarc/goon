package goon

import (
  "fmt"
  "strconv"
)

type ValueType int
const (
  TYPE_NIL ValueType = -1

  _  = iota
  TYPE_INT
  TYPE_BOOL
)

type Value struct {
  val interface{}
  val_type ValueType
}

func (v *Value) ToString() string {
  switch v.val_type {
  case TYPE_INT:
    return strconv.Itoa(v.val.(int));
  case TYPE_BOOL:
    if(v.val.(bool)) {
      return "true"
    } else {
      return "false"
    }
  }

  return fmt.Sprintf("Unknown %d: %s", v.val_type, v.val);
}

func (v *Value) Add(other *Value) *Value {
  if(v.val_type == TYPE_INT && other.val_type == TYPE_INT) {
    return &Value{v.val.(int) + other.val.(int), TYPE_INT};
  }

  return nil
}

func (v *Value) Subtract(other *Value) *Value {
  if(v.val_type == TYPE_INT && other.val_type == TYPE_INT) {
    return &Value{v.val.(int) - other.val.(int), TYPE_INT};
  }

  return nil
}

func (v *Value) Multiply(other *Value) *Value {
  if(v.val_type == TYPE_INT && other.val_type == TYPE_INT) {
    return &Value{v.val.(int) * other.val.(int), TYPE_INT};
  }

  return nil
}

func (v *Value) Divide(other *Value) *Value {
  if(v.val_type == TYPE_INT && other.val_type == TYPE_INT) {
    return &Value{v.val.(int) / other.val.(int), TYPE_INT};
  }

  return nil
}
