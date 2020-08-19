# OpenVPN members

This microservice was done to interact with multiple clients with dynamic IPs,
detecting ports exposed to the vpn.

The idea is, basically, a handy way to use another member services.

## How it works

It watches `./status.log` to retrieve all the connected clients.

Its recommended to simply create a symlink with the original file (ln -s
/var/log/openvpn/status.log ./status.log) or, in case of having more than 
one vpn running, symlinking the desired one.

It uses nmap to catch the exposed ports for each client. In case of desiring
to add more ports to scan (or changing the aliases), the definition is on main.go as a map.
