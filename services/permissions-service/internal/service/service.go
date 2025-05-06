package service

import (
	"context"
	"time"

	"github.com/MomusWinner/MicroChat/internal/chatdb"
	"github.com/MomusWinner/MicroChat/internal/proxyproto"
	"github.com/MomusWinner/MicroChat/services/permissions-service/internal/config"
	"github.com/Nerzal/gocloak/v13"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	proxyproto.UnimplementedCentrifugoProxyServer
	dbConn      *pgxpool.Pool
	cloakConn   *gocloak.GoCloak
	cloakSecret string
	cloakId     string
	cloakRealm  string
	storage     *chatdb.Queries

	token     *gocloak.JWT
	expiredAt time.Time
}

func New(conf *config.Config) (*Service, error) {
	cloakConn := gocloak.NewClient(conf.KeyCloakURL)

	connCfg, err := pgxpool.ParseConfig(conf.DatabaseURL)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), connCfg)
	if err != nil {
		return nil, err
	}

	return &Service{
		dbConn:      conn,
		cloakId:     conf.CloakClientId,
		cloakRealm:  conf.CloakRealm,
		cloakSecret: conf.CloakSecret,
		cloakConn:   cloakConn,
		storage:     chatdb.New(conn),
	}, nil
}
