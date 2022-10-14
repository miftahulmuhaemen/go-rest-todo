package common

type ResponseCode struct {
	RC       string `json:"rc"`
	Message  string `json:"message"`
	Messages string `json:"messages,omitempty"`
}

func GetSuccessMessage() ResponseCode {
	response := ResponseCode{
		RC:      "00",
		Message: "success",
	}
	return response
}

func GetErrorMessage(rc string, msg string) ResponseCode {
	var (
		response ResponseCode
	)

	switch rc {
	default:
		response = ResponseCode{
			RC:      rc,
			Message: "bad request",
		}

	case "51A":
		response = ResponseCode{
			RC:       rc,
			Message:  "missing authentication token",
			Messages: msg,
		}
	case "51B":
		response = ResponseCode{
			RC:       rc,
			Message:  "failed authentication",
			Messages: msg,
		}
	case "51C":
		response = ResponseCode{
			RC:       rc,
			Message:  "username already exist",
			Messages: msg,
		}
	case "51D":
		response = ResponseCode{
			RC:       rc,
			Message:  "invalid or expired token",
			Messages: msg,
		}
	case "52B":
		response = ResponseCode{
			RC:       rc,
			Message:  "route or data not found or unavailable",
			Messages: msg,
		}
	case "52C":
		response = ResponseCode{
			RC:       rc,
			Message:  "request validation mismatch",
			Messages: msg,
		}
	case "53S":
		response = ResponseCode{
			RC:       rc,
			Message:  "something went wrong",
			Messages: msg,
		}
	}
	return response
}
