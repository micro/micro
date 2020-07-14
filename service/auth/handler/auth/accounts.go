package auth

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/micro/v2/internal/namespace"
)

// List returns all auth accounts
func (a *Auth) List(ctx context.Context, req *pb.ListAccountsRequest, rsp *pb.ListAccountsResponse) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.Options{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.auth.Accounts.List", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.auth.Accounts.List", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth.Accounts.List", err.Error())
	}

	// setup the defaults incase none exist
	a.setupDefaultAccount(req.Options.Namespace)

	// get the records from the store
	key := strings.Join([]string{storePrefixAccounts, req.Options.Namespace, ""}, joinKey)
	recs, err := a.Options.Store.Read(key, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to read from store: %v", err)
	}

	// unmarshal the records
	var accounts = make([]*auth.Account, 0, len(recs))
	for _, rec := range recs {
		var r *auth.Account
		if err := json.Unmarshal(rec.Value, &r); err != nil {
			return errors.InternalServerError("go.micro.auth", "Error to unmarshaling json: %v. Value: %v", err, string(rec.Value))
		}
		accounts = append(accounts, r)
	}

	// serialize the accounts
	rsp.Accounts = make([]*pb.Account, 0, len(recs))
	for _, a := range accounts {
		rsp.Accounts = append(rsp.Accounts, serializeAccount(a))
	}

	return nil
}

// Delete an auth account
func (a *Auth) Delete(ctx context.Context, req *pb.DeleteAccountRequest, rsp *pb.DeleteAccountResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.auth.Auth.Delete", "Missing ID")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.Options{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.auth.Accounts.Delete", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.auth.Accounts.Delete", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth.Accounts.Delete", err.Error())
	}

	// check the account exists
	key := strings.Join([]string{storePrefixAccounts, req.Options.Namespace, req.Id}, joinKey)
	if _, err := a.Options.Store.Read(key); err == store.ErrNotFound {
		return errors.BadRequest("go.micro.auth.Accounts.Delete", "Account not found with this ID")
	} else if err != nil {
		return errors.BadRequest("go.micro.auth.Accounts.Delete", "Error querying accounts: %v", err)
	}

	// delete the account
	if err := a.Options.Store.Delete(key); err != nil {
		return errors.BadRequest("go.micro.auth.Accounts.Delete", "Error deleting account: %v", err)
	}

	return nil
}

func serializeAccount(a *auth.Account) *pb.Account {
	return &pb.Account{
		Id:       a.ID,
		Type:     a.Type,
		Scopes:   a.Scopes,
		Issuer:   a.Issuer,
		Metadata: a.Metadata,
	}
}
