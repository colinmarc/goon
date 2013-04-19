package goon

import (
  "fmt"
  "strconv"
)

type ValueType int
const (
  NilType ValueType = -1

  _  = iota
  IntType
  BoolType
)

type Value struct {
  val interface{}
  val_type ValueType
}

var NIL = &Value{nil, NilType}
var TRUE = &Value{true, BoolType}
var FALSE = &Value{false, BoolType}

func (v *Value) String() string {
  switch v.val_type {
  case NilType:
    return "nil"
  case IntType:
    return strconv.Itoa(v.val.(int));
  case BoolType:
    if (v.val.(bool)) {
      return "true"
    } else {
      return "false"
    }
  }

  return fmt.Sprintf("Unknown %d: %s", v.val_type, v.val);
}

func (v *Value) IsTruthy() bool {
  if (v.val_type == NilType) {
    return false
  }

  if (v.val_type == BoolType && v.val.(bool) == false) {
    return false
  }

  return true
}

func (v *Value) Or(other *Value) *Value {
  if v.IsTruthy() {
    return v
  } else if (other.IsTruthy()) {
    return other
  }

  return FALSE
}

func (v *Value) And(other *Value) *Value {
  if (v.IsTruthy() && other.IsTruthy()){
    return other
  }

  return FALSE
}

func (v *Value) Equals(other *Value) *Value {
  if (v.val_type == other.val_type && v.val == other.val) {
    return TRUE
  }

  return FALSE
}

func (v *Value) NotEquals(other *Value) *Value {
  if (v.val_type == other.val_type && v.val == other.val) {
    return FALSE
  }

  return TRUE
}

func (v *Value) Add(other *Value) *Value {
  if v.val_type == IntType && other.val_type == IntType {
    return &Value{v.val.(int) + other.val.(int), IntType};
  }

  return nil
}

func (v *Value) Subtract(other *Value) *Value {
  if v.val_type == IntType && other.val_type == IntType {
    return &Value{v.val.(int) - other.val.(int), IntType};
  }

  return nil
}

func (v *Value) Multiply(other *Value) *Value {
  if v.val_type == IntType && other.val_type == IntType {
    return &Value{v.val.(int) * other.val.(int), IntType};
  }

  return nil
}

func (v *Value) Divide(other *Value) *Value {
  if v.val_type == IntType && other.val_type == IntType {
    return &Value{v.val.(int) / other.val.(int), IntType};
  }

  return nil
}
