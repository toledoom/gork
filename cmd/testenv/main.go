package main

import (
	"fmt"
	"reflect"

	"github.com/toledoom/gork/internal/domain/battle"
	battlestorage "github.com/toledoom/gork/internal/storage/battle"
)

type Handepora struct{}

func main() {
	h := battle.Battle{}
	bs := battlestorage.DynamoStorage{}
	t := reflect.TypeOf(h)
	t2 := reflect.TypeOf(bs)
	fmt.Println(t.String())
	fmt.Println(t2.String())
}
