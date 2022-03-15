package main

import "fmt"

type Object interface {
  GetType() string
  GetId() string
}

func ObjectKey(obj Object) string {
  return fmt.Sprintf(
    "%d:%s:%d:%s",
    len(obj.GetType()),
    obj.GetType(),
    len(obj.GetId()),
    obj.GetId(),
  )
}
