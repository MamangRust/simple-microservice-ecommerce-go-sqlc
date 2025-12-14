package errorhandler

type ErrorHandler struct {
	UserQueryError   UserQueryErrorHandler
	UserCommandError UserCommandErrorHandler
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}
