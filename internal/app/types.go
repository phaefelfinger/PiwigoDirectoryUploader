package app

import (
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
)

type AppContext struct {
	Piwigo         *piwigo.PiwigoContext
	SessionId      string
	LocalRootPath  string
	ChunkSizeBytes int
}
