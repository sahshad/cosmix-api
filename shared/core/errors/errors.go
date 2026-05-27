package errors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Field   string `json:"field,omitempty"`
	Err     error  `json:"-"`
}

func (err *AppError) Error() string {
	return err.Message
}