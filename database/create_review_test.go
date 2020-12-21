package database

import (
	"context"
	"testing"

	"github.com/mtavano/ghreviews/testutil"
	"github.com/stretchr/testify/require"
)

func Test_CreateReview_Success(t *testing.T) {
	db := testutil.NewTestDatabase(t)
	st := &Store{db: db}

	record, err := st.CreateReview(context.TODO(), "foobar", "https://some.image.url/foo", "lorem ipsum dolot ae ben culbprit amus", nil)
	defer st.truncateTables([]string{"reviews"})

	require.NoError(t, err)
	require.Equal(t, "foobar", record.GithubUsername)
	require.Nil(t, record.Badge)
}
