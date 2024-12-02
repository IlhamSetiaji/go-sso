package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type IUserService interface {
	GetUsersZitadel(url string, accessToken string, payload map[string]interface{}) ([]map[string]interface{}, error)
}

type UserService struct {
	Log *logrus.Logger
}

func NewUserService(log *logrus.Logger) IUserService {
	return &UserService{
		Log: log,
	}
}

func UserServiceFactory(log *logrus.Logger) IUserService {
	return NewUserService(log)
}

func (s *UserService) GetUsersZitadel(url string, accessToken string, payload map[string]interface{}) ([]map[string]interface{}, error) {
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[UserService.GetUsersZitadel] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))

	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[UserService.GetUsersZitadel] Error when fetching client: " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[UserService.GetUsersZitadel] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[UserService.GetUsersZitadel] Error when fetching response: " + string(bodyBytes))
	}
	var response struct {
		Users []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[UserService.GetUsersZitadel] Error when decoding response: " + err.Error())
	}

	return response.Users, nil
}
