package web_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	_assert "github.com/stretchr/testify/assert"
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

func (m *mockDbConn) Close(ctx context.Context) error {
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
	assert := _assert.New(t)

	mdc := new(mockDbConn)
	mrw := new(mockResponseWriter)

	web := web.NewWith(&web.Args{Db: mdc})

	web.CreateNote(mrw, &http.Request{})
}
