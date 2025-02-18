package evaluator

import (
	"context"
	"fmt"
	"github/dyxgou/redis/internal/storage"
	"github/dyxgou/redis/internal/timer"
	"github/dyxgou/redis/pkg/ast"
	"log/slog"
)

const opSuccesful = "OK"

type Evaluator struct {
	s *storage.Storage
	t *timer.Ticker
}

func New(ctx context.Context) *Evaluator {
	e := &Evaluator{
		s: storage.New(),
	}

	t := timer.NewTicker()
	e.t = t
	go t.Init(ctx, e.deleteKey)

	return e
}

func (e *Evaluator) deleteKey(key string) {
	e.s.Delete(key)
}

func (e *Evaluator) Eval(cmd ast.Command) (string, error) {
	switch cmd := cmd.(type) {
	case *ast.GetCommand:
		return e.evalGetCommand(cmd)
	case *ast.GetDelCommand:
		return e.evalGetDelCommand(cmd)
	case *ast.GetSetCommand:
		return e.evalGetSetCommand(cmd)
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

func (e *Evaluator) evalGetDelCommand(gc *ast.GetDelCommand) (string, error) {
	val, ok := e.s.Get(gc.Key)
	if !ok {
		return storage.Nil.String(), nil
	}

	defer e.s.Delete(gc.Key)
	return val.String(), nil
}

func (e *Evaluator) evalGetSetCommand(gc *ast.GetSetCommand) (string, error) {
	if err := e.s.Set(gc.Key, gc.Value); err != nil {
		return "", err
	}

	res, _ := e.s.Get(gc.Key)

	return res.String(), nil
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
		e.t.Insert(sc.Key, sc.Ex)
	}

	return opSuccesful, nil
}
