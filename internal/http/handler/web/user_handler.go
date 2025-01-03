package web

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/middleware"
	userRequest "app/go-sso/internal/http/request/web/user"
	empUsecase "app/go-sso/internal/usecase/employee"
	roleUsecase "app/go-sso/internal/usecase/role"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/views"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type UserHandlerInterface interface {
	Index(ctx *gin.Context)
	StoreUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	// UserRoles(ctx *gin.Context)
}

func UserHandlerFactory(log *logrus.Logger, validator *validator.Validate) UserHandlerInterface {
	return &UserHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *UserHandler) Index(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-user")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	// page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	// pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// factory := usecase.FindAllPaginatedUseCaseFactory(h.Log)

	// req := &usecase.IFindAllPaginatedRequest{
	// 	Page:     page,
	// 	PageSize: pageSize,
	// }
	// resp, err := factory.Execute(req)

	factory := usecase.GetAllUsersUseCaseFactory(h.Log)

	resp, err := factory.Execute()
	if err != nil {
		h.Log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roleFactory := roleUsecase.GetAllRolesUseCaseFactory(h.Log)
	role, err := roleFactory.Execute()
	if err != nil {
		h.Log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	empFactory := empUsecase.FindAllEmployeesNotInUserUseCaseFactory(h.Log)
	empResp, err := empFactory.Execute()
	if err != nil {
		h.Log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	index := views.NewView("base", "views/users/index.html")
	data := map[string]interface{}{
		"Title":     "Julong Portal | Users",
		"Users":     resp.Users,
		"Roles":     role.Roles,
		"Employees": empResp.Employees,
	}

	index.Render(ctx, data)
}

func (h *UserHandler) StoreUser(ctx *gin.Context) {
	middleware.PermissionMiddleware("create-user")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}
	session := sessions.Default(ctx)
	payload := new(userRequest.CreateUserRequest)
	if err := ctx.ShouldBind(payload); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("payload: %v", payload)

	err := h.Validate.Struct(payload)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	var employeeID *uuid.UUID

	if payload.EmployeeID != "" {
		parsedID := uuid.MustParse(payload.EmployeeID)
		employeeID = &parsedID
	} else {
		employeeID = nil
	}

	var user = &entity.User{
		Name:        payload.Name,
		Username:    payload.Username,
		Email:       payload.Email,
		Gender:      payload.Gender,
		MobilePhone: payload.MobilePhone,
		Password:    string(hashedPasswordBytes),
		Status:      payload.Status,
		EmployeeID:  employeeID,
	}

	factory := usecase.CreateUserUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.ICreateUserUseCaseRequest{
		User:    user,
		RoleIDs: payload.RoleIDs,
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("user created: %v", response)
	session.Set("success", "User created successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	middleware.PermissionMiddleware("update-user")(ctx)
	session := sessions.Default(ctx)
	payload := new(userRequest.UpdateUserRequest)
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

	var employeeID *uuid.UUID

	if payload.EmployeeID != "" {
		parsedID := uuid.MustParse(payload.EmployeeID)
		employeeID = &parsedID
	} else {
		employeeID = nil
	}

	var user = &entity.User{
		ID:          uuid.MustParse(payload.ID),
		Name:        payload.Name,
		Username:    payload.Username,
		Email:       payload.Email,
		Gender:      payload.Gender,
		MobilePhone: payload.MobilePhone,
		Status:      payload.Status,
		EmployeeID:  employeeID,
	}
	factory := usecase.UpdateUserUseCaseFactory(h.Log)

	response, err := factory.Execute(usecase.IUpdateUserUseCaseRequest{
		User:    user,
		RoleIDs: payload.RoleIDs,
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("user updated: %v", response)
	session.Set("success", "User updated successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	middleware.PermissionMiddleware("delete-user")(ctx)
	session := sessions.Default(ctx)
	payload := new(userRequest.DeleteUserRequest)
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

	factory := usecase.DeleteUserUseCaseFactory(h.Log)

	err = factory.Execute(usecase.IDeleteUserUseCaseRequest{
		ID: uuid.MustParse(payload.ID),
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	h.Log.Printf("user deleted")
	session.Set("success", "User deleted successfully")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}
