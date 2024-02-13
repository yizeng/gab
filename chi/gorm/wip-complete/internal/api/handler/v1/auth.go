package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/api/handler/v1/request"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/config"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/domain"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/pkg/jwthelper"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/service"
)

type AuthService interface {
	Signup(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
}

type AuthHandler struct {
	conf *config.APIConfig
	svc  AuthService
}

func NewAuthHandler(conf *config.APIConfig, svc AuthService) *AuthHandler {
	return &AuthHandler{
		conf: conf,
		svc:  svc,
	}
}

// HandleSignup godoc
// @Summary      Signup a new user
// @Tags         auth
// @Produce      json
// @Param        request   body      request.SignupRequest true "request body"
// @Success      201      {object}   domain.User
// @Failure      400      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /auth/signup [post]
func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	req := request.SignupRequest{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, response.ErrBadRequest(err))

		return
	}

	user, err := h.svc.Signup(r.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserEmailExists) {
			_ = render.Render(w, r, response.ErrBadRequest(service.ErrUserEmailExists))

			return
		}

		err = fmt.Errorf("v1.HandleSignup -> h.svc.Signup -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// HandleLogin godoc
// @Summary      Login a user
// @Tags         auth
// @Produce      json
// @Param        request   body      request.LoginRequest true "request body"
// @Success      200      {object}   domain.User
// @Failure      401      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /auth/login [post]
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	req := request.LoginRequest{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, response.ErrBadRequest(err))

		return
	}

	user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrWrongPassword) {
			_ = render.Render(w, r, response.ErrWrongCredentials(err))

			return
		}

		err = fmt.Errorf("v1.HandleSignup -> h.svc.Login -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	token, err := jwthelper.GenerateToken([]byte(h.conf.JWTSigningKey), user.ID, r.UserAgent())
	if err != nil {
		err = fmt.Errorf("v1.HandleSignup -> middleware.GenerateToken() -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.LoginResponse{
		Token: token,
		User:  user,
	})
}
