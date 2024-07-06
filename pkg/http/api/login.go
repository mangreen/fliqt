package api

import "fliqt/pkg/model"

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	*model.User
	Token string `json:"token" binding:"required"`
}
