// windows only
package shellutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var Dir []string
var DebugOpt = true

func init() {
	wd, err := os.Getwd()
	if EvalError(err) {
		os.Exit(1)
	}
	Dir = strings.Split(wd, string(filepath.Separator))

	if filepath.IsAbs(wd) {
		Dir = append([]string{string(filepath.Separator)}, Dir...)
	}
	fmt.Println(Dir)
}
func EvalError(err error) bool {
	if err != nil {
		if DebugOpt {
			fmt.Println(err)
		}
		return true
	}
	return false
}

func Shell() {
	fmt.Println("start")
	test := "pwd"
	cmd, err := ConstructCmd(test)
	if EvalError(err) {
		os.Exit(1)
	}
	err = cmd.Exec()
	if EvalError(err) {
		os.Exit(2)
	}
	test = "cd"
	cmd, err = ConstructCmd(test)
	if EvalError(err) {
		os.Exit(1)
	}
	err = cmd.Exec()
	if EvalError(err) {
		os.Exit(2)
	}
	fmt.Println("end")
}

var (
	ErrBadCmd          = errors.New("bad command")
	ErrCmdNotSupported = errors.New("command not supported")
	ErrNotEnoughArgs   = errors.New("not enough arguments")
)

type SupportedCommand string

const (
	pwd SupportedCommand = "pwd"
	cd  SupportedCommand = "cd"
)

func ConstructCmd(cmd string) (Cmd, error) {
	args := strings.Split(cmd, " ")
	if len(args) < 1 {
		return nil, ErrBadCmd
	}
	switch SupportedCommand(args[0]) {
	case pwd:
		return NewPwd(args[1:]), nil
	case cd:
		return NewPwd(args[1:]), nil
	}
	return nil, ErrCmdNotSupported
}

type Cmd interface {
	Args() []string
	Exec() error
}
type Pwd struct {
	args []string
}

func NewPwd(args []string) Pwd {
	return Pwd{
		args: args,
	}
}

// Args implements Cmd.
func (c Pwd) Args() []string {
	return c.args
}

// Exec implements Cmd.
func (c Pwd) Exec() error {
	fmt.Println(filepath.Join(Dir...))
	return nil
}

type Cd struct {
	args []string
}

func NewCd(args []string) Pwd {
	return Pwd{
		args: args,
	}
}

// Args implements Cmd.
func (c Cd) Args() []string {
	return c.args
}

// Exec implements Cmd.
func (c Cd) Exec() error {
	if len(c.args) < 1 {
		return ErrNotEnoughArgs
	}
	return nil
}
