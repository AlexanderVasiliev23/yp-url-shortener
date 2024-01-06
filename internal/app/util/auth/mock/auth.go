package mock

import "context"

type UserContextFetcherMock struct {
	UserID int
	Err    error
}

func (f *UserContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.UserID, f.Err
}
