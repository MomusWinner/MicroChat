package service

import "github.com/MomusWinner/MicroChat/internal/proxyproto"

func RespondError(code uint32, msg string) (*proxyproto.ConnectResponse, error) {
	return &proxyproto.ConnectResponse{
		Error: &proxyproto.Error{
			Code:    code,
			Message: msg,
		},
	}, nil
}
