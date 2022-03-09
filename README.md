# srv-customer

Customer Service for Pegadaian Digital Service

## Local Development

### Prerequisites

1. Goland IDE or Visual Studio Code
2. Go 1.17
3. UNIX Shell
   > Use `wsl2` in Windows 10
4. Git
5. Make
6. Docker and Docker Compose CE

### Quick Start

```shell
# Set-up development environment
make configure

# Build server for local development
make servers

# Run database upgrade migration scripts
make db-up

# Start development server
make serve
```

### Configuration

> TODO

### Updating Dependencies

1. Run `go mod vendor` to download dependencies
2. Restart Service container

### Debugging

Attach debugger in your IDE to port 3001 (or as defined in `CUSTOMER_SVC_DEBUG_PORT`)

## Contributors

- Saggaf Arsyad <saggaf@nusantarabetastudio.com>
- Dinan <dinan@nusantarabetastudio.com>
- Imam Ponco <ponco@nusantarabetastudio.com>
