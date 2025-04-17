package dbtypes

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"
)

type DbBigInt struct {
	*big.Int
}

func (dbi *DbBigInt) Scan(value interface{}) error {
	var source string
	switch v := value.(type) {
	case []byte:
		source = string(v)
	case string:
		source = v
	case float64:
		source = fmt.Sprintf("%.0f", v)
	case int64:
		source = fmt.Sprintf("%d", v)
	default:
		return fmt.Errorf("unsupported scan type for DbAddress: %T", value)
	}

	if strings.TrimSpace(source) == "" {
		dbi.Int = nil
		return nil
	}

	parsedInt, ok := new(big.Int).SetString(source, 10)
	if !ok {
		return fmt.Errorf("can't parse %s as a big.Int", source)
	}
	dbi.Int = parsedInt

	return nil
}

func (dbi *DbBigInt) Value() (driver.Value, error) {
	if dbi.Int == nil {
		return nil, fmt.Errorf("can't store nil in not null column")
	}

	return dbi.Int.String(), nil
}
