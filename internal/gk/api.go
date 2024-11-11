package gk

import (
	"context"

	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

// AuthParamName name of param context
const AuthParamName = "UUID"

// API interface of API
type API interface {
	Init()
	GetRouter() any
	Signup(context.Context, *SignupRequest) (string, error)
	Signin(context.Context, *SigninRequest) (string, error)
	ListEntities(context.Context, *model.User) (*model.Entities, error)
}

// SignupRequest request of signup
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SigninRequest request of signin
type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ListEntitiesResponse list entities response
type ListEntitiesResponse struct {
	Items *model.Entities `json:"items"`
	Count int             `json:"count"`
}

// UploadEntityRequest request to upload entity
type UploadEntityRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        model.EntityType `json:"type"`
	Attr        model.MetaData   `json:"attr"`
}
