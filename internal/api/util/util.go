package util

import (
	"fmt"
	"strings"
)

// ParamPage 排序
//
//	{
//	   "page_num": 1,
//	   "page_size": 10,
//	   "sort": {
//	       "field": "created_at",
//	       "direction": "desc"
//	   }
//	}
type ParamPage struct {
	PageNum  int         `form:"page_num" json:"page_num" binding:"omitempty,min=1"`
	PageSize int         `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=200"` // 前端可选 15 30 50 100 200
	Sort     *SortOption `form:"sort" collection_format:"csv"`
}

// SortOption 排序选项
type SortOption struct {
	Field     string `json:"field"`     // 排序字段
	Direction string `json:"direction"` // 排序方向：asc/desc
}

func NewUtilManager() *ParamPage {
	return &ParamPage{}
}

// NormalizePagination 规范化分页参数
func (p *ParamPage) NormalizePagination(params *ParamPage) (int, int) {
	// 默认分页参数
	if params.PageNum <= 0 {
		params.PageNum = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 5
	}

	// 计算正确的偏移量
	offset := (params.PageNum - 1) * params.PageSize

	return params.PageSize, offset
}

// CalculateTotalPages 计算总页数
func (p *ParamPage) CalculateTotalPages(total int64, pageSize int) int {
	return (int(total) + pageSize - 1) / pageSize
}

func (p *ParamPage) GetSortSqlDemo(mapping map[string]string) string {
	if p.Sort == nil {
		return ""
	}

	// 检查字段是否在允许的映射中
	field, ok := mapping[p.Sort.Field]
	if !ok {
		return ""
	}

	// 确定排序方向
	direction := "ASC"
	if strings.ToLower(p.Sort.Direction) == "desc" {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", field, direction)
}

type APIResponse struct {
	Data any `json:"data,omitempty"` // 成功时返回业务数据
	Meta any `json:"meta,omitempty"` // 列表分页、游标等响应元信息
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`              // 机器可读错误码，例如 invalid_argument
	Message string `json:"message"`           // 人类可读错误描述
	Details any    `json:"details,omitempty"` // 可选的错误上下文
}

func ResponseSuccessful(data any) APIResponse {
	return Success(data)
}

func ResponseFailure(message string, err any) ErrorResponse {
	return Failure("request_failed", message, err)
}

func Success(data any) APIResponse {
	return APIResponse{
		Data: data,
	}
}

func SuccessWithMeta(data any, meta any) APIResponse {
	return APIResponse{
		Data: data,
		Meta: meta,
	}
}

func Failure(code string, message string, details any) ErrorResponse {
	if code == "" {
		code = "internal_error"
	}
	if message == "" {
		message = "request failed"
	}
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: normalizeError(details),
		},
	}
}

func normalizeError(err any) any {
	if err == nil {
		return nil
	}
	if e, ok := err.(error); ok {
		return e.Error()
	}
	return err
}
