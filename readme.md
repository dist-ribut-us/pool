## Pool

### Services and Applications
An application has a user interface, it's something that users care about. A
service runs in the background. A program is either a service or application.

All programs are issued a port to communicate on. Pool will only respond to
requests from that port. They are all issued their key at startup. A program
should not store their key anywhere else.