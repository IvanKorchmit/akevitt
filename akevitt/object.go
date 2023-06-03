package akevitt

import (
	"bytes"
	"encoding/gob"
)

type Object interface {
	Description() string                    // Retrieve description about that object
	Save(key uint64, engine *Akevitt) error // Save object into database
	OnLoad(engine *Akevitt) error
}

type GameObject interface {
	Object
	Name() string
	Create(engine *Akevitt, session *ActiveSession, params interface{}) error
	GetMap() map[string]Object
	OnRoomLookup() uint64
}

type Room interface {
	Object
	GetExits() []Exit
	SetExits(exits ...Exit)
	GetKey() uint64
}

type Exit interface {
	Object
	GetRoom() Room
	GetKey() uint64
	SetRoom(room Room)
	Enter(engine *Akevitt, session *ActiveSession) error
}

// Converts `T` to byte array
func serialize[T Object](v T) ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(v)

	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Converts byte array to T struct.
func deserialize[T Object](b []byte, engine *Akevitt) (T, error) {
	var result T
	var decodeBuffer bytes.Buffer

	decodeBuffer.Write(b)

	dec := gob.NewDecoder(&decodeBuffer)
	err := dec.Decode(&result)

	if err != nil {
		return result, err
	}

	return result, err
}

func CreateObject[T GameObject](engine *Akevitt, session *ActiveSession, object T, params interface{}) (T, error) {
	return object, object.Create(engine, session, params)
}
