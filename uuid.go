package jwtauthserver

type IDprovider interface {
	ID() (string, error)
}
