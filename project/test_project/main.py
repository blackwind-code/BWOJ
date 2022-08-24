import json
import random


def question(msg_i: json) -> json:
    def gen():
        a = random.randrange(0, 100)
        b = random.randrange(0, 100)
        return [a, b]

    msg_o = msg_i
    content = []

    for _ in range(100):
        content.append({"q": gen()})

    msg_o["content"] = content
    return msg_o


def judge(msg_i: json) -> json:
    def method(a: int, b: int):
        return a+b

    msg_o = msg_i

    for elem in msg_o["content"]:
        answer = method(elem["q"][0], elem["q"][1])
        if answer == elem["a"]:
        content.append({"q": gen()})

    return msg_o


def select(msg_i: json) -> json:
    qj_id = msg_i["id"]
    method = msg_i["method"]

    if method == "q":
        msg_o = question(msg_i)
        return msg_o
    elif method == "a":
        msg_o = judge(msg_i)
        return msg_o
    else:
        msg_o = msg_i
        msg_o["error"] = f"no such method: {method}"
        return msg_o


if __name__ == '__main__':
    msg_i = json.loads(input())  # get data from stdin
    msg_o = select(msg_i)
    print(json.dumps(msg_o))  # return data to stdout
