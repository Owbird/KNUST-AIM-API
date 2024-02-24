package models

type UserAuthPayload struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	StudentId string `json:"studentId"`
}

type UserCookies struct {
	Antiforgery string `json:"antiforgery"`
	Session     string `json:"session"`
	Identity    string `json:"identity"`
}

type UserResponse struct {
	Message string      `json:"message"`
	Cookies UserCookies `json:"cookies"`
}
