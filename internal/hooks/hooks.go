package hooks

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/chasinglogic/dfm/internal/utils"
)

var (
	ErrInvalidHook = errors.New("hook was not a string or did not match expected format")
)

type Hooks map[string][]any

func (h Hooks) Execute(dir, hookName string) error {
	value, ok := h[hookName]
	if !ok {
		return nil
	}

	for _, hook := range value {
		args, err := parse(hook)
		if err != nil {
			return err
		}

		if err := utils.RunIn(dir, args...); err != nil {
			return err
		}
	}

	return nil
}

func parse(hook any) ([]string, error) {
	switch v := hook.(type) {
	case string:
		return []string{"/bin/sh", "-c", v}, nil
	case map[string]any:
		args := []string{}
		switch interpreter := v["interpreter"].(type) {
		case string:
			args = append(args, strings.Split(interpreter, " ")...)
		default:
			return nil, ErrInvalidHook
		}

		switch script := v["script"].(type) {
		case string:
			args = append(args, script)
		default:
			return nil, ErrInvalidHook
		}

		return args, nil
	default:
		fmt.Println(reflect.TypeOf(hook))
		return nil, ErrInvalidHook
	}
}
