package service

import (
	vh "github.com/mikerumy/vhosting"
	storage "github.com/mikerumy/vhosting/pkg/storage/interfaces"
)

type AuthService struct {
	stor storage.Authorization
}

func NewAuthService(stor storage.Authorization) *AuthService {
	return &AuthService{stor: stor}
}

func (s *AuthService) POSTSession(session vh.Session) error {
	return s.stor.POSTSession(session)
}

func (s *AuthService) GETNamePass(namepass vh.NamePass) error {
	return s.stor.GETNamePass(namepass)
}

func (s *AuthService) DELETECurrentSession(cookieValue string) error {
	return s.stor.DELETECurrentSession(cookieValue)
}

func (s *AuthService) UPDATELoginTimestamp(username, timestamp string) error {
	return s.stor.UPDATELoginTimestamp(username, timestamp)
}

func (s *AuthService) UPDATEUserPassword(namepass vh.NamePass) error {
	return s.stor.UPDATEUserPassword(namepass)
}
