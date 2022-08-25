package pgs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/scryinfo/dot/dot"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

const (
	//ConnWrapperTypeID type id
	ConnWrapperTypeID = "ffc08507-dd5f-456c-84ea-cdae00b220bf"
)

type config struct {
	Host     dot.StringFromEnv `json:"host"`
	Port     dot.StringFromEnv `json:"port"`
	User     dot.StringFromEnv `json:"user"`
	Password dot.StringFromEnv `json:"password"`
	Database dot.StringFromEnv `json:"database"`
	ShowSQL  bool              `json:"showSql"`
}

//ConnWrapper connect wrapper
type ConnWrapper struct {
	db   *bun.DB
	conf config
}

func (c *ConnWrapper) Create(dot.Line) error {
	// dsn := "postgres://root:aBc123@localhost:5432/postgres?sslmode=disable&timeout=5"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?&timeout=5&sslmode=disable",
		string(c.conf.User),
		string(c.conf.Password),
		string(c.conf.Host),
		string(c.conf.Port),
		string(c.conf.Database),
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	db.AddQueryHook(bundebug.NewQueryHook())
	c.db = db

	return nil
}

func (c *ConnWrapper) AfterAllDestroy(dot.Line) {
	if c.db != nil {
		_ = c.db.Close()
		c.db = nil
	}
}

//GetDb get db
func (c *ConnWrapper) GetDb() *bun.DB {
	return c.db
}

//TestConn test the connect
func (c *ConnWrapper) TestConn() bool {
	n := -1
	c.db.Query("select 1")
	return n == 1
}

//construct dot
func newConnWrapper(conf []byte) (dot.Dot, error) {
	dconf := &config{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}
	d := &ConnWrapper{conf: *dconf}
	return d, err
}

//ConnWrapperTypeLives make all type lives
func ConnWrapperTypeLives() []*dot.TypeLives {
	return []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: ConnWrapperTypeID, NewDoter: newConnWrapper},
	}}
}

//GenerateConnWrapper this func is for test
func GenerateConnWrapper(conf string) *ConnWrapper {
	conn := &ConnWrapper{}
	_ = json.Unmarshal([]byte(conf), &conn.conf)
	_ = conn.Create(nil)
	return conn
}

//GenerateConnWrapperByDb this func is for test
func GenerateConnWrapperByDb(db *bun.DB) *ConnWrapper {
	conn := &ConnWrapper{db, config{}}
	return conn
}

type pgLogger struct{}

func (d pgLogger) BeforeQuery(c context.Context, _ *pgdriver.Listener) (context.Context, error) {
	return c, nil
}

func (d pgLogger) AfterQuery(_ context.Context, q *bun.DB) error {
	dot.Logger().Debug(func() string {
		// q.Formatter().FormatQuery()
		return ""
	})
	return nil
}
