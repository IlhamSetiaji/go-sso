package web

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/middleware"
	request "app/go-sso/internal/http/request/web/permission"
	appUsecase "app/go-sso/internal/usecase/application"
	usecase "app/go-sso/internal/usecase/permission"
	"app/go-sso/views"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PermissionHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type PermissionHandlerInterface interface {
	Index(ctx *gin.Context)
	StorePermission(ctx *gin.Context)
	UpdatePermission(ctx *gin.Context)
	DeletePermission(ctx *gin.Context)
}

func PermissionHandlerFactory(log *logrus.Logger, validator *validator.Validate) PermissionHandlerInterface {
	return &PermissionHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *PermissionHandler) Index(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-permission")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	factory := usecase.GetAllPermissionsUseCaseFactory(h.Log)

	resp, err := factory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	appFactory := appUsecase.GetAllApplicationsUseCaseFactory(h.Log)

	appResp, err := appFactory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	index := views.NewView("base", "views/permissions/index.html")
	data := map[string]interface{}{
		"Title":        "Julong Portal | Permissions",
		"Permissions":  resp.Permissions,
		"Applications": appResp.Applications,
	}

	index.Render(ctx, data)
}

func (h *PermissionHandler) StorePermission(ctx *gin.Context) {
	middleware.PermissionMiddleware("create-permission")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	payload := new(request.CreatePermissionRequest)
	if err := ctx.ShouldBind(payload); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	err := h.Validate.Struct(payload)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	var permission = &entity.Permission{
		Name:          payload.Name,
		ApplicationID: uuid.MustParse(payload.ApplicationID),
		GuardName:     payload.GuardName,
		Label:         payload.Label,
	}

	factory := usecase.StorePermissionUseCaseFactory(h.Log)
	res, err := factory.Execute(&usecase.IStorePermissionUseCaseRequest{
		Permission:    permission,
		ApplicationID: permission.ApplicationID,
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("permission created: %v", res)
	session.Set("success", "Permission created successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *PermissionHandler) UpdatePermission(ctx *gin.Context) {
	middleware.PermissionMiddleware("update-permission")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	payload := new(request.UpdatePermissionRequest)
	if err := ctx.ShouldBind(payload); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	err := h.Validate.Struct(payload)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	var permission = &entity.Permission{
		ID:            uuid.MustParse(payload.ID),
		Name:          payload.Name,
		ApplicationID: uuid.MustParse(payload.ApplicationID),
		GuardName:     payload.GuardName,
		Label:         payload.Label,
	}

	factory := usecase.UpdatePermissionUseCaseFactory(h.Log)
	res, err := factory.Execute(&usecase.IUpdatePermissionUseCaseRequest{
		ID:            permission.ID,
		Permission:    permission,
		ApplicationID: permission.ApplicationID,
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("permission updated: %v", res)
	session.Set("success", "Permission updated successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *PermissionHandler) DeletePermission(ctx *gin.Context) {
	middleware.PermissionMiddleware("delete-permission")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	payload := new(request.DeletePermissionRequest)
	if err := ctx.ShouldBind(payload); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	err := h.Validate.Struct(payload)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	factory := usecase.DeletePermissionUseCaseFactory(h.Log)
	err = factory.Execute(&usecase.IDeletePermissionUseCaseRequest{
		ID: uuid.MustParse(payload.ID),
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("permission deleted: %v", payload.ID)
	session.Set("success", "Permission deleted successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}
