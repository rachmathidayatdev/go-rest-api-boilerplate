package util

import (
	"context"
)

//SessionKey const
const SessionKey = 1

//Session struct
type Session struct {
	CID    string
	Logger Logger
}

//SessionCid func
func SessionCid(ctx context.Context) string {
	session, ok := ctx.Value(SessionKey).(*Session)

	// Handle if session middleware is not used
	if !ok {
		return ""
	}

	return session.CID
}

//SessionLogger func
func SessionLogger(ctx context.Context) Logger {
	session, ok := ctx.Value(SessionKey).(*Session)

	// Handle if session middleware is not used
	if !ok {
		return Log
	}

	return session.Logger
}

//NewSessionCtx func
func NewSessionCtx(cid string, log Logger) context.Context {
	session := Session{
		cid,
		log,
	}
	return context.WithValue(context.Background(), SessionKey, &session)
}
