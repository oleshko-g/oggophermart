package gophermart

import (
	"context"

	user "github.com/oleshko-g/oggophermart/internal/gen/user"
	"goa.design/clue/log"
)

// user service example implementation.
// The example methods log the requests and return zero values.
type usersrvc struct{}

// NewUser returns the user service implementation.
func NewUser() user.Service {
	return &usersrvc{}
}

// Register implements register.
func (s *usersrvc) Register(ctx context.Context, p *user.LoginPass) (res *user.UserServiceResult, err error) {
	res = &user.UserServiceResult{}
	log.Printf(ctx, "user.register")
	return
}

// Login implements login.
func (s *usersrvc) Login(ctx context.Context, p *user.LoginPass) (res *user.UserServiceResult, err error) {
	res = &user.UserServiceResult{}
	log.Printf(ctx, "user.login")
	return
}
