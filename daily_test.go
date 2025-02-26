package reporter_test

import (
	"github.com/adminium/reporter"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestDaily(t *testing.T) {

	d := reporter.NewDaily(
		reporter.WithChannel(reporter.FeiShu, os.Getenv("BOT_URL")),
		reporter.WithExpression("* * * * *"),
	)

	go func() {
		time.Sleep(2 * time.Minute)
		d.Stop()
	}()

	err := d.Start()
	require.NoError(t, err)

	return
}
