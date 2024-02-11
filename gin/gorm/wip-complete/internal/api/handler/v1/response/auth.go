package response

import "github.com/yizeng/gab/gin/wip-complete/internal/domain"

type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}
