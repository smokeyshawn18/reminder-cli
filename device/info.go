package device

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

// GetDeviceInfo fetches hostname, OS, arch, username, local IP addresses, and public IP address
func GetDeviceInfo() (hostname, osName, arch, username string, localIPs []string, publicIP string, err error) {
    hostname, err = os.Hostname()
    if err != nil {
        err = fmt.Errorf("error getting hostname: %w", err)
        return
    }

    currentUser, uErr := user.Current()
    if uErr != nil {
        err = fmt.Errorf("error getting current user: %w", uErr)
        return
    }
    username = currentUser.Username
    osName = runtime.GOOS
    arch = runtime.GOARCH

    localIPs, err = getLocalIPs()
    if err != nil {
        err = fmt.Errorf("error fetching local IPs: %w", err)
        return
    }

    publicIP, err = fetchPublicIP()
    if err != nil {
        err = fmt.Errorf("error fetching public IP: %w", err)
        return
    }

    return
}

// getLocalIPs enumerates all non-loopback IPv4 addresses on the machine
func getLocalIPs() ([]string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    var ips []string
    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp == 0 {
            continue // interface down
        }
        if iface.Flags&net.FlagLoopback != 0 {
            continue // skip loopback
        }
        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            if ip == nil || ip.IsLoopback() {
                continue
            }
            ip = ip.To4()
            if ip == nil {
                continue
            }
            ips = append(ips, ip.String())
        }
    }
    return ips, nil
}




// fetchPublicIP requests https://api.ipify.org to find the public IP address.
func fetchPublicIP() (string, error) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
