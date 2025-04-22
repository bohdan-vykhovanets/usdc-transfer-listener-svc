package dbtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

type DbHash common.Hash

func (dh *DbHash) Scan(value interface{}) error {
	var source string
	switch v := value.(type) {
	case []byte:
		source = string(v)
	case string:
		source = v
	default:
		return fmt.Errorf("unsupported scan type for DbHash: %T", value)
	}

	if !strings.HasPrefix(source, "0x") || len(source) != 66 {
		if source == "" {
			return nil
		}
		return fmt.Errorf("invalid hex hash: %s", source)
	}

	hash := common.HexToHash(source)
	*dh = DbHash(hash)
	return nil
}

func (dh DbHash) Value() (driver.Value, error) {
	hash := common.Hash(dh)
	return strings.ToLower(hash.Hex()), nil
}

func (dh DbHash) MarshalJSON() ([]byte, error) {
	hash := common.Hash(dh)
	return json.Marshal(hash.Hex())
}

func (dh *DbHash) UnmarshalJSON(data []byte) error {
	var source string
	if err := json.Unmarshal(data, &source); err != nil {
		return fmt.Errorf("can't unmarshal db hash: %w", err)
	}

	if !strings.HasPrefix(source, "0x") || len(source) != 66 {
		return fmt.Errorf("invalid hex hash: %s", source)
	}

	if source == "" {
		*dh = DbHash(common.Hash{})
		return nil
	}

	*dh = DbHash(common.HexToHash(source))
	return nil
}
