package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"strconv"
	"time"
)

func (p *Parser) parseDriverList(dat map[string]interface{}, timestamp time.Time) {

	for driverNum, info := range dat {
		if driverNum == "_kf" {
			continue
		}

		record := info.(map[string]interface{})

		current, exists := p.driverTimes[driverNum]

		if !exists {
			number, _ := strconv.Atoi(driverNum)

			line := 0
			rawLine, exists := record["Line"]
			if exists {
				line = int(rawLine.(float64))
			}

			fullName, _ := record["FullName"].(string)
			shortName, _ := record["Tla"].(string)
			teamName, _ := record["TeamName"].(string)
			teamColour, _ := record["TeamColour"].(string)

			current = Messages.Timing{
				Number:    number,
				Position:  line,
				Name:      fullName,
				ShortName: shortName,
				Team:      teamName,
				Color:     "#" + teamColour,
			}
		}

		p.driverTimes[driverNum] = current
	}
}
