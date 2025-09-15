package tokenizers

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) GenerateToken(username string, userID int64, userType string) (string, error) {
	args := m.Called(username, userID, userType)
	return args.String(0), args.Error(1)
}
