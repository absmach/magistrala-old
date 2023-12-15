// Code generated by mockery v2.38.0. DO NOT EDIT.

// Copyright (c) Abstract Machines

package postgres

import (
	context "context"

	clients "github.com/absmach/magistrala/pkg/clients"

	mock "github.com/stretchr/testify/mock"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// ChangeStatus provides a mock function with given fields: ctx, client
func (_m *MockRepository) ChangeStatus(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for ChangeStatus")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckSuperAdmin provides a mock function with given fields: ctx, adminID
func (_m *MockRepository) CheckSuperAdmin(ctx context.Context, adminID string) error {
	ret := _m.Called(ctx, adminID)

	if len(ret) == 0 {
		panic("no return value specified for CheckSuperAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, adminID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveAll provides a mock function with given fields: ctx, pm
func (_m *MockRepository) RetrieveAll(ctx context.Context, pm clients.Page) (clients.ClientsPage, error) {
	ret := _m.Called(ctx, pm)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAll")
	}

	var r0 clients.ClientsPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Page) (clients.ClientsPage, error)); ok {
		return rf(ctx, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Page) clients.ClientsPage); ok {
		r0 = rf(ctx, pm)
	} else {
		r0 = ret.Get(0).(clients.ClientsPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Page) error); ok {
		r1 = rf(ctx, pm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveAllBasicInfo provides a mock function with given fields: ctx, pm
func (_m *MockRepository) RetrieveAllBasicInfo(ctx context.Context, pm clients.Page) (clients.ClientsPage, error) {
	ret := _m.Called(ctx, pm)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAllBasicInfo")
	}

	var r0 clients.ClientsPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Page) (clients.ClientsPage, error)); ok {
		return rf(ctx, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Page) clients.ClientsPage); ok {
		r0 = rf(ctx, pm)
	} else {
		r0 = ret.Get(0).(clients.ClientsPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Page) error); ok {
		r1 = rf(ctx, pm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveAllByIDs provides a mock function with given fields: ctx, pm
func (_m *MockRepository) RetrieveAllByIDs(ctx context.Context, pm clients.Page) (clients.ClientsPage, error) {
	ret := _m.Called(ctx, pm)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAllByIDs")
	}

	var r0 clients.ClientsPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Page) (clients.ClientsPage, error)); ok {
		return rf(ctx, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Page) clients.ClientsPage); ok {
		r0 = rf(ctx, pm)
	} else {
		r0 = ret.Get(0).(clients.ClientsPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Page) error); ok {
		r1 = rf(ctx, pm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveByID provides a mock function with given fields: ctx, id
func (_m *MockRepository) RetrieveByID(ctx context.Context, id string) (clients.Client, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveByID")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (clients.Client, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) clients.Client); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveByIdentity provides a mock function with given fields: ctx, identity
func (_m *MockRepository) RetrieveByIdentity(ctx context.Context, identity string) (clients.Client, error) {
	ret := _m.Called(ctx, identity)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveByIdentity")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (clients.Client, error)); ok {
		return rf(ctx, identity)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) clients.Client); ok {
		r0 = rf(ctx, identity)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, identity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, client
func (_m *MockRepository) Save(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, client
func (_m *MockRepository) Update(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateIdentity provides a mock function with given fields: ctx, client
func (_m *MockRepository) UpdateIdentity(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for UpdateIdentity")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOwner provides a mock function with given fields: ctx, client
func (_m *MockRepository) UpdateOwner(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOwner")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRole provides a mock function with given fields: ctx, client
func (_m *MockRepository) UpdateRole(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRole")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSecret provides a mock function with given fields: ctx, client
func (_m *MockRepository) UpdateSecret(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSecret")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTags provides a mock function with given fields: ctx, client
func (_m *MockRepository) UpdateTags(ctx context.Context, client clients.Client) (clients.Client, error) {
	ret := _m.Called(ctx, client)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTags")
	}

	var r0 clients.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) (clients.Client, error)); ok {
		return rf(ctx, client)
	}
	if rf, ok := ret.Get(0).(func(context.Context, clients.Client) clients.Client); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Get(0).(clients.Client)
	}

	if rf, ok := ret.Get(1).(func(context.Context, clients.Client) error); ok {
		r1 = rf(ctx, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
