package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"strconv"
	"time"
)

func (p *Parser) parsePositionData(dat map[string]interface{}, timestamp time.Time) ([]Messages.Location, error) {

	result := make([]Messages.Location, 0)

	for _, record := range dat["Position"].([]interface{}) {
		timestampStr := record.(map[string]interface{})["Timestamp"].(string)
		dataTimestamp, err := parseTime(timestampStr)
		if err != nil {
			p.ParseTimeError(connection.PositionFile, timestamp, "Timestamp", err)
		}

		for key, entry := range record.(map[string]interface{})["Entries"].(map[string]interface{}) {
			driver, _ := strconv.ParseInt(key, 10, 8)
			//status := entry.(map[string]interface{})["Status"].(string)
			x := entry.(map[string]interface{})["X"].(float64)
			y := entry.(map[string]interface{})["Y"].(float64)
			z := entry.(map[string]interface{})["Z"].(float64)

			result = append(result, Messages.Location{
				Timestamp:    dataTimestamp,
				DriverNumber: int(driver),
				X:            x,
				Y:            y,
				Z:            z,
			})
		}
	}

	return result, nil
}
