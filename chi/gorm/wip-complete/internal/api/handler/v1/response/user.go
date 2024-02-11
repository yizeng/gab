package response

import (
	"net/http"

	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/domain"
)

type UserResponse struct {
	*domain.User
}

func NewUser(user *domain.User) *UserResponse {
	resp := &UserResponse{
		User: user,
	}

	return resp
}

func (resp *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
