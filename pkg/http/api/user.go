package api

type UserCreateRequest struct {
	ID         string `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
	ManagerID  string `json:"manager_id"`
}

type UserUpdateRequest struct {
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
	ManagerID  string `json:"manager_id"`
}
