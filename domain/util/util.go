package util

import (
	"bytes"
	"github.com/danielcoman/diff"
	"github.com/google/uuid"
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

func (u Util) GenerateUUID(data string) uuid.UUID {
	// Just a rand root uuid
	namespace, _ := uuid.Parse("e8b764da-5fe5-51ed-8af8-c5c6eca28d7a")
	return uuid.NewSHA1(namespace, []byte(data))
}
