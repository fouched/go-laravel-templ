package session

import (
	"github.com/alexedwards/scs/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieSecure   string
	CookieName     string
	CookieDomain   string
	SessionType    string
}

func (s *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// how long should sessions last
	minutes, err := strconv.Atoi(s.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	// should cookies persist
	persist = strings.ToLower(s.CookiePersist) == "true"

	// must cookies be secure
	secure = strings.ToLower(s.CookieSecure) == "true"

	// create session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = s.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = s.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// which session store
	switch strings.ToLower(s.SessionType) {
	case "redis":
	case "mysql", "mariadb":
	case "postgres", "postgresql":
	default:
		// cookie
	}

	return session
}
