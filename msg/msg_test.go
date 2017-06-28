package msg

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/grant-nelson/Gozer/common"
)

func TestKind(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(Error.String(), "Error")
	t.CheckStr(Warning.String(), "Warning")
	t.CheckStr(Info.String(), "Info")
	t.CheckStr(Debug.String(), "Debug")
	t.CheckStr((MessageKind)(12345).String(), "Unknown")
}

func TestMessage(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(NewError("Example ", "test", " error", " message").String(),
		`Error: Example test error message`)
	t.CheckStr(NewWarning("Example ", "test", " warning", " message").String(),
		`Warning: Example test warning message`)
	t.CheckStr(NewInfo("Example ", "test", " info", " message").String(),
		`Info: Example test info message`)
	t.CheckStr(NewDebug("Example ", "test", " debug", " message").String(),
		`Debug: Example test debug message`)
	t.CheckStr((*Message)(nil).String(),
		`nil`)

	msg := NewError("Example error message with data").
		Add("first", "numb").
		Add("second", "coming undone").
		Add("third", "digger")
	t.CheckStr(msg.String(),
		`Error: Example error message with data:`,
		`  first:  numb`,
		`  second: coming undone`,
		`  third:  digger`)
}

func TestDataSetter(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*DataSetter)(nil).String(), `nil`)
	d1 := NewDataSetter("prodigy", "last")
	t.CheckStr(d1.String(), `DataSetter(prodigy: last)`)
	d2 := NewDataSetter("pendulum", "rest")
	t.CheckStr(d2.String(), `DataSetter(pendulum: rest)`)
	t.CheckStr(d1.Process(d2.Process(NewError("caravan"))).String(),
		`Error: caravan:`,
		`  pendulum: rest`,
		`  prodigy:  last`)
	t.CheckStr(d2.Process(d1.Process(NewDebug("rush"))).String(),
		`Debug: rush:`,
		`  pendulum: rest`,
		`  prodigy:  last`)
	t.CheckStr(d1.Process(NewWarning("mighty")).String(),
		`Warning: mighty:`,
		`  prodigy: last`)
}

func TestFilter(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*Filter)(nil).String(), `nil`)
	f1 := NewFilter(Error)
	t.CheckStr(f1.String(), `Filter(Error)`)
	f2 := NewFilter(Debug)
	t.CheckStr(f2.String(), `Filter(Debug)`)
	t.CheckStr(f1.Process(NewDebug("official")).String(),
		`Debug: official`)
	t.CheckStr(f1.Process(NewError("still")).String(),
		`nil`)
	t.CheckStr(f1.Process(NewWarning("sail")).String(),
		`Warning: sail`)
	t.CheckStr(f2.Process(NewDebug("official")).String(),
		`nil`)
	t.CheckStr(f2.Process(NewError("still")).String(),
		`Error: still`)
	t.CheckStr(f2.Process(NewWarning("sail")).String(),
		`Warning: sail`)
}

func TestCounter(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*Counter)(nil).String(), `nil`)
	c1 := NewCounter(Error)
	c2 := NewCounter(Debug)
	t.CheckStr(c1.String(), `Error: 0`)
	t.CheckStr(c2.String(), `Debug: 0`)

	t.CheckStr(c1.Process(NewDebug("official")).String(),
		`Debug: official`)
	t.CheckStr(c1.Process(NewError("still")).String(),
		`Error: still`)
	t.CheckStr(c1.Process(NewWarning("sail")).String(),
		`Warning: sail`)
	t.CheckStr(c1.String(), `Error: 1`)
	t.CheckStr(c2.String(), `Debug: 0`)

	t.CheckStr(c2.Process(NewDebug("official")).String(),
		`Debug: official`)
	t.CheckStr(c2.Process(NewError("still")).String(),
		`Error: still`)
	t.CheckStr(c2.Process(NewWarning("sail")).String(),
		`Warning: sail`)
	t.CheckStr(c1.String(), `Error: 1`)
	t.CheckStr(c2.String(), `Debug: 1`)
}

func TestLogIO(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*LogIO)(nil).String(), `nil`)
	buf := &bytes.Buffer{}
	log := NewLogIO(buf)
	t.CheckStr(log.String(), `LogIO(Errors, Warnings)`)

	log.Process(NewError("performing"))
	log.Process(NewDebug("filtered"))
	t.CheckStr(buf.String(),
		`Error: performing`,
		``) // From the newline at the end of the log

	log.Debug = true
	log.Info = true
	t.CheckStr(log.String(), `LogIO(Errors, Warnings, Info, Debug)`)
	log.Process(NewDebug("not filtered"))
	t.CheckStr(buf.String(),
		`Error: performing`,
		`Debug: not filtered`,
		``) // From the newline at the end of the log

	log.Process(NewInfo("one"))
	log.Process(NewWarning("clover"))
	t.CheckStr(buf.String(),
		`Error: performing`,
		`Debug: not filtered`,
		`Info: one`,
		`Warning: clover`,
		``) // From the newline at the end of the log
}

func TestLogIOPrint(tt *testing.T) {
	t := common.NewTester(tt)
	result := ""
	tempPrintln := println
	println = func(msg ...interface{}) (int, error) {
		result = fmt.Sprint(msg...)
		return len(result), nil
	}
	defer func() {
		println = tempPrintln
	}()
	log := NewLogIO(nil)

	log.Process(NewError("right"))
	t.CheckStr(result,
		`Error: right`)
}

func TestLogger(tt *testing.T) {
	t := common.NewTester(tt)
	log := NewLogger()
	log.PushData("rock", "classical")
	t.CheckInt(log.ErrorCount(), 0, "initial error count")

	log.Error("tea")
	log.Debug("coffee").Add("stack", "override")
	t.CheckStr(log.String(),
		`Error: tea:`,
		`  rock: classical`,
		`Debug: coffee:`,
		`  rock:  classical`,
		`  stack: override`)
	t.CheckInt(log.ErrorCount(), 1, "first error count")

	log.Pop()
	log.Warning("break")
	t.CheckStr(log.String(),
		`Error: tea:`,
		`  rock: classical`,
		`Debug: coffee:`,
		`  rock:  classical`,
		`  stack: override`,
		`Warning: break`)
	log.Clear()
	t.CheckInt(log.ErrorCount(), 0, "cleared error count")

	log.Push(NewFilter(Info))
	log.Error("salt")
	log.Info("pepper")
	log.Warning("crazy")
	log.Info("truth")
	t.CheckStr(log.String(),
		`Error: salt`,
		`Warning: crazy`)
	t.CheckInt(log.ErrorCount(), 1, "second error count")

	log.Process(NewFilter(Warning))
	log.Process(NewDataSetter("danger", "nuts"))
	t.CheckStr(log.String(),
		`Error: salt:`,
		`  danger: nuts`)
	t.CheckInt(log.ErrorCount(), 1, "third error count")

	log.Pop()
	log.Pop() // Does nothing
	log.Pop() // Does nothing
	log.Info("pepper")
	log.Info("truth")
	log.Add(nil)
	log.Process(NewFilter(Error))
	t.CheckStr(log.String(),
		`Info: pepper`,
		`Info: truth`)
	t.CheckInt(log.ErrorCount(), 0, "filtered error count")

	(*Logger)(nil).Error("no effect")
	t.CheckInt((*Logger)(nil).ErrorCount(), 0, "error count on nil")
	t.CheckStr((*Logger)(nil).String(), ``)
}
