package v1

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
	"github.com/GearFramework/GophKeeper/internal/pkg/user"
	"github.com/gin-gonic/gin"
)

func (a *RestAPI) addEntity(gr *gin.RouterGroup) {
	gr.GET("/entities", func(ctx *gin.Context) {
		a.ListEntityHandler(ctx)
	})
	gr.POST("/entities", func(ctx *gin.Context) {
		a.CreateEntityHandler(ctx)
	})
	gr.PUT("/entities/:guid", func(ctx *gin.Context) {
		a.UploadEntityHandler(ctx)
	})
	gr.GET("/entities/:guid", func(ctx *gin.Context) {
		a.ViewEntityHandler(ctx)
	})
	gr.GET("/entities/download/:guid", func(ctx *gin.Context) {
		a.DownloadEntityHandler(ctx)
	})
	gr.DELETE("/entities/:guid", func(ctx *gin.Context) {
		a.DeleteEntityHandler(ctx)
	})
}

// NewListEntitiesResponse return response struct on list entities request
func NewListEntitiesResponse(entities *model.Entities) *gk.ListEntitiesResponse {
	return &gk.ListEntitiesResponse{
		Items: entities,
		Count: len(*entities),
	}
}

// NewCreateEntityRequest return upload entity request
func NewCreateEntityRequest(ctx *gin.Context) (*gk.UploadEntityRequest, error) {
	if !strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
		logger.Log.Errorf(
			"invalid request header: Content-Type %s\n",
			ctx.Request.Header.Get("Content-Type"),
		)
		return nil, ErrUploadBadRequest
	}
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()
	var req gk.UploadEntityRequest
	err := json.NewDecoder(ctx.Request.Body).Decode(&req)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, ErrUploadBadRequest
	}
	return &req, nil
}

func (a *RestAPI) hasAuthUser(ctx *gin.Context) (*model.User, error) {
	UUID, exists := ctx.Get(gk.AuthParamName)
	fmt.Println(UUID, exists)
	if !exists || UUID == nil {
		return nil, ErrUnauthorized
	}
	usr, err := user.Repo.GetByUUID(ctx, UUID.(string))
	fmt.Println(usr, err)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return usr, nil
}

// ListEntityHandler handle of rest list user entities
func (a *RestAPI) ListEntityHandler(ctx *gin.Context) {
	logger.Log.Info("EntityListHandler called")
	usr, err := a.hasAuthUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	entities, err := a.ListEntities(ctx, usr)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	if len(*entities) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, NewListEntitiesResponse(entities))
}

// CreateEntityHandler handle creating user entity
func (a *RestAPI) CreateEntityHandler(ctx *gin.Context) {
	logger.Log.Info("CreateEntityHandler called")
	usr, err := a.hasAuthUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	req, err := NewCreateEntityRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	ent, err := a.CreateEntity(ctx, usr, req)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	ctx.Data(http.StatusOK, "text/plain", []byte(ent.GUID))
}

// UploadEntityHandler handle of upload file
func (a *RestAPI) UploadEntityHandler(ctx *gin.Context) {
	logger.Log.Info("UploadEntityHandler called")
	usr, err := a.hasAuthUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	guid := ctx.Param("guid")
	if guid == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ent, err := a.GetEntity(ctx, usr, guid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if ent.Type != model.BinaryData {
		ctx.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	a.router.MaxMultipartMemory = 128 << 20
	file, handler, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Log.Errorf("invalid FormFile \"file\": %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logger.Log.Infof("starting upload file %s", handler.Filename)
	err = a.UploadEntity(ctx, usr, guid, file)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	ctx.Status(http.StatusCreated)
}

// ViewEntityHandler handle of rest download user entity
func (a *RestAPI) ViewEntityHandler(ctx *gin.Context) {
	logger.Log.Info("ViewEntityHandler called")
	usr, err := a.hasAuthUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	guid := ctx.Param("guid")
	if guid == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logger.Log.Infof("get entity guid %s by user %s:%s", guid, usr.Username, usr.UUID)
	ent, err := a.GetEntity(ctx, usr, guid)
	if err != nil {
		logger.Log.Warn(err.Error())
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	ctx.JSON(http.StatusOK, ent)
}

// DownloadEntityHandler handle of rest download user entity
func (a *RestAPI) DownloadEntityHandler(ctx *gin.Context) {
	logger.Log.Info("DownloadEntityHandler called")
	usr, err := a.hasAuthUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	guid := ctx.Param("guid")
	if guid == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ent, reader, err := a.DownloadEntity(ctx, usr, guid)
	if err != nil {
		logger.Log.Warn(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if ent.Type != model.BinaryData {
		ctx.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	mimeType := mime.TypeByExtension(ent.Attr.Binary.Extension)
	logger.Log.Infof("get download file %s as mime %s", ent.Attr.Binary.OriginalFilename, mimeType)
	ctx.Writer.Header().Set("Content-Type", mimeType)
	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+ent.Attr.Binary.OriginalFilename)
	_, err = reader.(*os.File).WriteTo(ctx.Writer)
	if err != nil {
		logger.Log.Warn(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

// DeleteEntityHandler handle of rest request on deleting entity
func (a *RestAPI) DeleteEntityHandler(ctx *gin.Context) {
	logger.Log.Info("DownloadEntityOptions called")
	usr, err := a.hasAuthUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	guid := ctx.Param("guid")
	if guid == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = a.DeleteEntity(ctx, usr, guid)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	ctx.Status(http.StatusOK)
}
