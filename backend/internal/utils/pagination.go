package utils

import (
    "math"
)

// PaginationParams holds pagination request parameters
type PaginationParams struct {
    Page  int    `json:"page"`
    Limit int    `json:"limit"`
    Sort  string `json:"sort"`
    Order string `json:"order"`
}

// PaginationResponse contains paginated results metadata
type PaginationResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    Total      int64       `json:"total"`
    TotalPages int         `json:"total_pages"`
    HasNext    bool        `json:"has_next"`
    HasPrev    bool        `json:"has_prev"`
}

// ValidatePagination ensures pagination params are valid
func ValidatePagination(params *PaginationParams) *PaginationParams {
    // Set defaults if not provided
    if params.Page < 1 {
        params.Page = 1
    }
    
    if params.Limit < 1 {
        params.Limit = 10
    }
    
    // Cap maximum limit to prevent abuse
    if params.Limit > 100 {
        params.Limit = 100
    }
    
    // Default sort order
    if params.Order == "" {
        params.Order = "desc"
    }
    
    if params.Sort == "" {
        params.Sort = "created_at"
    }
    
    return params
}

// CalculateOffset computes database offset from page
func CalculateOffset(page, limit int) int {
    return (page - 1) * limit
}

// BuildPaginationResponse creates pagination metadata
func BuildPaginationResponse(data interface{}, page, limit int, total int64) *PaginationResponse {
    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    
    return &PaginationResponse{
        Data:       data,
        Page:       page,
        Limit:      limit,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }
}