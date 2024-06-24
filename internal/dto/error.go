package dto

type CustomError struct {
	statusCode int
	errorStr   string
	request    any
}

type ErrInfo struct {
	ErrorMessage string
	StatusCode   int
}

func (c *CustomError) Error() string {
	return c.errorStr
}

func NewCustomError(errorInfo ErrInfo, req any) *CustomError {
	return &CustomError{
		errorStr:   errorInfo.ErrorMessage,
		statusCode: errorInfo.StatusCode,
		request:    req,
	}
}

func (c *CustomError) GetStatus() int {
	return c.statusCode
}

func (c *CustomError) GetRequestInfo() any {
	return c.request
}
