package models

type WSRequestHeaders struct {
	XUserID   string `json:"X-User-ID"`
	XDeviceID string `json:"X-Device-ID"`
	XLanguage string `json:"X-Language"`
}
