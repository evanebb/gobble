# gobble

A small REST API to easily manage and generate iPXE configs, and match them to specific clients. It promotes reuse of OS distributions and kernel parameters.

# General concepts

- Can be deployed in addition to your existing netboot infrastructure.
- Simple REST API to create, assign and delete distro's, profiles and systems.
- Built on the iPXE configuration language.
- Does not provide a DHCP, TFTP server or iPXE. Those will have to be set up separately, and pointed to this application.

# What does it use?

- PostgreSQL database for storage (using `pgx` for database communication)
- JSON-based REST API with `go-chi`
- UUIDs to identify resources using Google's `uuid` package

# Process

- DHCP server instructs client to chainload iPXE, points to TFTP server
- TFTP server only contains the iPXE image (both UEFI and BIOS, not my problem)
- Server loads the iPXE ROM and starts PXE booting again
- DHCP server sees that client has now loaded iPXE, points it to Gobble server; URL would be http://gobble.example.local/api/pxe-config?mac=$servermac
- Server checks if a system with the specified MAC address exists, errors if it doesn't
- Server renders the iPXE config on the fly from the profile assigned to the system and returns it
- This config contains the kernel, initrd and custom kernel parameters that were assigned. This points to a TFTP, HTTP, NFS, etc. server, which is all out of the control of this application.
- Done!

# TODO
- [x] Wrapper methods for sending responses
- [x] Build Docker image and docker-compose.yml
- [ ] Input validation in factories for models
- [ ] Uniqueness checks (for names, MAC addresses) inside code instead of just in the database
- [ ] Check whether or not supplied linked distro/profile exists before creating instead of letting the foreign key handle it
- [ ] Documentation
- [ ] Tests
- [ ] Proper error handling instead of cascading them to the user through the API
- [ ] Replace primitive types with custom types where necessary (e.g. kernel and initrd)
