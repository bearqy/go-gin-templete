package util

import (
	"errors"
	"testing"
)

func TestNormalizePagination(t *testing.T) {
	params := &ParamPage{}
	limit, offset := params.NormalizePagination(params)

	if limit != 5 || offset != 0 {
		t.Fatalf("expected default limit 5 offset 0, got %d %d", limit, offset)
	}
}

func TestGetSortSQLDemoUsesWhitelist(t *testing.T) {
	params := &ParamPage{Sort: &SortOption{Field: "created_at", Direction: "desc"}}
	orderBy := params.GetSortSqlDemo(map[string]string{"created_at": "created_at"})

	if orderBy != "created_at DESC" {
		t.Fatalf("unexpected order by: %q", orderBy)
	}

	params.Sort.Field = "unsafe_field"
	if orderBy := params.GetSortSqlDemo(map[string]string{"created_at": "created_at"}); orderBy != "" {
		t.Fatalf("expected unsafe field to be ignored, got %q", orderBy)
	}
}

func TestResponseSuccessful(t *testing.T) {
	response := ResponseSuccessful("hello")

	if response.Data != "hello" {
		t.Fatalf("expected data, got %v", response.Data)
	}
	if response.Meta != nil {
		t.Fatalf("expected nil meta, got %v", response.Meta)
	}
}

func TestSuccessWithMeta(t *testing.T) {
	meta := map[string]any{"page": 1, "page_size": 20, "total": 100}
	response := SuccessWithMeta([]string{"hello"}, meta)

	if response.Data == nil {
		t.Fatal("expected data")
	}
	if response.Meta == nil {
		t.Fatal("expected meta")
	}
}

func TestResponseFailure(t *testing.T) {
	response := ResponseFailure("", errors.New("invalid request"))

	if response.Error.Code != "request_failed" {
		t.Fatalf("expected default error code, got %q", response.Error.Code)
	}
	if response.Error.Message != "request failed" {
		t.Fatalf("expected default failure message, got %q", response.Error.Message)
	}
	if response.Error.Details != "invalid request" {
		t.Fatalf("expected normalized error details, got %v", response.Error.Details)
	}
}
