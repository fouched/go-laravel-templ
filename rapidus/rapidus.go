package rapidus

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/fouched/rapidus/cache"
	"github.com/fouched/rapidus/session"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

type Rapidus struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Session       *scs.SessionManager
	DB            Database
	config        config // no reason to export this
	EncryptionKey string
	Cache         cache.Cache
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
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

	// connect to database if specified
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := r.OpenDB(os.Getenv("DATABASE_TYPE"), r.BuildDSN())
		//TODO: build in retry
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		r.DB = Database{
			Type: os.Getenv("DATABASE_TYPE"),
			Pool: db,
		}
	}

	if os.Getenv("CACHE") == "redis" {
		r.Cache = r.createRedisClient()
	}

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
		database: databaseConfig{
			dsn:      r.BuildDSN(),
			database: os.Getenv("DATABASE_TYPE"),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

	// create session
	s := session.Session{
		CookieLifetime: r.config.cookie.lifetime,
		CookiePersist:  r.config.cookie.persist,
		CookieSecure:   r.config.cookie.secure,
		CookieName:     r.config.cookie.name,
		CookieDomain:   r.config.cookie.domain,
		SessionType:    r.config.sessionType,
		DBPool:         r.DB.Pool,
	}
	r.Session = s.InitSession()

	// encryption key
	r.EncryptionKey = os.Getenv("KEY")

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

	defer r.DB.Pool.Close()

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

func (r *Rapidus) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	}

	return dsn
}

func (r *Rapidus) createRedisClient() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn: redis.NewClient(&redis.Options{
			Addr:     r.config.redis.host,
			Password: r.config.redis.password,
			DB:       0, // use default DB
		}),
		Prefix: r.config.redis.prefix,
	}
	return &cacheClient
}
