package mock

import "context"

type UserContextFetcherMock struct {
	userID int
	err    error
}

func (f *UserContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.userID, f.err
}
