package web_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/velvetreactor/metrics/pkg/web"
)

type mockDbConn struct{ mock.Mock }

func (m *mockDbConn) Exec(
	_ context.Context,
	query string,
	values ...any,
) (pgconn.CommandTag, error) {
	args := m.Called(query, values)

	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *mockDbConn) Close(_ context.Context) error {
	args := m.Called()

	return args.Error(0)
}

type mockResponseWriter struct{ mock.Mock }

func (m *mockResponseWriter) Header() http.Header {
	args := m.Called()

	return args.Get(0).(http.Header)
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	args := m.Called(data)

	return args.Get(0).(int), args.Error(1)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

func TestCreateNote(t *testing.T) {
	//assert := _assert.New(t)

	mdc := new(mockDbConn)
	mrw := new(mockResponseWriter)

	req := &http.Request{
		Form: url.Values{
			"body": []string{
				"asdfasdf",
			},
		},
		PostForm: url.Values{
			"body": []string{
				"asdfasdf",
			},
		},
	}

	mdc.On(
		"Exec",
		"INSERT INTO notes (body) VALUES ($1)",
		[]interface{}{"asdfasdf"},
	).Return(pgconn.NewCommandTag("INSERT 0 1"), nil)
	mdc.On("Close").Return(nil)

	web := web.NewWith(&web.Args{Db: mdc})
	web.CreateNote(mrw, req)

	mdc.AssertExpectations(t)
	mrw.AssertExpectations(t)
}
