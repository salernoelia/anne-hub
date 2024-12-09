package models

// Define a struct to represent the custom headers
type WSRequestHeaders struct {
    XUserID    string `json:"X-User-ID"`
    XDeviceID  string `json:"X-Device-ID"`
    XLanguage  string `json:"X-Language"`
}