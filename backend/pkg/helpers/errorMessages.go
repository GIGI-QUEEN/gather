package helpers

var InternalErrorMsg = ErrorMsg{
	ErrorDescription: "Internal server error",
	ErrorType:        "INTERNAL_SERVER_ERROR",
}
var NotFoundErrorMsg = ErrorMsg{
	ErrorDescription: "Page not found",
	ErrorType:        "NOT_FOUND_ERROR",
}
var UnauthorizedErrorMsg = ErrorMsg{
	ErrorDescription: "Restricted: user not authorized",
	ErrorType:        "UNAUTHORIZED_ERROR",
}

var ForbiddenErrorMsg = ErrorMsg{
	ErrorDescription: "Forbidden",
	ErrorType:        "FORBIDDEN_ERROR",
}
var EmptyCategoryErrorMsg = ErrorMsg{
	ErrorDescription: "Category is empty",
	ErrorType:        "EMPTY_CATEGORY",
}

var MethodNotAllowedMsg = ErrorMsg{
	ErrorDescription: "Method not allowed",
	ErrorType:        "METHOD_NOT_ALLOWED",
}

var WrongFileTypeErrorMsg = ErrorMsg{
	ErrorDescription: "Wrong file type",
	ErrorType:        "WRONG_FILE_TYPE",
}
