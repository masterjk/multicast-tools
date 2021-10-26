package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"golang.org/x/net/ipv4"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	ifs, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error retrieving all available interfaces")
		os.Exit(-1)
	}

	fmt.Println("Found network interfaces:")
	for _, ifentry := range ifs {
		fmt.Printf("- Index: %d; Name: %s; MTU: %d; Flags: %s; Hardware Address: %s\n", ifentry.Index, ifentry.Name, ifentry.MTU, ifentry.Flags, ifentry.HardwareAddr)
		if addrs, err := ifentry.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Printf("  - IP: %s\n", ipnet.IP.String())
					}
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("Starting multicast-receiver...")
	fmt.Println()

	ifname := flag.String("ifname", "", "interface name (ex: eth0)")
	groupAddressStr := flag.String("group-address", "", "multicast group address (range: 232.0.0.0/8)")
	port := flag.Int("port", 0, "multicast port")
	sourceIPStr := flag.String("source-ip", "", "multicast source IP")
	flag.Parse()

	if len(*ifname) == 0 || len(*groupAddressStr) == 0 || *port == 0 || len(*sourceIPStr) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	fmt.Printf("- Interface     : %s\n", *ifname)
	fmt.Printf("- Group Address : %s\n", *groupAddressStr)
	fmt.Printf("- Port          : %d\n", *port)
	fmt.Printf("- Source IP     : %s\n", *sourceIPStr)
	fmt.Println()

	groupAddress := net.ParseIP(*groupAddressStr)
	if groupAddress == nil {
		log.Errorf("Error parsing groupAddress: ", *groupAddressStr)
		os.Exit(-1)
	}

	sourceIP := net.ParseIP(*sourceIPStr)
	if sourceIP == nil {
		log.Errorf("Error parsing source IP: ", *sourceIPStr)
		os.Exit(-1)
	}

	ifi, err := net.InterfaceByName(*ifname)
	if err != nil {
		log.Errorf("InterfaceByName: ", err)
		os.Exit(-1)
	}

	conn, err := net.ListenPacket("udp4", fmt.Sprintf("%s:%d", *groupAddressStr, *port))
	if err != nil {
		log.Errorf("ListenPacket: ", err)
		os.Exit(1)
	}
	log.Infof("Listening on: %s", conn.LocalAddr())

	pconn := ipv4.NewPacketConn(conn)

	grp := &net.UDPAddr{
		IP:   groupAddress,
		Port: *port,
	}

	src := &net.IPAddr{
		IP: sourceIP,
	}

	if err := pconn.SetMulticastInterface(ifi); err != nil {
		log.Errorf("SetMulticastInterface: ", err)
		os.Exit(-1)
	}

	if err := pconn.JoinSourceSpecificGroup(ifi, grp, src); err != nil {
		log.Errorf("JoinSourceSpecificGroup: ", err)
		os.Exit(-1)
	}

	for {
		buf := make([]byte, 512)
		n, cm, addr, er := pconn.ReadFrom(buf)
		log.Debugf("Message: %s; Size: %d;  ControlMessage: %s; Address: %s; Error: %v", string(buf[:n]), n, cm, addr, er)
	}
}
