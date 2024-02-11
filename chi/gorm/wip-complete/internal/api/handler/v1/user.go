package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/domain"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/pkg/jwthelper"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/service"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (domain.User, error)
}

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// HandleGetUser godoc
// @Summary      Get a user
// @Tags         users
// @Produce      json
// @Param        userID   path       int  true "user ID"
// @Success      200      {object}   domain.User
// @Failure      401      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /users/{userID} [get]
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		_ = render.Render(w, r, response.ErrInvalidInput("userID", rawUserID))

		return
	}

	if userID <= 0 {
		_ = render.Render(w, r, response.ErrNotFound("user", "ID", userID))

		return
	}

	claims, err := jwthelper.RetrieveClaimsFromContext(r.Context())
	if err != nil {
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	if uint(userID) != claims.UserID {
		_ = render.Render(w, r, response.ErrPermissionDenied(fmt.Errorf("can't view user %v by user %v", userID, claims.UserID)))

		return
	}

	user, err := h.svc.GetUser(r.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			_ = render.Render(w, r, response.ErrNotFound("user", "ID", userID))

			return
		}

		err = fmt.Errorf("v1.HandleGetUser -> h.svc.GetUser -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	render.Status(r, http.StatusOK)
	if err = render.Render(w, r, response.NewUser(&user)); err != nil {
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}
}
