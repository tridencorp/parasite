package config

import (
	"encoding/json"
	"os"
)

// Load json config. We are unmarshalling content of json key to given destination.
// We are reading the whole file each time Load() is called. It simplifies the entire design
// and in most cases, it's only called once at the beginning of the program,
// so performance is not a concern here.
func Load(file, key string, dst any) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	conf := map[string]json.RawMessage{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}

	err = json.Unmarshal(conf[key], dst)
	if err != nil {
		return err
	}

	return nil
}

