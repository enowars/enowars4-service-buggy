# enowars4-service-buggy
[![Build Status](https://droneci.sect.tu-berlin.de/api/badges/enowars/enowars4-service-buggy/status.svg)](https://droneci.sect.tu-berlin.de/enowars/enowars4-service-buggy)


buggy is a webservice written in [go](https://golang.org/).

## Service
### Functionality
The Buggy Store offers two buggys. On each product page you can have a look at the description, comment something and preorder the buggy.
You can also send a support ticket from the user profile to the Buggy-Store admin-team.

### Technology
The Front- and Backend are written in go, making use of the gohtml templating functionality. The users, comments, messages etc. are stored in a mysql database.
#### Vulnerabilities
##### Tickets
Each ticket gets a hash assigned to it.

Because only the username as well as the current time (in seconds) is used, this hash can be forged if the username is known.
Exploiting involves getting the usernames via the orders and brute forcing the timestamp. The orders themselves have a hash assigned to them which again can be forged.

The vulnerability can be fixed quite easily be simply changing the hash digest contents (i.e. you could hardcode a string in place of the username). Replaying is possible only when analyzing all requests to `/orders/` as there are only 48 uniquely different hashes possible.

##### Status
The `/user/` endpoint displays user information.

Because `db.InsertUser` returns true even when the insert was faulty, you can get an authenticated session with an empty `User` struct.
Browsing `/user/` (with any username) will then cause all Users to be printed due to a faulty if-comparison in `keepUser`:
```go
if usernameSession != "" && (user.Username != usernameSession || usernameURL != user.Username) {
```

In the `status` part of the user information flags are leaked.

Fixing this vuln can be done at multiple places, for instance you can just remove the `usernameSession != ""` part of if-comparison.

##### PoC
PoC exploits can be found in the `checker.py`.
