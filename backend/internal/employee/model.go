package employee

type Employee struct {
	ID         string  `json:"_id"`
	EmpID      string  `json:"empID"`
	CreatedBy  string  `json:"createdBy"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
}

type CreateEmployeeRequest struct {
	EmpID      string  `json:"empID"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
}

type UpdateEmployeeRequest struct {
	Name       *string  `json:"name,omitempty"`
	Email      *string  `json:"email,omitempty"`
	Department *string  `json:"department,omitempty"`
	Salary     *float64 `json:"salary,omitempty"`
}
