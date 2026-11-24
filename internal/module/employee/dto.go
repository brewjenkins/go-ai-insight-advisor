package employee

type CreateEmployeeRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type EmployeeResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func toEmployeeResponse(e *Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
	}
}
