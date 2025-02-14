package evaluator

import (
	"fmt"
	"github/dyxgou/redis/internal/storage"
	"github/dyxgou/redis/internal/timer"
	"github/dyxgou/redis/pkg/ast"
	"log/slog"
)

const opSuccesful = "OK"

type Evaluator struct {
	s *storage.Storage
	t *timer.Timer
}

func New() *Evaluator {
	return &Evaluator{
		s: storage.New(),
		t: timer.New(),
	}
}

func (e *Evaluator) Eval(cmd ast.Command) (string, error) {
	switch cmd := cmd.(type) {
	case *ast.GetCommand:
		return e.evalGetCommand(cmd)
	case *ast.SetCommand:
		return e.evalSetCommand(cmd)
	}

	return "", fmt.Errorf("command not supported for evaluation. got=%T", cmd)
}

func (e *Evaluator) evalGetCommand(gc *ast.GetCommand) (string, error) {
	val, ok := e.s.Get(gc.Key)
	slog.Info("GET Command", "val", val)

	if !ok {
		return storage.Nil.String(), nil
	}

	return val.String(), nil
}

func (e *Evaluator) evalSetCommand(sc *ast.SetCommand) (string, error) {
	exists := e.s.Exists(sc.Key)
	if sc.Xx && !exists {
		return "", fmt.Errorf("flag XX set when key does not exists. key=%q", sc.Key)
	}

	if sc.Nx && exists {
		return "", fmt.Errorf("flag NX set when key does exists. key=%q", sc.Key)
	}

	if err := e.s.Set(sc.Key, sc.Value); err != nil {
		return "", err
	}

	if sc.Ex >= 1 {
		e.t.Insert(timer.NewTimestamp(sc.Key, sc.Ex))
	}

	return opSuccesful, nil
}
