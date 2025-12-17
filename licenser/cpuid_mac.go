package licenser

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

func getMac() (string, error) {
	// Get the MAC address of this unit
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	inter := interfaces[0]
	mac := inter.HardwareAddr.String() // Get the machine MAC address
	return mac, nil
}

func getCpuId() (string, error) {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	str := string(out)
	// Match a regular expression of one or more blank characters
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")
	reg := regexp.MustCompile(`^ProcessorId\s+(\S+)`)
	id := reg.FindStringSubmatch(str)
	if len(id) > 0 {
		return id[1], nil
	}
	return "", fmt.Errorf("not found")
}
