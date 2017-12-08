package profiler

import (
	"time"

	log "github.com/AKovalevich/scrabbler/log/logrus"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Do.Debugf("%s took %s", name, elapsed)
}