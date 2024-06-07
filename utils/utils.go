package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func ScanRows2map(rows *sql.Rows) (res []map[string]interface{}) {
	defer rows.Close()
	cols, _ := rows.Columns()
	cache := make([]interface{}, len(cols))
	// 为每一列初始化一个指针
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	for rows.Next() {
		rows.Scan(cache...)
		row := make(map[string]interface{})
		for i, val := range cache {
			// 处理数据类型
			v := *val.(*interface{})
			switch v.(type) {
			case []uint8:
				v = string(v.([]uint8))
			case nil:
				v = ""
			}
			row[cols[i]] = v
		}
		res = append(res, row)
	}
	rows.Close()
	return res
}

func ParseNestedJSON(jsonData []byte, keys ...string) (any, error) {
	var result any
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if m, ok := result.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				result = val
			} else {
				return nil, fmt.Errorf("Key '%s' not found", key)
			}
		} else {
			return nil, fmt.Errorf("Value is not a map[string]interface{}")
		}
	}

	return result, nil
}

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
