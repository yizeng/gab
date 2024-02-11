package response

import "github.com/yizeng/gab/chi/gorm/wip-complete/internal/domain"

type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}
