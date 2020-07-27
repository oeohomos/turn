package turn

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/pion/transport/vnet"
)

// RelayAddressGeneratorForRange can be used to return static IP address each time a relay is created.
// This can be used when you have a single static IP address that you want to use
type RelayAddressGeneratorForRange struct {
	// RelayAddress is the IP returned to the user when the relay is created
	RelayAddress net.IP

	// Address is passed to Listen/ListenPacket when creating the Relay
	Address string

	Net *vnet.Net

	MinPort int
	MaxPort int
}

// Validate is caled on server startup and confirms the RelayAddressGenerator is properly configured
func (r *RelayAddressGeneratorForRange) Validate() error {
	if r.Net == nil {
		r.Net = vnet.NewNet(nil)
	}

	switch {
	case r.RelayAddress == nil:
		return errRelayAddressInvalid
	case r.Address == "":
		return errListeningAddressInvalid
	default:
		return nil
	}
}

// AllocatePacketConn generates a new PacketConn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorForRange) AllocatePacketConn(network string, requestedPort int) (net.PacketConn, net.Addr, error) {
	requestedPort = GetFreePortInRange(r.MinPort, r.MaxPort)
	if requestedPort == 0 {
		return nil, nil, errors.New("没有可用端口")
	}
	conn, err := r.Net.ListenPacket(network, r.Address+":"+strconv.Itoa(requestedPort))
	if err != nil {
		return nil, nil, err
	}

	// Replace actual listening IP with the user requested one of RelayAddressGeneratorForRange
	relayAddr := conn.LocalAddr().(*net.UDPAddr)
	relayAddr.IP = r.RelayAddress

	return conn, relayAddr, nil
}

// AllocateConn generates a new Conn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorForRange) AllocateConn(network string, requestedPort int) (net.Conn, net.Addr, error) {
	return nil, nil, fmt.Errorf("TODO")
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePortInRange(min, max int) int {
	for i := (max - min) * 2; i >= 0; i-- {
		randPort := rand.Intn(max+1-min) + min
		if CheckFreePort(randPort) {
			return randPort
		}
	}
	return 0
}

func CheckFreePort(port int) bool {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return false
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return false
	}
	defer l.Close()
	return true
}
