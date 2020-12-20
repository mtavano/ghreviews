package database

import (
	"context"
	"testing"

	ghreview "github.com/mtavano/ghreviews"
	"github.com/mtavano/ghreviews/testutil"
	"github.com/stretchr/testify/require"
)

func Test_CreateReview_Success(t *testing.T) {
	db := testutil.NewTestDatabase(t)
	st := &Store{db: db}

	review := &ghreview.GhReview{
		GithubUsername: "foobar",
		Content:        "lorem ipsum dolot ae ben culbprit amus",
	}

	record, err := st.CreateReview(context.TODO(), review)
	defer st.truncateTables([]string{"reviews"})

	require.NoError(t, err)
	require.Equal(t, "foobar", record.GithubUsername)
	require.Nil(t, record.Badge)
}
