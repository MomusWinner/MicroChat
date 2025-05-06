package service

import "github.com/MomusWinner/MicroChat/internal/proxyproto"

const (
	CODE_INDERNAL_ERROR    = 100
	CODE_PERMISSION_DENIED = 103
	CODE_BAD_REQUEST       = 107
	CODE_UNAUTHORIZED      = 101
)

func RespondSubscribeError(code uint32, msg string) (*proxyproto.SubscribeResponse, error) {
	return &proxyproto.SubscribeResponse{
		Error: &proxyproto.Error{
			Code:    code,
			Message: msg,
		},
	}, nil
}

func RespondPublishError(code uint32, msg string) (*proxyproto.PublishResponse, error) {
	return &proxyproto.PublishResponse{
		Error: &proxyproto.Error{
			Code:    code,
			Message: msg,
		},
	}, nil
}
