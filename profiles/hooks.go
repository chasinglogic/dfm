package profiles

import (
	"os"
	"os/exec"

	"github.com/chasinglogic/dfm/logger"
	"github.com/google/shlex"
)

type Hooks map[string][]interface{}

func (h Hooks) RunHook(name, dir string, dryRun bool) error {
	hooks, ok := h[name]
	if !ok {
		logger.Debug.Printf("no hook defined for: %s", name)
		return nil
	}

	for _, hook := range hooks {
		var command []string
		var err error

		if hookStr, ok := hook.(string); ok {
			command, err = shlex.Split(hookStr)
			if err != nil {
				return err
			}

		} else if hookMap, ok := hook.(map[string]interface{}); ok {
			interpreter, hasInterpreter := hookMap["interpreter"]
			script, hasScript := hookMap["script"]
			if !hasScript || !hasInterpreter {
				logger.Debug.Printf("%s hook is missing script or interpreter", name)
				continue
			}

			if _, ok := interpreter.(string); !ok {
				logger.Debug.Printf("%s hook interpreter is not a string", name)
				continue
			}

			if _, ok := script.(string); !ok {
				logger.Debug.Printf("%s hook script is not a string", name)
				continue
			}

			interpreterArgs, err := shlex.Split(interpreter.(string))
			if err != nil {
				return err
			}

			command = append(
				interpreterArgs,
				script.(string),
			)
		}

		logger.Debug.Printf("Executing command: %s\n", command)
		proc := exec.Command(command[0], command[1:]...)
		proc.Stdout = os.Stdout
		proc.Stdin = os.Stdin
		proc.Stderr = os.Stderr
		proc.Dir = dir

		proc.Start()
		err = proc.Wait()

		if err != nil {
			return err
		}
	}

	return nil
}
