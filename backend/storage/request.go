package storage

import "time"

const (
	SortDirectionAsc = iota
	SortDirectionDesc

	DefaultPageSize = 25
)

// BaseListRequest contains basic parameters for listing requests
type BaseListRequest struct {
	SortField     string `json:"sort_field"`
	SortDirection int    `json:"sort_direction"`
	Limit         int64  `json:"limit"`
	Page          int64  `json:"page"`
}

// ListUsersRequest contains supported fields for User list request
type ListUsersRequest struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	Status           int    `json:"status"`
	NotificationType int    `json:"notification_type"`
	BaseListRequest
}

// ListFeelingsRequest contains supported fields for Feeling list request
type ListFeelingsRequest struct {
	Title string `json:"title"`
	BaseListRequest
}

// ListStatusesRequest contains supported fields for Status list request
type ListStatusesRequest struct {
	UserID        int64     `json:"user_id"`
	Feelings      []string  `json:"feelings"`
	DatetimeStart time.Time `json:"datetime_start,omitempty"`
	DatetimeEnd   time.Time `json:"datetime_end,omitempty"`
	BaseListRequest
}

// FeelingsFrequencyRequest contains supported fields for frequent feelings list request
type FeelingsFrequencyRequest struct {
	UserID        int64     `json:"user_id"`
	DatetimeStart time.Time `json:"datetime_start,omitempty"`
	DatetimeEnd   time.Time `json:"datetime_end,omitempty"`
	BaseListRequest
}

// FeelingFrequencyItem represents item of response of feelings frequency response
type FeelingFrequencyItem struct {
	Feeling   Feeling `json:"feeling"`
	Frequency int64   `json:"frequency"`
}

// NewBaseListRequest returns new BaseListRequest with default parameters
func NewBaseListRequest() BaseListRequest {
	return BaseListRequest{
		SortDirection: SortDirectionDesc,
		Limit:         DefaultPageSize,
	}
}

// NewListUsersRequest returns model, which contains supported fields for User listing request
func NewListUsersRequest() ListUsersRequest {
	return ListUsersRequest{
		Status:           -1,
		NotificationType: -1,
		BaseListRequest:  NewBaseListRequest(),
	}
}

// NewListFeelingsRequest returns model, which contains supported fields for Feeling listing request
func NewListFeelingsRequest() ListFeelingsRequest {
	return ListFeelingsRequest{
		BaseListRequest: NewBaseListRequest(),
	}
}

// NewListStatusesRequest returns model, which contains supported fields for Status listing request
func NewListStatusesRequest() ListStatusesRequest {
	return ListStatusesRequest{
		BaseListRequest: NewBaseListRequest(),
	}
}

// NewFeelingsFrequencyRequest returns model, which contains supported fields for frequent feelings list request
func NewFeelingsFrequencyRequest() FeelingsFrequencyRequest {
	return FeelingsFrequencyRequest{
		BaseListRequest: NewBaseListRequest(),
	}
}
