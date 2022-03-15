package main

import "fmt"

import "github.com/google/uuid"

func main() {
  store := NewStore("foo")

  out := &Template{
    Id: uuid.NewString(),
    Test: "round trip",
  }

  store.Write(out)

  in := &Template{
    Id: out.Id,
  }

  store.Read(in)

  fmt.Printf("%+v\n", in)
}
