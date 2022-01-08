// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	links "github.com/magmel48/go-web/internal/db/links"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, shortID, originalURL
func (_m *Repository) Create(ctx context.Context, shortID string, originalURL string) (*links.Link, bool, error) {
	ret := _m.Called(ctx, shortID, originalURL)

	var r0 *links.Link
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *links.Link); ok {
		r0 = rf(ctx, shortID, originalURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*links.Link)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, string, string) bool); ok {
		r1 = rf(ctx, shortID, originalURL)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, shortID, originalURL)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateBatch provides a mock function with given fields: ctx, originalURLs
func (_m *Repository) CreateBatch(ctx context.Context, originalURLs []string) ([]links.Link, error) {
	ret := _m.Called(ctx, originalURLs)

	var r0 []links.Link
	if rf, ok := ret.Get(0).(func(context.Context, []string) []links.Link); ok {
		r0 = rf(ctx, originalURLs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]links.Link)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, originalURLs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByShortID provides a mock function with given fields: ctx, shortID
func (_m *Repository) FindByShortID(ctx context.Context, shortID string) (*links.Link, error) {
	ret := _m.Called(ctx, shortID)

	var r0 *links.Link
	if rf, ok := ret.Get(0).(func(context.Context, string) *links.Link); ok {
		r0 = rf(ctx, shortID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*links.Link)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shortID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
