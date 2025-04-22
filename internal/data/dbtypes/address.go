package dbtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

type DbAddress common.Address

func (da *DbAddress) Scan(value interface{}) error {
	var source string
	switch v := value.(type) {
	case []byte:
		source = string(v)
	case string:
		source = v
	default:
		return fmt.Errorf("unsupported scan type for DbAddress: %T", value)
	}

	if !common.IsHexAddress(source) {
		if source == "" {
			return nil
		}
		return fmt.Errorf("invalid hex address: %s", source)
	}

	address := common.HexToAddress(source)
	*da = DbAddress(address)
	return nil
}

func (da DbAddress) Value() (driver.Value, error) {
	address := common.Address(da)
	return strings.ToLower(address.Hex()), nil
}

func (da DbAddress) MarshalJSON() ([]byte, error) {
	address := common.Address(da)
	return json.Marshal(address.Hex())
}

func (da *DbAddress) UnmarshalJSON(data []byte) error {
	var source string
	if err := json.Unmarshal(data, &source); err != nil {
		return fmt.Errorf("can't unmarshal db address: %w", err)
	}

	if !common.IsHexAddress(source) {
		return fmt.Errorf("invalid hex address: %s", source)
	}

	if source == "" {
		*da = DbAddress(common.Address{})
		return nil
	}

	*da = DbAddress(common.HexToAddress(source))
	return nil
}
