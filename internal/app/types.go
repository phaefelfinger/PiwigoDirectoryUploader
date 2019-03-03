package app

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
)

type appContext struct {
	Piwigo         *piwigo.PiwigoContext
	SessionId      string
	LocalRootPath  string
	ChunkSizeBytes int
}
