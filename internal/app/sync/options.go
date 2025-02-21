package sync

import "context"

type Option func(*App)

func WithDownloader(downloader func(ctx context.Context, url, destination string) error) func(app *App) {
	return func(app *App) {
		app.downloader = downloader
	}
}
