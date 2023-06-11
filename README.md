![Build and test](https://github.com/evanebb/gobble/actions/workflows/go.yml/badge.svg)
# gobble
:warning: **This project is made for fun and learning purposes. As such, the code quality or support cannot be guaranteed.** :warning:

A small REST API to manage and generate iPXE configs, and match them to specific clients using their MAC address.

# General concepts

- Simple REST API to create, assign and delete profiles and systems.
- Uses the iPXE configuration language.
- Does not provide a DHCP, TFTP server or iPXE. Those will have to be set up separately, and pointed to this application.
- Can be deployed in addition to your existing netboot infrastructure.

# Process

- DHCP server instructs client to chainload iPXE, points to TFTP server
- TFTP server only contains the iPXE rom(s)
- Server loads the iPXE rom and starts PXE booting again
- DHCP server sees that client has now loaded iPXE and points it to this application; an example URL would be http://gobble.example.local/api/pxe-config?mac=$servermac
- The application retrieves the system using the provided MAC address, and renders the iPXE config on the fly from the profile assigned to it and returns it
- This config contains the kernel, initrd and custom kernel parameters that were assigned. This points to a TFTP, HTTP, NFS, etc. server, which is all out of the control of this application.
- Done!
