package actions

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"go.flipt.io/flipt/rpc/flipt"
	"strconv"
)

func checkFeature(ctx context.Context, key string) (*string, error) {
	res, err := ToggleService.Evaluate(ctx, &flipt.EvaluationRequest{
		FlagKey:  key,
		EntityId: uuid.NewV4().String(), // random requester
	})

	if err != nil {
		return nil, err
	}

	return &res.Value, nil
}

func isEnabled(ctx context.Context, key string) (bool, error) {
	res, err := checkFeature(ctx, key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(*res)
}
