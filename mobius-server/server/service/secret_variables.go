package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/notawar/mobius/mobius-server/server/contexts/ctxerr"
	"github.com/notawar/mobius/mobius-server/server/mobius"
)

const (
	SecretVariablePrefix = "MOBIUS_SECRET_"
	SecretVariableMaxLen = 255
)

// //////////////////////////////////////////////////////////////////////////////
// Secret variables
// //////////////////////////////////////////////////////////////////////////////

type secretVariablesRequest struct {
	DryRun          bool                    `json:"dry_run"`
	SecretVariables []mobius.SecretVariable `json:"secrets"`
}

type secretVariablesResponse struct {
	Err error `json:"error,omitempty"`
}

func (r secretVariablesResponse) Error() error { return r.Err }

func secretVariablesEndpoint(ctx context.Context, request interface{}, svc mobius.Service) (mobius.Errorer, error) {
	req := request.(*secretVariablesRequest)
	err := svc.CreateSecretVariables(ctx, req.SecretVariables, req.DryRun)
	return secretVariablesResponse{Err: err}, nil
}

func (svc *Service) CreateSecretVariables(ctx context.Context, secretVariables []mobius.SecretVariable, dryRun bool) error {
	// Do authorization check first so that we don't have to worry about it later in the flow.
	if err := svc.authz.Authorize(ctx, &mobius.SecretVariable{}, mobius.ActionWrite); err != nil {
		return err
	}

	privateKey := svc.config.Server.PrivateKey
	if testSetEmptyPrivateKey {
		privateKey = ""
	}

	if len(privateKey) == 0 {
		return ctxerr.Wrap(ctx,
			&mobius.BadRequestError{Message: "Couldn't save secret variables. Missing required private key. Learn how to configure the private key here: https://mobiusmdm.com/learn-more-about/mobius-server-private-key"})
	}

	// Preprocess: strip MOBIUS_SECRET_ prefix from variable names
	for i, secretVariable := range secretVariables {
		secretVariables[i].Name = mobius.Preprocess(strings.TrimPrefix(secretVariable.Name, SecretVariablePrefix))
	}

	// Validation
	for _, secretVariable := range secretVariables {
		if len(secretVariable.Name) == 0 {
			return ctxerr.Wrap(ctx,
				mobius.NewInvalidArgumentError("name", "secret variable name cannot be empty"))
		}
		if len(secretVariable.Name) > SecretVariableMaxLen {
			return ctxerr.Wrap(ctx,
				mobius.NewInvalidArgumentError("name", fmt.Sprintf("secret variable name is too long: %s", secretVariable.Name)))
		}
	}

	if dryRun {
		return nil
	}

	if err := svc.ds.UpsertSecretVariables(ctx, secretVariables); err != nil {
		return ctxerr.Wrap(ctx, err, "saving secret variables")
	}
	return nil
}
