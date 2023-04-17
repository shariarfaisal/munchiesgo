package api

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p Point) Value() driver.Value {
	return fmt.Sprintf("(%f,%f)", p.Lat, p.Lng)
}

func getPoint(v interface{}) *Point {
	if v == nil {
		return nil
	}

	var pointData []byte
	switch v.(type) {
	case []byte:
		pointData = v.([]byte)
	case string:
		pointData = []byte(v.(string))

	default:
		return nil
	}

	var point Point
	_, err := fmt.Sscanf(string(pointData), "(%f,%f)", &point.Lat, &point.Lng)
	if err != nil {
		return nil
	}

	return &point
}

func getBrandsSlug(name string, vendorId int64) string {
	return strings.ReplaceAll(name, " ", "-") + "-" + fmt.Sprintf("%d", vendorId)
}

func getPrefix(name string) string {
	fields := strings.Split(name, " ")
	prefix := ""

	if len(name) < 3 {
		return name
	}

	if len(fields) == 1 {
		return name[:3]
	}

	for _, field := range fields {
		prefix += string(field[0])
	}

	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	return prefix
}

var days = map[string]int32{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
}

var daysByIndex = map[int32]string{
	0: "sunday",
	1: "monday",
	2: "tuesday",
	3: "wednesday",
	4: "thursday",
	5: "friday",
	6: "saturday",
}
