package oauthtelos

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/mattermost/mattermost-server/v6/einterfaces"
	"github.com/mattermost/mattermost-server/v6/model"
)

type TelosProvider struct {
}

type TelosUser struct {
	Sub               string `json:"sub"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

func init() {
	provider := &TelosProvider{}
	einterfaces.RegisterOAuthProvider(model.UserAuthServiceTelos, provider)
}

func userFromTelosUser(tu *TelosUser) *model.User {
	user := &model.User{}
	username := tu.PreferredUsername
	if username == "" {
		username = tu.Email
	}
	user.Username = model.CleanUsername(username)
	user.FirstName = tu.GivenName
	user.LastName = tu.FamilyName
	user.Email = tu.Email
	user.Email = strings.ToLower(user.Email)
	user.AuthData = &tu.Sub
	user.AuthService = model.UserAuthServiceTelos

	return user
}

/*
{
    "sub":"0a8ff279-c29d-40b3-a74e-6f9dee3c6d6f",
    "email_verified":false,
    "name":"Clay Gulick",
    "preferred_username":"clay@teloshs.com",
    "given_name":"Clay",
    "family_name":"Gulick",
    "email":"clay@teloshs.com"
}
*/
func telosUserFromJSON(data io.Reader) (*TelosUser, error) {
	decoder := json.NewDecoder(data)
	var tu TelosUser
	err := decoder.Decode(&tu)
	if err != nil {
		return nil, err
	}
	return &tu, nil
}

func (tu *TelosUser) IsValid() error {

	if tu.Email == "" {
		return errors.New("user e-mail should not be empty")
	}

	return nil
}

func (m *TelosProvider) GetUserFromJSON(data io.Reader, tokenUser *model.User) (*model.User, error) {
	tu, err := telosUserFromJSON(data)
	if err != nil {
		return nil, err
	}
	if err = tu.IsValid(); err != nil {
		return nil, err
	}

	return userFromTelosUser(tu), nil
}

func (m *TelosProvider) GetSSOSettings(config *model.Config, service string) (*model.SSOSettings, error) {
	return &config.TelosLoginSettings, nil
}

func (m *TelosProvider) GetUserFromIdToken(idToken string) (*model.User, error) {
	return nil, nil
}

func (m *TelosProvider) IsSameUser(dbUser, oauthUser *model.User) bool {
	return dbUser.AuthData == oauthUser.AuthData
}
