package dbtypes

import (
	"database/sql/driver"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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

func (da *DbAddress) Value() (driver.Value, error) {
	if da == nil {
		return nil, fmt.Errorf("can't store nil in not null column")
	}

	address := common.Address(*da)
	return address.Hex(), nil
}
