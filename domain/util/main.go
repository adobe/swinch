package util

import (
	"bytes"
	"github.com/danielcoman/diff"
	log "github.com/sirupsen/logrus"
)

type Util struct {
}

func (u Util) Changes(oldData, newData []byte) bool {
	changes := bytes.Compare(oldData, newData)
	if changes == 0 {
		log.Infof("No changes detected")
		return false
	}

	return true
}

func (u Util) DiffChanges(oldData, newData []byte) {
	log.Infof(diff.LineDiff(string(oldData), string(newData)))
}
