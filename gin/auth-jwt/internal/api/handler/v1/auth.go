package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yizeng/gab/gin/auth-jwt/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/auth-jwt/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/auth-jwt/internal/domain"
	"github.com/yizeng/gab/gin/auth-jwt/internal/service"
)

type AuthService interface {
	Signup(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
}

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

// HandleSignup godoc
// @Summary      Signup a new user
// @Tags         auth
// @Produce      json
// @Param        request   body      request.SignupRequest true "request body"
// @Success      201      {object}   domain.User
// @Failure      400      {object}   response.ErrResponse
// @Failure      500      {object}   response.ErrResponse
// @Router       /auth/signup [post]
func (h *AuthHandler) HandleSignup(ctx *gin.Context) {
	req := request.SignupRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err))

		return
	}

	if err := req.Validate(); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err))

		return
	}

	user, err := h.svc.Signup(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserEmailExists) {
			response.RenderError(ctx, response.NewBadRequest(service.ErrUserEmailExists))

			return
		}

		err = fmt.Errorf("v1.HandleSignup -> h.svc.Signup -> %w", err)
		response.RenderError(ctx, response.NewInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// HandleLogin godoc
// @Summary      Login a user
// @Tags         auth
// @Produce      json
// @Param        request   body      request.LoginRequest true "request body"
// @Success      200      {object}   domain.User
// @Failure      401      {object}   response.ErrResponse
// @Failure      500      {object}   response.ErrResponse
// @Router       /auth/login [post]
func (h *AuthHandler) HandleLogin(ctx *gin.Context) {
	req := request.LoginRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err))

		return
	}

	if err := req.Validate(); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err))

		return
	}

	user, err := h.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			response.RenderError(ctx, response.NewWrongCredentials(err))

			return
		}

		err = fmt.Errorf("v1.HandleSignup -> h.svc.Login -> %w", err)
		response.RenderError(ctx, response.NewInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusOK, user)
}
