![Build and test](https://github.com/evanebb/gobble/actions/workflows/go.yml/badge.svg)
# gobble
:warning: **This project is made for fun and learning purposes. As such, the code quality or support cannot be guaranteed.** :warning:

A small REST API to manage and generate iPXE configs, and match them to specific clients using their MAC address.

Check out the [quickstart article on the wiki](https://github.com/evanebb/gobble/wiki/Quickstart) for an example setup.

# General concepts
- Simple REST API to create, assign and delete profiles and systems.
- Uses the iPXE scripting language.
- Does not include a DHCP, TFTP server or iPXE firmware. Those will have to be set up separately, giving you the flexibility to choose whatever you want or already have.

# Process
- A client boots from the network, and starts sending DHCP requests.
- The DHCP server instructs the client to chainload iPXE, and points it towards the TFTP server and proper iPXE firmware file to do so.
- The TFTP server serves the iPXE firmware to the client.
- The client loads the iPXE firmware; it starts PXE booting and sending DHCP requests again.
- The DHCP server sees that the client has now loaded iPXE, and points it towards Gobble to retrieve an iPXE script; an example URL is http://gobble.example.local/api/pxe-config?mac=$servermac
- Gobble looks up the registered system using the provided MAC address, and renders the iPXE config on the fly from the profile assigned to it.
- This config contains the kernel, initrd and custom kernel parameters that were assigned. This points to a TFTP, HTTP, NFS, etc. server, which is all out of the control of this application.
- Done!
