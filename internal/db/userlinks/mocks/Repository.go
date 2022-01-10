// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	userlinks "github.com/magmel48/go-web/internal/db/userlinks"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userID, linkID
func (_m *Repository) Create(ctx context.Context, userID *string, linkID int) error {
	ret := _m.Called(ctx, userID, linkID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *string, int) error); ok {
		r0 = rf(ctx, userID, linkID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByLinkID provides a mock function with given fields: ctx, userID, linkID
func (_m *Repository) FindByLinkID(ctx context.Context, userID *string, linkID int) (*userlinks.UserLink, error) {
	ret := _m.Called(ctx, userID, linkID)

	var r0 *userlinks.UserLink
	if rf, ok := ret.Get(0).(func(context.Context, *string, int) *userlinks.UserLink); ok {
		r0 = rf(ctx, userID, linkID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userlinks.UserLink)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *string, int) error); ok {
		r1 = rf(ctx, userID, linkID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, userID
func (_m *Repository) List(ctx context.Context, userID *string) ([]userlinks.UserLink, error) {
	ret := _m.Called(ctx, userID)

	var r0 []userlinks.UserLink
	if rf, ok := ret.Get(0).(func(context.Context, *string) []userlinks.UserLink); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]userlinks.UserLink)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}