from random import randint, choice

first_name = [
    "John",
    "David",
    "Michael",
    "Chris",
    "Mike",
    "Mark",
    "Paul",
    "Daniel",
    "James",
    "Maria",
]

last_name = [
    "Smith",
    "Jones",
    "Johnson",
    "Lee",
    "Brown",
    "Williams",
    "Rodriguez",
    "Garcia",
    "Gonzalez",
    "Lopez",
]

delimiter = [
    ".",
    "_",
    "-",
    "",
]


def random_username():
    return choice(first_name) + choice(delimiter) + choice(last_name) + str(randint(1000000, 999999999))
