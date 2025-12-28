package oggophermart

import (
	"context"

	genUser "github.com/oleshko-g/oggophermart/internal/gen/user"
	"goa.design/clue/log"
)

// user service example implementation.
// The example methods log the requests and return zero values.
type userSvc struct{}

var _ genUser.Service = (*userSvc)(nil)

// NewUser returns the user service implementation.
func NewUser() genUser.Service {
	return &userSvc{}
}

// Register implements register.
func (s *userSvc) Register(ctx context.Context, p *genUser.LoginPass) (res *genUser.UserServiceResult, err error) {
	res = &genUser.UserServiceResult{}
	log.Printf(ctx, "user.register")
	return
}

// Login implements login.
func (s *userSvc) Login(ctx context.Context, p *genUser.LoginPass) (res *genUser.UserServiceResult, err error) {
	res = &genUser.UserServiceResult{}
	log.Printf(ctx, "user.login")
	return
}
