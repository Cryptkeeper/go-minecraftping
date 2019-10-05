// Package minecraftping is a simple library to ping Minecraft Java Edition servers.
package minecraftping

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	// DefaultPort is the default Minecraft Java Edition network port.
	DefaultPort = 25565

	// LatestProtocolVersion is the latest released Minecraft Java Edition protocol version.
	// This field is manually updated and may be out of date.
	// Corresponding Minecraft version: 1.14.4
	// More protocol versions: https://wiki.vg/Protocol_version_numbers
	LatestProtocolVersion = 498
)

// requestPacket is a cached copy of the Request packet.
// It contains only it's length (1) and the packet's ID (0).
var requestPacket = []byte{1, 0}

// Response is a representation of the Minecraft Java Edition server's ping response.
// More information: https://wiki.vg/Server_List_Ping#Response
type Response struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description json.RawMessage `json:"description"`
	Favicon     string          `json:"favicon"`
}

// Ping connects and pings the Minecraft Java Edition server at the specified address and port.
// protocolVersion dictates which protocol version to attempt the ping with as the response is protocol version dependent.
// More information: https://wiki.vg/Server_List_Ping
func Ping(address string, port uint16, protocolVersion int, timeout time.Duration) (Response, error) {
	var resp Response

	deadline := time.Now().Add(timeout)

	conn, err := net.DialTimeout("tcp", address+":"+strconv.Itoa(int(port)), timeout)
	if err != nil {
		return resp, err
	}
	defer conn.Close()

	if err := conn.SetDeadline(deadline); err != nil {
		return resp, err
	}

	// Construct and write Handshake packet to open connection and then write Request packet.
	// More information: https://wiki.vg/Server_List_Ping
	handshake := makeHandshakePacket(address, port, protocolVersion)
	conn.Write(handshake)

	conn.Write(requestPacket)

	reader := bufio.NewReader(conn)

	// Read and discard the length of the incoming packet.
	_, err = binary.ReadUvarint(reader)
	if err != nil {
		return resp, err
	}

	// Read the packet ID and validate it as 0.
	packetId, err := binary.ReadUvarint(reader)
	if err != nil {
		return resp, err
	}
	if packetId != 0 {
		return resp, fmt.Errorf("received invalid packetId (expected 0!) %d", packetId)
	}

	// Read the length of the incoming JSON payload (as a uvarint). Read the following bytes into a buffer and then
	// unmarshal the []byte into its struct representation Response.
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return resp, err
	}
	payload := make([]byte, length)
	if _, err = io.ReadFull(reader, payload); err != nil {
		return resp, err
	}
	if err = json.Unmarshal(payload, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func makeHandshakePacket(address string, port uint16, protocolVersion int) []byte {
	var buf bytes.Buffer

	buf.Write([]byte("\x00"))

	putVarInt(&buf, int32(protocolVersion))

	putVarInt(&buf, int32(len(address)))
	buf.WriteString(address)

	binary.Write(&buf, binary.BigEndian, port)

	putVarInt(&buf, 1)

	// Prepend the buffer with it's length as a uvarint
	var out bytes.Buffer

	putVarInt(&out, int32(buf.Len()))
	out.Write(buf.Bytes())

	return out.Bytes()
}

// Allocate a []byte buffer of binary.MaxVarintlen32 and write value as a uvarint32. Trim and write to buf.
func putVarInt(buf *bytes.Buffer, value int32) {
	bytes := make([]byte, binary.MaxVarintLen32)
	bytesWritten := binary.PutUvarint(bytes, uint64(value))

	buf.Write(bytes[:bytesWritten])
}
