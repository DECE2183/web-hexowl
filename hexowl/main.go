package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/operators"
	"github.com/dece2183/hexowl/utils"
)

func main() {
	builtin.FuncsInit()

	js.Global().Set("hexowlPrompt", js.FuncOf(prompt))

	for {
		<-make(chan bool)
	}
}

func calculate(words []utils.Word) (err error, output string) {
	operator, err := operators.Generate(words, make(map[string]interface{}))
	if err != nil {
		return
	}

	val, err := operators.Calculate(operator, make(map[string]interface{}))
	if err != nil {
		return
	}

	if val != nil {
		switch v := val.(type) {
		case string:
			output += fmt.Sprintf("\n\t%s\r\n", v)
		case bool:
			output += fmt.Sprintf("\n\tResult:\t%v\r\n", v)
		case float32, float64:
			output += fmt.Sprintf("\n\tResult:\t%f\r\n", val)
			output += fmt.Sprintf("\t\t0x%X\r\n", utils.ToNumber[uint64](val))
			output += fmt.Sprintf("\t\t0b%b\r\n", utils.ToNumber[uint64](val))
		case int64, uint64:
			output += fmt.Sprintf("\n\tResult:\t%d\r\n", val)
			output += fmt.Sprintf("\t\t0x%X\r\n", utils.ToNumber[uint64](val))
			output += fmt.Sprintf("\t\t0b%b\r\n", utils.ToNumber[uint64](val))
		case []interface{}:
			output += fmt.Sprintf("\n\tResult:\t%v\r\n", val)
			if len(v) > 0 {
				var hstr, bstr string
				switch v[0].(type) {
				case float32, float64, int64, uint64:
					for _, el := range v {
						hstr += fmt.Sprintf("0x%X ", utils.ToNumber[uint64](el))
						bstr += fmt.Sprintf("0b%b ", utils.ToNumber[uint64](el))
					}
					output += fmt.Sprintf("\t\t[%s]\r\n", hstr[:len(hstr)-1])
					output += fmt.Sprintf("\t\t[%s]\r\n", bstr[:len(bstr)-1])
				}
			}
		default:
			output += fmt.Sprintf("\n\tResult:\t%v\r\n", val)
		}
	}

	return
}

func prompt(this js.Value, inputs []js.Value) interface{} {
	var err error
	var out string

	words := utils.ParsePrompt(inputs[0].String())

	if len(words) > 0 {
		calcBeginTime := time.Now()
		err, out = calculate(words)
		calcTime := time.Since(calcBeginTime)

		if err != nil {
			out += fmt.Sprintf("\n\tError occurred: %s\n\n", err)
		} else {
			out += fmt.Sprintf("\n\tTime:\t%d ms\r\n\n", calcTime.Milliseconds())
		}
	}

	return out
}
