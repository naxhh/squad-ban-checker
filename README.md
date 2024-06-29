# Squad ban cleaner

A very simple tool to check the `Bans.cfg` file for expired bans.

I noticed that sometimes squad will not cleanup them.

## Running

Place the `Bans.cfg` on the root fail and `go run main.go`

It will produce 2 files, `active_bans.cfg` contains all the still valid ban entries, `expired_bans.cfg` contains all expired bans.