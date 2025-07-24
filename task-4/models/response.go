package models

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TokenResponse JWT令牌响应
type TokenResponse struct {
	Token string `json:"token"`
} 