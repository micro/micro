package handler

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/auth"
	"github.com/micro/go-micro/util/log"

	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/store"

	pb "github.com/micro/go-micro/auth/service/proto"
)

// New returns an instance of Handler
func New() *Handler {
	return &Handler{
		store: *cmd.DefaultOptions().Store,
	}
}

var (
	// Duration the service account is valid for
	Duration = time.Hour * 24
)

// Handler processes RPC calls
type Handler struct {
	store store.Store
}

// Generate creates a new service account in the store
func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if req.Account == nil {
		return errors.BadRequest("go.micro.auth", "service account required")
	}

	parent := req.Account.Parent
	if parent == nil {
		return errors.BadRequest("go.micro.auth", "parent required")
	}
	if parent.Name == "" || parent.Type == "" {
		return errors.BadRequest("go.micro.auth", "invalid parent")
	}

	// generate the token
	token, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	// key for the store
	key := fmt.Sprintf("%v/%v", prefixForResource(parent), token.String())

	// construct the service account
	sa := auth.Account{
		Created:  time.Now(),
		Expiry:   time.Now().Add(Duration),
		Metadata: req.Account.Metadata,
	}

	// add the roles
	sa.Roles = make([]*auth.Role, len(req.Account.Roles))
	for i, r := range req.Account.Roles {
		sa.Roles[i] = &auth.Role{Name: r.Name}

		if r.Resource != nil {
			sa.Roles[i].Resource = &auth.Resource{
				Name: r.Resource.Name,
				Type: r.Resource.Type,
			}
		}
	}

	// encode the data to bytes
	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	if err := e.Encode(sa); err != nil {
		return err
	}

	// write to the store
	err = h.store.Write(&store.Record{
		Key:    key,
		Value:  buf.Bytes(),
		Expiry: Duration,
	})
	if err != nil {
		return err
	}
	log.Infof("Created service account: %v", key)

	// encode the response
	rsp.Account = &pb.Account{
		Created:  sa.Created.Unix(),
		Expiry:   sa.Expiry.Unix(),
		Metadata: sa.Metadata,
		Token:    token.String(),
		Roles:    req.Account.Roles,
	}

	return nil
}

// Validate retrieves a service account from the store
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	if req.Token == "" {
		return errors.BadRequest("go.micro.auth", "token required")
	}

	// lookup the record by token
	records, err := h.store.Read(req.Token, store.ReadSuffix())
	if err == store.ErrNotFound || len(records) == 0 {
		return errors.Unauthorized("go.micro.auth", "invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "error reading store")
	}

	// decode the result
	b := bytes.NewBuffer(records[0].Value)
	decoder := gob.NewDecoder(b)
	var sa auth.Account
	err = decoder.Decode(&sa)

	// encode the response
	rsp.Account = &pb.Account{
		Created:  sa.Created.Unix(),
		Expiry:   sa.Expiry.Unix(),
		Metadata: sa.Metadata,
		Token:    req.Token,
		Roles:    make([]*pb.Role, len(sa.Roles)),
	}
	for i, r := range sa.Roles {
		rsp.Account.Roles[i] = &pb.Role{Name: r.Name}

		if r.Resource != nil {
			rsp.Account.Roles[i].Resource = &pb.Resource{
				Name: r.Resource.Name,
				Type: r.Resource.Type,
			}
		}
	}

	log.Infof("Validated service account: %v", records[0].Key)
	return nil
}

// Revoke deletes the service account
func (h *Handler) Revoke(ctx context.Context, req *pb.RevokeRequest, rsp *pb.RevokeResponse) error {
	if req.Token == "" {
		return errors.BadRequest("go.micro.auth", "token required")
	}

	records, err := h.store.Read(req.Token, store.ReadSuffix())
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "error reading store")
	}
	for _, r := range records {
		if err := h.store.Delete(r.Key); err != nil {
			return errors.InternalServerError("go.micro.auth", "error deleting from store")
		}
		log.Infof("Revoked service account: %v", r.Key)
	}

	return nil
}

// prefixForResource is used is the store's key name, e.g. user/asim@micro.mu || service/go.micro.srv.auth
func prefixForResource(r *pb.Resource) string {
	return fmt.Sprintf("%v/%v", r.Type, r.Name)
}
