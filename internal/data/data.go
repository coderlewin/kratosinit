package data

import (
	"context"
	"github.com/coderlewin/kratosinit/internal/biz"
	"github.com/coderlewin/kratosinit/internal/conf"
	"github.com/coderlewin/kratosinit/internal/data/gorm_gen/dal"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewMysql, NewTransaction, NewUserRepo, NewRedis)

// Data .
type Data struct {
	db    *dal.Query
	cache redis.Cmdable
}

type contextTxKey struct{}

func (d *Data) InTx(ctx context.Context, f func(ctx context.Context) error) error {
	return d.db.Transaction(func(tx *dal.Query) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return f(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *dal.Query {
	tx, ok := ctx.Value(contextTxKey{}).(*dal.Query)
	if ok {
		return tx
	}
	return d.db
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

// NewData .
func NewData(c *conf.Data, qry *dal.Query, cache redis.Cmdable) (*Data, error) {
	d := &Data{db: qry, cache: cache}
	return d, nil
}

func NewMysql(c *conf.Data, logger log.Logger) (*dal.Query, func(), error) {
	logx := log.NewHelper(logger)
	db, err := gorm.Open(mysql.Open(c.Database.Source))
	sqldb, err := db.DB()
	if err != nil {
		logx.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	sqldb.SetMaxIdleConns(int(c.Database.MaxIdleConn))
	sqldb.SetMaxOpenConns(int(c.Database.MaxOpenConn))
	sqldb.SetConnMaxLifetime(c.Database.ConnMaxLifetime.AsDuration())

	query := dal.Use(db)
	dal.SetDefault(db)

	return query, func() {
		if err := sqldb.Close(); err != nil {
			logx.Errorf("failed to close db connection: %v", err)
		}
	}, nil
}

func NewRedis(c *conf.Data, logger log.Logger) (redis.Cmdable, func(), error) {
	logx := log.NewHelper(logger)
	client := redis.NewClient(&redis.Options{
		Addr:            c.Redis.GetAddr(),
		Password:        c.Redis.GetPassword(),
		DB:              int(c.Redis.GetDb()),
		ReadTimeout:     c.Redis.GetReadTimeout().AsDuration(),
		WriteTimeout:    c.Redis.GetWriteTimeout().AsDuration(),
		MaxIdleConns:    10,
		MaxActiveConns:  20,
		ConnMaxIdleTime: time.Second * 3600,
		ConnMaxLifetime: time.Second * 1800,
	})

	cleanFunc := func() {
		if err := client.Close(); err != nil {
			logx.Errorf("failed to close db connection: %v", err)
		}
	}
	return client, cleanFunc, nil
}
