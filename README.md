# BWOJ

BWOJ (Blackwind Online Judge) is simple online judge server.

# How to run

1. Install Go

[Go](https://golang.org/) is a modern language developed by Google, suitable for modern computer architecture (multi-core, I/O async...). AnchorVPN uses golang.
```bash
apt-get install golang-go
```

2. Run

[tmux](https://github.com/tmux/tmux/wiki) is a widely used terminal multiplexer. For running AnchorVPN in the background, we need this terminal multiplexer. Other options for process background running is `screen` or adding `&` on the back site of command, but it is not my recommendation. You can install tmux by:
```bash
apt-get install tmux
```
Create a new session using:
```bash
tmux new -s <session name>
```
For example, `<session name>` = anchorvpn. Then type this:
```bash
chmod 755 ./run.sh
./run.sh
```
Detach tmux session by pressing `ctrl + b` then `d`. You're done!

If you want to attach the background session:
```bash
tmux attach -t <session name>
```

3. logging
The logfile is `vpn.log` and the current client database is `clients.db`. The type of the database is [BoltDB](https://github.com/boltdb/bolt). If you have `ccze` then you can keep watching logs by this command:
```bash
watch -n 0.3 -c "cat vpn.log | last -20 | ccze -A"
```

# RESTful API endpoint

```
GET /api/{:project_id}/{:question_id}
```

Post answer to question #{:id}
```
POST /api/{:project_id}/{:question_id}
```