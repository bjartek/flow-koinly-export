package flowgraph

import (
	"fmt"
)

type RawArgument struct {
	Value interface{}
	Type  string
}

func (self RawArgument) GetValue(_ int) string {
	return fmt.Sprintf("%v", self.Value)

}
