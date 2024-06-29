package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Ban struct {
	PlayerSteamID string
	Expiration    int64
	Admin         string
	AdminSteamID  string
	Comment       string
	RawLine       string
}

func main() {
	bans, err := readBansFromFile("Bans.cfg")
	if err != nil {
		fmt.Println("Failed to read bans:", err)
		return
	}

	expiredCount := 0
	activeCount := 0

	expiredFile, err := os.Create("expired_bans.cfg")
	if err != nil {
		fmt.Println("Failed to create expired bans file:", err)
		return
	}
	defer expiredFile.Close()

	activeFile, err := os.Create("active_bans.cfg")
	if err != nil {
		fmt.Println("Failed to create active bans file:", err)
		return
	}
	defer activeFile.Close()

	for _, ban := range bans {
		currentUnixTime := time.Now().Unix()

		if ban.Expiration != 0 && ban.Expiration < currentUnixTime {
			expiredCount++
			_, err := expiredFile.WriteString(ban.RawLine + "\n")
			if err != nil {
				fmt.Println("Failed to write expired ban:", err)
			}
		} else {
			activeCount++
			_, err := activeFile.WriteString(ban.RawLine + "\n")
			if err != nil {
				fmt.Println("Failed to write active ban:", err)
			}
		}
	}

	fmt.Println("Expired bans:", expiredCount)
	fmt.Println("Active bans:", activeCount)
	fmt.Println("Bans successfully processed!")
}

func readBansFromFile(filename string) ([]Ban, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bans []Ban
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		ban, err := parseBan(line)
		if err != nil {
			fmt.Println("Failed to parse ban:", err)
			continue
		}

		bans = append(bans, ban)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return bans, nil
}

func parseBan(line string) (Ban, error) {
	var ban Ban

	parts := strings.Split(line, ":")

	// Simple ban format: id:timestamp
	if len(parts) == 2 {
		expirationUnix, err := toTimestamp(parts[1])
		if err != nil {
			return ban, fmt.Errorf("failed to parse ban expiration: %s", err)
		}

		ban.PlayerSteamID = strings.TrimSpace(parts[0])
		ban.Expiration = expirationUnix
		ban.RawLine = line

		return ban, nil
	}

	// Complex ban format: <Admin> [SteamID <Admin ID>] Banned:<Banned ID>:<Timestamp> //<Comment>
	if len(parts) == 3 {
		expirationUnix, err := toTimestamp(parts[2])
		if err != nil {
			return ban, fmt.Errorf("failed to parse ban expiration: %s", err)
		}

		ban.PlayerSteamID = strings.TrimSpace(parts[1])
		ban.Expiration = expirationUnix
		ban.RawLine = line

		return ban, nil
	}

	return ban, nil
}

func toTimestamp(str string) (int64, error) {
	expirationStr := trimNonNumericCharacters(str)
	return strconv.ParseInt(expirationStr, 10, 64)
}

func trimNonNumericCharacters(str string) string {
	newStr := ""

	for i := len(str) - 1; i >= 0; i-- {
		if str[i] >= '0' && str[i] <= '9' {
			newStr = string(str[i]) + newStr
		} else {
			// if we find something that is not a number just stop processing.
			break
		}
	}

	return newStr
}
