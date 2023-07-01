package share

type Status struct {
	Code    int32
	Message string
}

func statusMessage(errorCode int32) string {
	switch errorCode {
	case 200:
		return "success"
	case 400:
		return "bad request"
	case 401:
		return "unauthorized"
	case 404:
		return "not found"
	default:
		return "UnKnownError"
	}
}

func GenerateStatus(errorCode int32, Error string) Status {
	stt := Status{
		Code:    errorCode,
		Message: Error + " " + statusMessage(errorCode),
	}
	return stt
}
