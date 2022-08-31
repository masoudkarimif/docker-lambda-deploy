package action

import (
	"context"
)

type Action struct{}

type storage interface {
	UpdateCode(ctx context.Context) error
}

type function interface {
	UpdateCode(ctx context.Context) error
}

func (*Action) Run(ctx context.Context, fc function, st storage) error {
	err := st.UpdateCode(ctx)
	if err != nil {
		return err
	}

	err = fc.UpdateCode(ctx)
	if err != nil {
		return err
	}

	return nil
}
