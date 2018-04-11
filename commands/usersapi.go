package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/makpoc/hades-api/sheet/models"
	"github.com/makpoc/hdsbrute"
)

type UserGetter interface {
	GetUsers() (models.Users, error)
	GetUser(string) (models.User, error)
}

type UserAPI struct {
	backendURL    string
	backendSecret string
}

func NewUserApi(backendURL, backendSecret string) UserAPI {
	return UserAPI{backendURL, backendSecret}
}

func (u UserAPI) getFullURL() string {
	url := fmt.Sprintf("%s/api/v1/users", u.backendURL)
	if u.backendSecret != "" {
		url = fmt.Sprintf("%s?secret=%s", url, u.backendSecret)
	}

	return url
}

func (u UserAPI) GetUsers() (models.Users, error) {
	var result models.Users

	users, err := getSheetUsers(u.getFullURL())
	if err != nil {
		log.Printf("Failed to get Sheet Users: %v", err)

		return result, fmt.Errorf("failed to get sheet users", err)
	}

	return users, nil
}

func (u UserAPI) GetUser(userName string) (*models.User, error) {
	users, err := u.GetUsers()
	if err != nil {
		return nil, err
	}

	userName = strings.TrimSpace(userName)

	fmt.Printf("Seaching for User: %s\n", userName)
	for _, user := range users {
		fmt.Printf("Checking user %v\n", user)
		if strings.ToLower(userName) == strings.ToLower(user.Name) || hdsbrute.TrimMentionPrefix(userName) == user.DiscordID {
			return &user, nil
		}
	}

	for _, user := range users {
		if matchPartialUser(user.Name, userName) {
			return &user, nil
		}
	}
	return &models.User{}, fmt.Errorf("user %s not found in sheet database", userName)
}

func matchPartialUser(given string, wanted string) bool {
	return strings.Contains(strings.TrimSpace(strings.ToLower(given)), strings.ToLower(wanted))
}

func getSheetUsers(url string) (models.Users, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info. Status was: %s", resp.Status)
	}

	var users models.Users
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return users, nil
}
