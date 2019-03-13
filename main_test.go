package main_test

import (
	"."
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var testInstance main.App

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testInstance.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestMain(m *testing.M) {
	testInstance = main.App{}
	testInstance.Initialize()
	code := m.Run()
	os.Exit(code)
}

func TestLeaderboard(t *testing.T) {
	req, _ := http.NewRequest("GET", "/leaderboard?offset=2", nil)
	response := executeRequest(req)
	decoder := json.NewDecoder(response.Body)
	var sortedUsers []main.User
	err := decoder.Decode(&sortedUsers)
	if err != nil {
		t.Errorf("Trouble with decoding: %s", err)
	}

	if len(sortedUsers) != 2 {
		t.Errorf("Wrong len. Expected 2, got %d", len(sortedUsers))
	}
	if sortedUsers[0].Points < sortedUsers[1].Points {
		t.Error("Expexted sortes list of users ( 2 users )")
	}
	checkResponseCode(t, http.StatusOK, response.Code)
	main.Users = []main.User{}
}

func TestSignup(t *testing.T) {
	reqGet, _ := http.NewRequest("GET", "/signup", nil)
	reqPost, _ := http.NewRequest("POST", "/signup", strings.NewReader(`{"nickname":"tested","email":"tested@gmail.com","password":"qqq"}`))
	responsePost := executeRequest(reqPost)
	responseGet := executeRequest(reqGet)
	checkResponseCode(t, http.StatusOK, responsePost.Code)
	checkResponseCode(t, http.StatusNotFound, responseGet.Code)
	usersCreated := false
	for _, user := range main.Users {
		if user.Nickname == "tested" && user.Email == "tested@gmail.com" {
			usersCreated = true
		}
	}
	if !usersCreated {
		t.Error("New users not created")
	}
}