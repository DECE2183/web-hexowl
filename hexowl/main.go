package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/input/syntax"
	"github.com/dece2183/hexowl/operators"
	"github.com/dece2183/hexowl/utils"
)

var virtOut bytes.Buffer
var sysDesc = builtin.System{
	Stdout: &virtOut,
	ClearScreen: func() {
		js.Global().Call("clearOutput")
	},
	ListEnvironments: func() ([]string, error) {
		return nil, nil
	},
	WriteEnvironment: envWrite,
	ReadEnvironment:  envRead,
}

func main() {
	builtin.SystemInit(sysDesc)

	js.Global().Set("hexowlPrompt", js.FuncOf(prompt))

	for {
		<-make(chan bool)
	}
}

func calculate(words []utils.Word) {
	var output string
	var err error

	var operator *operators.Operator
	var val interface{}

	calcBeginTime := time.Now()

	operator, err = operators.Generate(words, make(map[string]interface{}))
	if err != nil {
		goto result
	}

	val, err = operators.Calculate(operator, make(map[string]interface{}))
	if err != nil {
		goto result
	}

	output = virtOut.String()
	virtOut.Reset()

	if val != nil {
		var resultStr string

		switch v := val.(type) {
		case string:
			output += fmt.Sprintf("\n\t%s\r\n", v)
			goto result
		case float32, float64:
			resultStr = fmt.Sprintf(
				"\t%f\r\n\t\t0x%X\r\n\t\t0b%b\r\n",
				v,
				utils.ToNumber[uint64](val),
				utils.ToNumber[uint64](val),
			)
		case int64, uint64:
			resultStr = fmt.Sprintf(
				"\t%d\r\n\t\t0x%X\r\n\t\t0b%b\r\n",
				v,
				utils.ToNumber[uint64](val),
				utils.ToNumber[uint64](val),
			)
		case []interface{}:
			resultStr = fmt.Sprintf("\t%v\r\n", v)
			if len(v) > 0 {
				var hstr, bstr string
				switch v[0].(type) {
				case float32, float64, int64, uint64:
					for _, el := range v {
						hstr += fmt.Sprintf("0x%X ", utils.ToNumber[uint64](el))
						bstr += fmt.Sprintf("0b%b ", utils.ToNumber[uint64](el))
					}
					resultStr += fmt.Sprintf("\t\t[%s]\r\n", hstr[:len(hstr)-1])
					resultStr += fmt.Sprintf("\t\t[%s]\r\n", bstr[:len(bstr)-1])
				}
			}
		default:
			resultStr = fmt.Sprintf("\t%v\r\n", v)
		}

		output += "\n\tResult:" + syntax.Highlight(resultStr)
	}

result:
	calcTime := time.Since(calcBeginTime)

	if err != nil {
		output += fmt.Sprintf("\n\tError occurred: %s\n\n", err)
	} else {
		output += fmt.Sprintf("\n\tTime:\t%d ms\r\n\n", calcTime.Milliseconds())
	}

	js.Global().Call("printOutput", output)
}

func prompt(this js.Value, inputs []js.Value) interface{} {
	input := inputs[0].String()
	js.Global().Call("printOutput", ">: "+input)

	words := utils.ParsePrompt(input)

	if len(words) > 0 {
		go calculate(words)
	}

	return nil
}
