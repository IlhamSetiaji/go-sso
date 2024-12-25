package web

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/middleware"
	request "app/go-sso/internal/http/request/web/role"
	appUsecase "app/go-sso/internal/usecase/application"
	usecase "app/go-sso/internal/usecase/role"
	"app/go-sso/views"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RoleHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type RoleHandlerInterface interface {
	Index(ctx *gin.Context)
	AssignRoleToPermissions(ctx *gin.Context)
	ResignRoleFromPermission(ctx *gin.Context)
	StoreRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
}

func RoleHandlerFactory(log *logrus.Logger, validator *validator.Validate) RoleHandlerInterface {
	return &RoleHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *RoleHandler) Index(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-role")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	factory := usecase.GetAllRolesUseCaseFactory(h.Log)

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

	index := views.NewView("base", "views/roles/index.html")
	data := map[string]interface{}{
		"Title":        "Julong Portal | Roles",
		"Roles":        resp.Roles,
		"Applications": appResp.Applications,
	}

	index.Render(ctx, data)
}

func (h *RoleHandler) StoreRole(ctx *gin.Context) {
	middleware.PermissionMiddleware("create-role")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	payload := new(request.CreateRoleRequest)
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

	var role = &entity.Role{
		Name:      payload.Name,
		GuardName: payload.GuardName,
		Status:    payload.Status,
	}

	factory := usecase.StoreRoleUseCaseFactory(h.Log)
	res, err := factory.Execute(&usecase.IStoreRoleUseCaseRequest{
		Role:          role,
		ApplicationID: uuid.MustParse(payload.ApplicationID),
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("role created: %v", res)
	session.Set("success", "Role created successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *RoleHandler) UpdateRole(ctx *gin.Context) {
	middleware.PermissionMiddleware("update-role")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	payload := new(request.UpdateRoleRequest)
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

	factory := usecase.UpdateRoleUseCaseFactory(h.Log)
	res, err := factory.Execute(&usecase.IUpdateRoleUseCaseRequest{
		ID:            uuid.MustParse(payload.ID),
		Role:          &entity.Role{Name: payload.Name, GuardName: payload.GuardName, Status: payload.Status},
		ApplicationID: uuid.MustParse(payload.ApplicationID),
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("role updated: %v", res)
	session.Set("success", "Role updated successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *RoleHandler) AssignRoleToPermissions(ctx *gin.Context) {
	// middleware.PermissionMiddleware("assign-role-to-permissions")(ctx)
	// if ctx.IsAborted() {
	// 	ctx.Abort()
	// 	return
	// }

	session := sessions.Default(ctx)
	payload := new(usecase.IAssignRoleToPermissionIDsUsecaseRequest)
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

	factory := usecase.AssignRoleToPermissionIDsUsecaseFactory(h.Log)
	_, err = factory.Execute(&usecase.IAssignRoleToPermissionIDsUsecaseRequest{
		RoleID:        payload.RoleID,
		PermissionIDs: payload.PermissionIDs,
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("role assigned to permissions: %v", payload.RoleID)
	session.Set("success", "Role assigned to permissions successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *RoleHandler) ResignRoleFromPermission(ctx *gin.Context) {
	// middleware.PermissionMiddleware("resign-role-from-permission")(ctx)
	// if ctx.IsAborted() {
	// 	ctx.Abort()
	// 	return
	// }

	session := sessions.Default(ctx)
	payload := new(usecase.IResignRoleFromPermissionUsecaseRequest)
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

	factory := usecase.ResignRoleFromPermissionUsecaseFactory(h.Log)
	_, err = factory.Execute(&usecase.IResignRoleFromPermissionUsecaseRequest{
		RoleID:       payload.RoleID,
		PermissionID: payload.PermissionID,
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("role resigned from permission: %v", payload.RoleID)
	session.Set("success", "Role resigned from permission successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *RoleHandler) DeleteRole(ctx *gin.Context) {
	middleware.PermissionMiddleware("delete-role")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	payload := new(request.DeleteRoleRequest)
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

	factory := usecase.DeleteRoleUseCaseFactory(h.Log)
	err = factory.Execute(&usecase.IDeleteRoleUseCaseRequest{
		ID: uuid.MustParse(payload.ID),
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("role deleted: %v", payload.ID)
	session.Set("success", "Role deleted successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}
