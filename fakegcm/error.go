package fakegcm

import "net/http"

type Error struct {
	Code  int
	Error string
}

var (
	AuthenticationError       = Error{http.StatusUnauthorized, ""}
	DeviceMessageRateExceeded = Error{http.StatusOK, "DeviceMessageRateExceeded"}
	InternalServerError       = Error{http.StatusOK, "InternalServerError"}
	InvalidApnsCredential     = Error{http.StatusOK, "InvalidApnsCredential"}
	InvalidDataKey            = Error{http.StatusOK, "InvalidDataKey"}
	InvalidJSON               = Error{http.StatusBadRequest, ""}
	InvalidPackageName        = Error{http.StatusOK, "InvalidPackageName"}
	InvalidParameters         = Error{http.StatusBadRequest, "InvalidParameters"}
	InvalidRegistrationToken  = Error{http.StatusOK, "InvalidRegistration"}
	InvalidTimeToLive         = Error{http.StatusOK, "InvalidTtl"}
	MessageTooBig             = Error{http.StatusOK, "MessageTooBig"}
	MismatchedSender          = Error{http.StatusOK, "MismatchSenderId"}
	MissingRegistrationToken  = Error{http.StatusOK, "MissingRegistration"}
	Timeout                   = Error{http.StatusOK, "Unavailable"}
	TopicsMessageRateExceeded = Error{http.StatusOK, "TopicsMessageRateExceeded"}
	UnregisteredDevice        = Error{http.StatusOK, "NotRegistered"}
)
