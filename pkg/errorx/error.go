package errorx

func NewError(code ResCode) *ResponseError {
	Msg := getErrorMsg(code)
	if Msg == unknownMsg {
		return &ResponseError{
			Code:    ErrUnknown,
			Message: unknownMsg,
		}
	}
	return &ResponseError{
		Code:    code,
		Message: Msg,
	}
}

func NewErrorWithMsg(code ResCode, message string) *ResponseError {
	return &ResponseError{
		Code:    code,
		Message: message,
	}
}

func NewSuccess() *ResponseSuccess {
	return &ResponseSuccess{
		Code:    Success,
		Message: getErrorMsg(Success),
		Data:    nil,
	}
}

func (e ResponseError) Error() string {
	return e.Message
}

func getErrorMsg(code ResCode) string {
	if _, ok := resMap[code]; !ok {
		return unknownMsg
	}
	return resMap[code]
}
