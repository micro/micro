package accounts

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2/auth"
	accPb "github.com/micro/go-micro/v2/auth/service/proto/accounts"
	pb "github.com/micro/go-micro/v2/auth/service/proto/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	memStore "github.com/micro/go-micro/v2/store/memory"
)

const (
	storePrefix = "accounts/"
)

// Accounts processes RPC calls
type Accounts struct {
	Options auth.Options
}

// Init the auth
func (a *Accounts) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&a.Options)
	}

	// use the default store as a fallback
	if a.Options.Store == nil {
		a.Options.Store = store.DefaultStore
	}

	// noop will not work for auth
	if a.Options.Store.String() == "noop" {
		a.Options.Store = memStore.NewStore()
	}
}

// List returns all auth accounts
func (a *Accounts) List(ctx context.Context, req *accPb.ListAccountsRequest, rsp *accPb.ListAccountsResponse) error {
	// get the records from the store
	recs, err := a.Options.Store.Read(storePrefix, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to read from store: %v", err)
	}

	// unmarshal the records
	rsp.Accounts = make([]*pb.Account, 0, len(recs))
	for _, rec := range recs {
		var r *pb.Account
		if err := json.Unmarshal(rec.Value, &r); err != nil {
			return errors.InternalServerError("go.micro.auth", "Error to unmarshaling json: %v. Value: %v", err, string(rec.Value))
		}
		rsp.Accounts = append(rsp.Accounts, r)
	}

	return nil
}
