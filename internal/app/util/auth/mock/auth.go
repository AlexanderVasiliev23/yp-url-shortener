package mock

import "context"

// UserContextFetcherMock missing godoc.
type UserContextFetcherMock struct {
	Err    error
	UserID int
}

// GetUserIDFromContext missing godoc.
func (f *UserContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.UserID, f.Err
}
