package auth0

import (
	"testing"
)

func init() {
	RequireAuth0()
}

func TestGetAccessToken(t *testing.T) {
	apiClient, _ := NewAuth0APIClient()
	err := apiClient.getAccessToken()
	if err != nil {
		t.Errorf("auth0 oauth access token retrieval failed; %s", err.Error())
		return
	}
	if apiClient.Token == nil {
		t.Error("auth0 oauth access token retrieval failed; access token unset on api client")
	}
	if apiClient.TokenExpiresAt == nil {
		t.Error("auth0 oauth access token retrieval failed; access token expiration unset on api client")
	}
}

func TestExportUsers(t *testing.T) {
	users, _ := ExportUsers()
	if users == nil {
		t.Error("failed to export auth0 users")
	}
	if len(users) == 0 {
		t.Error("failed to export auth0 users")
	}
}

// func TestGetUser(t *testing.T) {
// 	user, _ := GetUser("")
// 	if user == nil {
// 		t.Error("failed to get auth0 user")
// 	}
// }
