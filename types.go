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

var NIL = &Value{nil, TYPE_NIL}
var TRUE = &Value{true, TYPE_BOOL}
var FALSE = &Value{false, TYPE_BOOL}

func (v *Value) ToString() string {
  switch v.val_type {
  case TYPE_NIL:
    return "nil"
  case TYPE_INT:
    return strconv.Itoa(v.val.(int));
  case TYPE_BOOL:
    if (v.val.(bool)) {
      return "true"
    } else {
      return "false"
    }
  }

  return fmt.Sprintf("Unknown %d: %s", v.val_type, v.val);
}

func (v *Value) IsTruthy() bool {
  if (v.val_type == TYPE_NIL) {
    return false
  }

  if (v.val_type == TYPE_BOOL && v.val.(bool) == false) {
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
    return v
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
  if v.val_type == TYPE_INT && other.val_type == TYPE_INT {
    return &Value{v.val.(int) + other.val.(int), TYPE_INT};
  }

  return nil
}

func (v *Value) Subtract(other *Value) *Value {
  if v.val_type == TYPE_INT && other.val_type == TYPE_INT {
    return &Value{v.val.(int) - other.val.(int), TYPE_INT};
  }

  return nil
}

func (v *Value) Multiply(other *Value) *Value {
  if v.val_type == TYPE_INT && other.val_type == TYPE_INT {
    return &Value{v.val.(int) * other.val.(int), TYPE_INT};
  }

  return nil
}

func (v *Value) Divide(other *Value) *Value {
  if v.val_type == TYPE_INT && other.val_type == TYPE_INT {
    return &Value{v.val.(int) / other.val.(int), TYPE_INT};
  }

  return nil
}
