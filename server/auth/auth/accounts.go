package auth

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
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

	// setup the defaults incase none exist
	a.setupDefaultAccount(req.Options.Namespace)

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, req.Options.Namespace, "auth.Accounts.List"); err != nil {
		return err
	}

	// get the records from the store
	key := strings.Join([]string{storePrefixAccounts, req.Options.Namespace, ""}, joinKey)
	recs, err := a.Options.Store.Read(key, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("auth.Accounts.List", "Unable to read from store: %v", err)
	}

	// unmarshal the records
	var accounts = make([]*auth.Account, 0, len(recs))
	for _, rec := range recs {
		var r *auth.Account
		if err := json.Unmarshal(rec.Value, &r); err != nil {
			return errors.InternalServerError("auth.Accounts.List", "Error to unmarshaling json: %v. Value: %v", err, string(rec.Value))
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
		return errors.BadRequest("auth.Accounts.Delete", "Missing ID")
	}

	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("auth.Accounts.Delete", "Unauthorized")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.Options{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request can access this namespace
	if err := namespace.AuthorizeAdmin(ctx, req.Options.Namespace, "auth.Accounts.Delete"); err != nil {
		return err
	}

	// check the account exists
	accToDelete, err := a.getAccountForID(req.Id, req.Options.Namespace, "auth.Accounts.Delete")
	if err != nil {
		return err
	}

	if req.Id == acc.ID || req.Id == acc.Name {
		return errors.BadRequest("auth.Accounts.Delete", "Can't delete your own account")
	}

	// delete the refresh token linked to the account
	tok, err := a.refreshTokenForAccount(req.Options.Namespace, accToDelete.ID)
	if err != nil {
		return errors.InternalServerError("auth.Accounts.Delete", "Error finding refresh token")
	}
	refreshKey := strings.Join([]string{storePrefixRefreshTokens, req.Options.Namespace, accToDelete.ID, tok}, joinKey)
	if err := a.Options.Store.Delete(refreshKey); err != nil {
		return errors.InternalServerError("auth.Accounts.Delete", "Error deleting refresh token: %v", err)
	}

	key := strings.Join([]string{storePrefixAccounts, req.Options.Namespace, accToDelete.ID}, joinKey)
	// delete the account
	if err := a.Options.Store.Delete(key); err != nil {
		return errors.BadRequest("auth.Accounts.Delete", "Error deleting account: %v", err)
	}
	keyByName := strings.Join([]string{storePrefixAccountsByName, req.Options.Namespace, accToDelete.Name}, joinKey)
	// delete the account
	if err := a.Options.Store.Delete(keyByName); err != nil {
		return errors.BadRequest("auth.Accounts.Delete", "Error deleting account: %v", err)
	}

	// Clear the namespace cache, since the accounts for this namespace could now be empty
	a.Lock()
	delete(a.namespaces, req.Options.Namespace)
	a.Unlock()

	return nil
}

func hasScope(scope string, scopes []string) bool {
	for _, s := range scopes {
		if scope == s {
			return true
		}
	}
	return false
}

// ChangeSecret by providing a refresh token and a new secret
func (a *Auth) ChangeSecret(ctx context.Context, req *pb.ChangeSecretRequest, rsp *pb.ChangeSecretResponse) error {
	if len(req.NewSecret) == 0 {
		return errors.BadRequest("auth.Auth.ChangeSecret", "New secret should not be blank")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.Options{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace, "auth.Accounts.ChangeSecret"); err != nil {
		return err
	}

	acc, err := a.getAccountForID(req.Id, req.Options.Namespace, "auth.Accounts.ChangeSecret")
	if err != nil {
		return err
	}

	if !secretsMatch(acc.Secret, req.OldSecret) {
		return errors.BadRequest("auth.Accounts.ChangeSecret", "Secret not correct")
	}

	// hash the secret
	secret, err := hashSecret(req.NewSecret)
	if err != nil {
		return errors.InternalServerError("auth.Accounts.ChangeSecret", "Unable to hash password: %v", err)
	}
	acc.Secret = secret

	// marshal to json
	bytes, err := json.Marshal(acc)
	if err != nil {
		return errors.InternalServerError("auth.Accounts.ChangeSecret", "Unable to marshal json: %v", err)
	}

	key := strings.Join([]string{storePrefixAccounts, acc.Issuer, acc.ID}, joinKey)
	// write to the store
	if err := a.Options.Store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError("auth.Accounts.ChangeSecret", "Unable to write account to store: %v", err)
	}
	usernameKey := strings.Join([]string{storePrefixAccountsByName, acc.Issuer, acc.Name}, joinKey)
	if err := a.Options.Store.Write(&store.Record{Key: usernameKey, Value: bytes}); err != nil {
		return errors.InternalServerError("auth.Accounts.ChangeSecret", "Unable to write account to store: %v", err)
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
		Name:     a.Name,
	}
}
