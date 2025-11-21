package dto

// UpdateUserRoleRequest represents request to update user role
type UpdateUserRoleRequest struct {
	RoleName string `json:"role_name" validate:"required,oneof=admin teacher student"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users      []UserDTO `json:"users"`
	Total      int64     `json:"total"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}
