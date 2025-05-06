package service

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MomusWinner/MicroChat/internal/chatdb"
	"github.com/MomusWinner/MicroChat/internal/proxyproto"
	"github.com/Nerzal/gocloak/v13"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) Subscribe(ctx context.Context, request *proxyproto.SubscribeRequest) (*proxyproto.SubscribeResponse, error) {
	useruuid, _ := parseUUID(request.User)
	log.Printf("UserID: %v", useruuid)
	userId := pgtype.UUID{
		Bytes: useruuid,
		Valid: true,
	}

	log.Printf("UserID: %v", userId)

	_, err := s.storage.GetUserByID(ctx, userId)

	if errors.Is(err, sql.ErrNoRows) {
		log.Print(err)

		cloakUser, err := s.GetCloakUserById(ctx, request.User)
		if err != nil {
			log.Print(err)
			return RespondSubscribeError(CODE_UNAUTHORIZED, err.Error())
		}
		log.Printf("User: %v", cloakUser)
		log.Print(chatdb.CreateUserParams{
			ID:         userId,
			Username:   *cloakUser.Username,
			GivenName:  *cloakUser.FirstName,
			FamilyName: *cloakUser.LastName,
		})

		err = s.storage.CreateUser(ctx, chatdb.CreateUserParams{
			ID:         userId,
			Username:   *cloakUser.Username,
			GivenName:  *cloakUser.FirstName,
			FamilyName: *cloakUser.LastName,
		})

		if err != nil {
			return RespondSubscribeError(CODE_BAD_REQUEST, err.Error())
		}
	}

	// log.Print("User")
	// log.Print(request.User)
	count, err := s.storage.UserCanSubscribe(ctx, chatdb.UserCanSubscribeParams{
		UserID:  userId,
		Channel: request.Channel,
	})

	if err != nil {
		return RespondSubscribeError(CODE_INDERNAL_ERROR, err.Error())
	}

	if count < 1 {
		return RespondSubscribeError(CODE_PERMISSION_DENIED, "permission denied")
	}

	return &proxyproto.SubscribeResponse{}, nil
}

func (s *Service) Publish(ctx context.Context, request *proxyproto.PublishRequest) (*proxyproto.PublishResponse, error) {
	useruuid, _ := parseUUID(request.User)
	log.Printf("UserID: %v", useruuid)
	userId := pgtype.UUID{
		Bytes: useruuid,
		Valid: true,
	}

	log.Printf("UserID: %v", userId)

	_, err := s.storage.GetUserByID(ctx, userId)

	if errors.Is(err, sql.ErrNoRows) {
		log.Print(err)

		cloakUser, err := s.GetCloakUserById(ctx, request.User)
		if err != nil {
			log.Print(err)
			return RespondPublishError(CODE_UNAUTHORIZED, err.Error())
		}
		log.Printf("User: %v", cloakUser)
		log.Print(chatdb.CreateUserParams{
			ID:         userId,
			Username:   *cloakUser.Username,
			GivenName:  *cloakUser.FirstName,
			FamilyName: *cloakUser.LastName,
		})

		err = s.storage.CreateUser(ctx, chatdb.CreateUserParams{
			ID:         userId,
			Username:   *cloakUser.Username,
			GivenName:  *cloakUser.FirstName,
			FamilyName: *cloakUser.LastName,
		})

		if err != nil {
			return RespondPublishError(CODE_BAD_REQUEST, err.Error())
		}
	}

	// log.Print("User")
	// log.Print(request.User)
	count, err := s.storage.UserCanPublish(ctx, chatdb.UserCanPublishParams{
		UserID:  userId,
		Channel: request.Channel,
	})

	if err != nil {
		return RespondPublishError(CODE_INDERNAL_ERROR, err.Error())
	}

	if count < 1 {
		return RespondPublishError(CODE_PERMISSION_DENIED, "permission denied")
	}

	return &proxyproto.PublishResponse{}, nil
}

func (s *Service) GetCloakUserById(ctx context.Context, userId string) (*gocloak.User, error) {
	if s.token == nil || s.expiredAt.After(time.Now()) {
		token, err := s.cloakConn.LoginClient(ctx, s.cloakId, s.cloakSecret, s.cloakRealm)
		if err != nil {
			log.Fatal(err)
			return nil, nil
		}
		s.token = token
		s.expiredAt = time.Now().Add(time.Second * time.Duration(s.token.ExpiresIn))
	}
	user, err := s.cloakConn.GetUserByID(ctx, s.token.AccessToken, s.cloakRealm, userId)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}
	return user, nil
}

func parseUUID(src string) (dst [16]byte, err error) {
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		// assume invalid.
		return dst, fmt.Errorf("cannot parse UUID %v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}
