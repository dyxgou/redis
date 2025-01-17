package evaluator

import (
	"context"
	"fmt"
	"github/dyxgou/redis/internal/storage"
	"github/dyxgou/redis/pkg/ast"
)

type evaluator struct {
	ctx context.Context
	s   *storage.Storage
}

func new(ctx context.Context) *evaluator {
	return &evaluator{
		ctx: ctx,
		s:   storage.New(),
	}
}

func (e *evaluator) eval(n ast.Command) (storage.Value, error) {
	switch n := n.(type) {
	case *ast.GetCommand:
		return e.evalGet(n)
	}

	return nil, nil
}

func (e *evaluator) evalGet(gc *ast.GetCommand) (storage.Value, error) {
	v, ok := e.s.Get(gc.Key)
	if !ok {
		return storage.Nil, nil
	}

	return v, nil
}

func (e *evaluator) evalSet(sc *ast.SetCommand) (storage.Value, error) {
	if sc.Nx && e.s.Exists(sc.Key) {
		return nil, fmt.Errorf("key=%q already has a value stored", sc.Key)
	}

	if sc.Xx && !e.s.Exists(sc.Key) {
		return nil, fmt.Errorf("key=%q does not has a value stored", sc.Key)
	}

	if sc.Ex > 0 {

	}

	e.s.Set(sc.Key, sc.Value)
	return nil, nil
}
