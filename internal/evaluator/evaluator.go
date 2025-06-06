package evaluator

import (
	"context"
	"fmt"
	"github/dyxgou/redis/internal/storage"
	"github/dyxgou/redis/internal/timer"
	"github/dyxgou/redis/pkg/ast"
	"strconv"
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
	case *ast.GetExCommand:
		return e.evalGetExCommand(cmd)
	case *ast.IncrCommand:
		return e.evalIncrCommand(cmd)
	case *ast.IncrByCommand:
		return e.evalIncrByCommand(cmd)
	case *ast.DecrCommand:
		return e.evalDecrCommand(cmd)
	case *ast.DecrByCommand:
		return e.evalDecrByCommand(cmd)
	case *ast.ExistsCommand:
		return e.evalExistsCommand(cmd)
	}

	return "", fmt.Errorf("command not supported for evaluation. got=%T", cmd)
}

func (e *Evaluator) evalIncrCommand(inc *ast.IncrCommand) (string, error) {
	val, ok := e.s.Get(inc.Key)
	if !ok {
		return "", fmt.Errorf("key=%q not found", inc.Key)
	}

	switch intVal := val.(type) {
	case *storage.Int:
		intVal.Value++
		return strconv.Itoa(intVal.Value), nil
	case *storage.Int64:
		intVal.Value++
		return strconv.FormatInt(intVal.Value, 10), nil
	}

	return "", fmt.Errorf("val kind is not numeric. val=%q", val.String())
}

func (e *Evaluator) evalExistsCommand(ec *ast.ExistsCommand) (string, error) {
	ok := e.s.Exists(ec.Key)

	return strconv.FormatBool(ok), nil
}

func (e *Evaluator) evalDecrCommand(dec *ast.DecrCommand) (string, error) {
	val, ok := e.s.Get(dec.Key)
	if !ok {
		return "", fmt.Errorf("key=%q not found", dec.Key)
	}

	switch intVal := val.(type) {
	case *storage.Int:
		intVal.Value--
		return strconv.Itoa(intVal.Value), nil
	case *storage.Int64:
		intVal.Value--
		return strconv.FormatInt(intVal.Value, 10), nil
	}

	return "", fmt.Errorf("val kind is not numeric. val=%q", val.String())
}

func (e *Evaluator) evalDecrByCommand(dec *ast.DecrByCommand) (string, error) {
	val, ok := e.s.Get(dec.Key)
	if !ok {
		return "", fmt.Errorf("key=%q not found", dec.Key)
	}

	switch intVal := val.(type) {
	case *storage.Int:
		intVal.Value -= dec.Decrement
		return strconv.Itoa(intVal.Value), nil
	case *storage.Int64:
		intVal.Value -= int64(dec.Decrement)
		return strconv.FormatInt(intVal.Value, 10), nil
	}

	return "", fmt.Errorf("val kind is not numeric. val=%q", val.String())
}

func (e *Evaluator) evalIncrByCommand(inc *ast.IncrByCommand) (string, error) {
	val, ok := e.s.Get(inc.Key)
	if !ok {
		return "", fmt.Errorf("key=%q not found", inc.Key)
	}

	switch intVal := val.(type) {
	case *storage.Int:
		intVal.Value += inc.Increment
		return strconv.Itoa(intVal.Value), nil
	case *storage.Int64:
		intVal.Value += int64(inc.Increment)
		return strconv.FormatInt(intVal.Value, 10), nil
	}

	return "", fmt.Errorf("val kind is not numeric. val=%q", val.String())
}

func (e *Evaluator) evalGetCommand(gc *ast.GetCommand) (string, error) {
	val, ok := e.s.Get(gc.Key)

	if !ok {
		return storage.Nil.String(), nil
	}

	return val.String(), nil
}

func (e *Evaluator) evalGetExCommand(gc *ast.GetExCommand) (string, error) {
	if gc.Ex < 1 {
		return "", fmt.Errorf("flag EX should have a value greater than 1. EX=%d", gc.Ex)
	}
	e.t.Insert(gc.Key, gc.Ex)
	val, ok := e.s.Get(gc.Key)

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
