package response

type Status int

const (
	StatusOK                  Status = 200
	StatusCreated             Status = 201
	StatusBadRequest          Status = 400
	StatusUnauthorized        Status = 401
	StatusForbidden           Status = 403
	StatusNotFound            Status = 404
	StatusInternalServerError Status = 500
)
