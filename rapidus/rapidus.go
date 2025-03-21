package rapidus

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/fouched/rapidus/session"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

type Rapidus struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	Routes   *chi.Mux
	Session  *scs.SessionManager
	config   config // no reason to export this
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
}

func (r *Rapidus) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := r.Init(pathConfig)
	if err != nil {
		return err
	}

	err = r.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := r.startLoggers()
	r.InfoLog = infoLog
	r.ErrorLog = errorLog
	r.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	r.Version = version
	r.RootPath = rootPath
	r.Routes = r.routes().(*chi.Mux)

	// setup config
	r.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"), // we probably don't need this with templ
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
	}

	// create session
	s := session.Session{
		CookieLifetime: r.config.cookie.lifetime,
		CookiePersist:  r.config.cookie.persist,
		CookieSecure:   r.config.cookie.secure,
		CookieName:     r.config.cookie.name,
		CookieDomain:   r.config.cookie.domain,
		SessionType:    r.config.sessionType,
	}
	r.Session = s.InitSession()

	return nil
}

func (r *Rapidus) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it does not exist
		err := r.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Rapidus) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     r.ErrorLog,
		Handler:      r.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	r.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	r.ErrorLog.Fatal(err)
}

func (r *Rapidus) checkDotEnv(path string) error {
	err := r.CreateFileIfNotExist(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (r *Rapidus) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
