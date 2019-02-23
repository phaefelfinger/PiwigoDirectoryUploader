package app

import "haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/authentication"

type AppContext struct {
	Piwigo         *authentication.PiwigoContext
	SessionId      string
	LocalRootPath  string
	ChunkSizeBytes int
}
