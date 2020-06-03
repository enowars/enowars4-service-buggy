import random
import string

from enochecker import *


def random_string(amount):
    return "".join(random.choice(string.ascii_letters + string.digits) for _ in range(amount))


class BuggyChecker(BaseChecker):
    port = 7890
    flag_count = 1
    noise_count = 1
    havoc_count = 1

    def putflag(self) -> None:
        self.logger.debug("Starting putflag")
        try:
            username = random_string(20)
            password = random_string(20)

            # Register account
            response = self.http_post(route="/register", data={"username":username, "pw":password},
                    allow_redirects=False)
            self.logger.debug("register done")

            if response.status_code != 302:
                raise BrokenServiceException("registration failed")

            cookies = response.cookies
            comment = "Awesome!"
            # Post Comment
            buggy = random.choice(["super", "mega"])
            response = self.http_post(route=f"/{buggy}-buggy", data={"comment":comment},
                    cookies=cookies, allow_redirects=False)
            self.logger.debug("comment written")

            if response.status_code != 302:
                raise BrokenServiceException("commenting failed")

            cookies = response.cookies
            subject = random_string(20)
            # Write ticket
            response = self.http_post(route="/tickets", data={"subject":subject, "message":self.flag},
                    cookies=cookies, allow_redirects=False)
            self.logger.debug("ticket written")

            if response.status_code != 302:
                raise BrokenServiceException("ticket failed")

            try:
                hash = response.headers["location"].split("/")[-1].strip()
            except Exception:
                raise BrokenServiceException("ticket redirect failed")

            assert_equals(64, len(hash), "ticket redirect failed")

            self.logger.debug(f"saving hash : {hash}")
            self.team_db[self.flag] = (hash, username, password)

        except Exception:
            raise BrokenServiceException("checker failed")

    def getflag(self) -> None:
        self.logger.debug("Starting getflag")
        try:
            try:
                (hash, user, password) = self.team_db[self.flag]
            except KeyError:
                return Result.MUMBLE

            # Login
            response = self.http_post(route="/login", data={"username": user,
                                "pw": password}, allow_redirects=False)
            self.logger.debug("logged in")

            assert_equals(302, response.status_code, "login failed")

            # TODO: View comment?

            # View ticket
            response = self.http_get(route=f"/tickets/{hash}")
            self.logger.debug("ticket loaded")

            assert_equals(200, response.status_code, "view ticket failed")
            assert_in(self.flag, response.text, "flag not found")

        except Exception:
            raise BrokenServiceException("checker failed")

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
        pass

app = BuggyChecker.service

if __name__ == "__main__":
        run(BuggyChecker)
