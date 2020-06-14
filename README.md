# enowars4-service-buggy
[![Build Status](https://droneci.sect.tu-berlin.de/api/badges/enowars/enowars4-service-buggy/status.svg)](https://droneci.sect.tu-berlin.de/enowars/enowars4-service-buggy)


buggy is a webservice written in [go](https://golang.org/).

**Note that the service and checker are still work in progress, more functionality and vulnerabilities could be added in the future.**

## Service
### Functionality
The Buggy Store offers two buggys. On each product page you can have a look at the description, comment something and preorder the buggy.
You can also send a support ticket from the user profile to the Buggy-Store admin team.

### Technology
The Front- and Backend are written in go, making use of the gohtml templating functionality. The users, comments, messages etc. are stored in a mysql database.
#### Vulnerability
Each ticket gets a (sha256-)hash assigned to it.
```Go
str := acc.User.Username + strconv.FormatInt(time.Now().Unix(), 10)
sha := sha256.Sum256([]byte(str))
hash := hex.EncodeToString(sha[:])
```
Because the username as well as the current time (in seconds) is used, this hash can be forged if the username is known.
Exploiting would involve getting the usernames via the comment sections and brute forcing the timestamp. Because the comments include a timestamp, the amount of enumeration that has to be done is very small though.

The vulnerability can be fixed quite easily be simply changing the sha256 digest contents (i.e. you could hardcode a string in place of the username). Replaying is not possible as the usernames as well as the timestamps are always different which leads to different hashes which need to be accessed.

## Checker
### putflag
The Checker registers an account, comments something on one of the buggys and sends a support ticket.

### getflag
The Checker logs in and checks if the ticket (and comment) is still there.

### putnoise
See putflag, just with a fake flag.

### getnoise
See getflag, just with a fake flag.

### havoc
The Checker registers an account, comments and preorders both buggys, checks for the welcome and preorder messages, and sends a support ticket.
