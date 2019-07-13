# go-minecraftping
A wrapper for querying [Minecraft Java Edition](https://minecraft.net) servers using the vanilla [Server List Ping](https://wiki.vg/Server_List_Ping) protocol.

## Usage
### Installation
Install using ```go get github.com/Cryptkeeper/go-minecraftping```

### Example Usage
```golang
package main

import (
	"fmt"
	"github.com/Cryptkeeper/go-minecraftping"
	"log"
	"time"
)

func main() {
	resp, err := minecraftping.Ping("myip", 25565, protocolVersion, time.Second * 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d/%d players are online.", resp.Players.Online, resp.Players.Max)
}
```

(The default Minecraft port, 25565, is also available as a const, ```minecraftping.DefaultPort```.)

```protocolVersion``` is ever changing as Minecraft updates. See [protocol version numbers](https://wiki.vg/Protocol_version_numbers) for a complete and updated listing. If compatible with the sent protocol version, the server will reply with the same protocol version in the ```Response``` object, otherwise it will send its required protocol version (keep in mind, some servers may be compatible with _multiple_ protocol versions.)

Minecraft's latest protocol version is available as a const, ```minecraftping.LatestProtocolVersion``` to help reduce magic numbers in basic usages of the library.

### Response
The response structure is described in [```minecraftping.Response```](https://github.com/Cryptkeeper/go-minecraftping/blob/master/minecraftping.go#L40)

## Behavior Notes
1. This does not support Minecraft's [legacy ping protocol](https://wiki.vg/Server_List_Ping#1.6) for pre-Minecraft version 1.6 servers.
2. The ```description``` field of ```Response``` is provided as a ```json.RawMessage``` object. This is because the field's schema follows the [Chat schema](https://wiki.vg/Chat) (a Minecraft specific schema) that I'm not willing to support at this functionality level.
3. This does not support the ```Ping``` or ```Pong``` behavior of the [Server List Ping](https://wiki.vg/Server_List_Ping) protocol. If you wish to determine the latency of the connection do you should do so manually. 