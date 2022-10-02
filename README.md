# gorong2

gorong2 is just a command that lets you expose your local apps to be accessed from internet

## Features
- Self-hosted
- Tunnel protocol (Edge type)
  - TCP
  - UDP
- Load balancer (round-robin)
- Internal communication protocol (Backbone type)
  - Websocket (support over TLS too)

## How to use
Prerequisite
- One server has a public IP
- Able to whitelist some ports from firewall on the server

Steps
1. Connect to your server via SSH
2. Download the precompiled binary `gorong2-server`. You can use `wget` command to download the binary. Don't forget to give executable permission on the downloaded binary (`chmod +x /path/to/gorong2-server`)
3. Run the `gorong2-server`. Make sure to change the token with your random-generated string
   ```shell
   /path/to/gorong2-server --token your-random-string # it will run on port :19000
   # Output: 2022/10/02 10:56:38 Token your-random-string
   ```
4. Whitelist some ports from firewall. For instance, you can open port 19000 (backbone) and 19080 (edge server)
5. On your local device, download the precompiled binary `gorong2-agent`. Don't forget to give the executable permission
6. Run the `gorong2-agent`
   ```shell
   /path/to/gorong2-agent --backbone server_public_ip:19000 \
      --agent localhost:80 \
      --server :19080 \
      --token your-random-string
   # make sure to change server_public_ip
   ```
7. Run your local app on port 80
8. Then, access your local app via `server_public_ip:19080`

Notes
- By default, the `gorong2-server` will not run the backbone server via secured websocket (Websocket over TLS). So, you need to create your certificates and enable it by yourself. Please check `--backbone-*` arguments
- The above instruction will enable TCP tunnel. You can also enable UDP tunnel by passing `--tunneltype udp` on `gorong2-agent`
- To make the load balancer works, you need to run the other agents with the same `--tunnelid`

## Architecture
```
┌─────────────────────┐  ┌──────────────┐  ┌──────────────────┐  ┌──────────────────────────────────────┐
│      Internet       │  │    Server    │  │     Internet     │  │            Local Network             │
│                     │  │    Network   │  │                  │  │                                      │
│ ┌────────┐       ┌──┴──┴──┐       ┌───┴──┴───┐          ┌───┴──┴───┐       ┌────────┐       ┌───────┐ │
│ │ Public │       │ Edge   │       │ Backbone │          │ Backbone │       │ Edge   │       │ Local │ │
│ │ Client │◄─────►│ Server │◄─────►│ Server   │◄────────►│ Agent    │◄─────►│ Agent  │◄─────►│ App   │ │
│ └────────┘       └──┬──┬──┘       └───┬──┬───┘          └───┬──┬───┘       └────────┘       └───────┘ │
│                     │  │              │  │                  │  │                                      │
│                     │  │              │  │                  │  │                                      │
└─────────────────────┘  └──────────────┘  └──────────────────┘  └──────────────────────────────────────┘
```

## Compile the binary
To compile the binary by yourself, you can do the following steps
1. Make sure you have installed golang with version >= v1.18
2. Build the `gorong2-server`
   ```shell
   go build -o ./build/gorong2-server ./cmd/server
   # or
   make build-server
   ```
3. Build the `gorong2-agent`
   ```shell
   go build -o ./build/gorong2-agent ./cmd/agent
   # or
   make build-agent
   ```
4. The compiled binaries are located in `./build`

## Things need to be considered
Since it is still in the development phase, there may be a lot of bugs like performance issues, memory leaks, security issues, and more. So, use gorong2 on your own responsibility. 