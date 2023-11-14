package middleware

import (
	"github.com/fasthttp/session/v2"
	"github.com/valyala/fasthttp"
	"mmplat/internal/util"
)

type UserSession struct {
	Session *session.Session
}

const throttle = 5

// Csrf csrf
func (s *UserSession) Csrf(ctx *fasthttp.RequestCtx) error {
	// inits session
	storage, _ := s.Session.Get(ctx)
	storage.Set("csrf", util.RandStr(util.DefRndLen))
	defer s.Session.Save(ctx, storage)

	return nil
}

// Throttle TODO finish
func (s *UserSession) Throttle(ctx *fasthttp.RequestCtx) error {
	/*storage, _ := s.Session.Get(ctx)
	throttleAuth := storage.Get("throttleAuth")
	if throttleAuth >= throttle {
		return errors.New("login attemps exceded limit")
	}*/
	return nil
}