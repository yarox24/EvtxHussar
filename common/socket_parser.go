package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

func IPPortFormatter(version int, IP net.IP, Port uint16, ScopeId uint32) string {

	out := ""

	if version == 4 {
		out += IP.String()
	} else if version == 6 {
		out += "[" + IP.String()
		if ScopeId > 0 {
			out += "%" + strconv.Itoa(int(ScopeId))
		}
		out += "]"
	} else {
		panic("IPPortFormatter - wrong version")
	}

	if Port != 0 {
		out += ":" + strconv.Itoa(int(Port))
	}

	return out

	// return ip.String() + ":" + strconv.Itoa(int(decoded_socket.Sin_port)), nil
}

func DecodeASIPv4(socket_bytes []uint8) (string, error) {
	type sockaddr_in struct {
		Sin_family uint16  // e.g. AF_INET
		Sin_port   uint16  // e.g. htons(3490)
		Sin_addr   uint32  // unsigned long s_addr;  // load with inet_aton()
		Sin_zero   [8]byte // zero this if you want to
	}

	var decoded_socket sockaddr_in
	err := binary.Read(bytes.NewReader(socket_bytes), binary.BigEndian, &decoded_socket)

	if err != nil {
		return "", err
	}

	ip := net.IPv4(
		byte(decoded_socket.Sin_addr>>24),
		byte(decoded_socket.Sin_addr>>16),
		byte(decoded_socket.Sin_addr>>8),
		byte(decoded_socket.Sin_addr),
	)

	return IPPortFormatter(4, ip, decoded_socket.Sin_port, 0), nil
}

func DecodeASIPv6(socket_bytes []uint8) (string, error) {
	type sockaddr6_in struct {
		Sin6_family   uint16   // AF_INET6
		Sin6_port     uint16   // port number
		Sin6_flowinfo uint32   // IPv6 flow information
		Sin6_addr     [16]byte // IPv6 address
		Sin6_scope_id uint32   // Scope ID
	}

	var decoded_socket6 sockaddr6_in
	err := binary.Read(bytes.NewReader(socket_bytes), binary.BigEndian, &decoded_socket6)

	if err != nil {
		return "", err
	}

	// Ipv6
	ip6 := net.IP(decoded_socket6.Sin6_addr[:])

	// Scope ID [Offset from 24-27] Not sure if implementation fully cover all cases. Probably not
	binary.Read(bytes.NewReader(socket_bytes[24:28]), binary.LittleEndian, &decoded_socket6.Sin6_scope_id)

	return IPPortFormatter(6, ip6, decoded_socket6.Sin6_port, decoded_socket6.Sin6_scope_id), nil
}

func ConvertSocketUint8ToText(socket_bytes []uint8) (string, error) {

	// IPv4 socket
	//	struct sockaddr_in {
	//		short   sin_family; 		// e.g. AF_INET
	//		u_short sin_port;   		// e.g. htons(3490)
	//  	struct  in_addr sin_addr; 	// unsigned long s_addr;  // load with inet_aton()
	//		char    sin_zero[8]; 		// zero this if you want to
	//	};

	// IPv6 socket
	//struct sockaddr_in6 {
	//	unsigned short 	sin6_family;   // AF_INET6
	//	uint16_t        sin6_port;     // port number
	//	uint32_t        sin6_flowinfo; // IPv6 flow information
	//struct in6_addr sin6_addr;     // IPv6 address
	//	uint32_t        sin6_scope_id; // Scope ID (new in 2.4)
	//};
	//
	//struct in6_addr {
	//	unsigned char   s6_addr[16];   // IPv6 address
	//};

	if len(socket_bytes) == 16 {
		// IPv4 case (IP + Port)
		return DecodeASIPv4(socket_bytes)
	} else if len(socket_bytes) == 28 {
		// IPv6 case (IP + Port)
		return DecodeASIPv6(socket_bytes)
	} else if len(socket_bytes) == 128 {
		// IPv4 or Ipv6 (Dynamic)
		if socket_bytes[0] == 2 && socket_bytes[1] == 0 {
			// IPv4 Case
			return DecodeASIPv4(socket_bytes)
		} else if socket_bytes[0] == 23 && socket_bytes[1] == 0 {
			// IPv6 Case
			return DecodeASIPv6(socket_bytes)
		} else {
			// Unsupported case
			return "", errors.New("ConvertSocketUint8ToText - Unsupported 128 len case")

		}
	} else if len(socket_bytes) == 0 {
		return "", nil
	} else {
		// Unsupported case
		return "", errors.New("ConvertSocketUint8ToText - Unsupported case")
	}

}
