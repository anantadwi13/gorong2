package common

import (
	"sync"

	"github.com/rs/xid"
)

type id struct {
	id      xid.ID
	idStr   string
	idStrMu sync.RWMutex
}

func GenerateId() Id {
	return &id{id: xid.New()}
}

func ParseId(idString string) (Id, error) {
	idRaw, err := xid.FromString(idString)
	if err != nil {
		return nil, err
	}
	return &id{id: idRaw}, nil
}

func (i *id) Raw() interface{} {
	return i.id
}

func (i *id) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	otherId, ok := other.(Id)
	if !ok {
		return false
	}
	return otherId.String() == i.String()
}

func (i *id) String() string {
	i.idStrMu.RLock()
	if i.idStr != "" {
		i.idStrMu.RUnlock()
		return i.idStr
	}
	i.idStrMu.RUnlock()

	i.idStrMu.Lock()
	defer i.idStrMu.Unlock()

	i.idStr = i.id.String()
	return i.idStr
}
