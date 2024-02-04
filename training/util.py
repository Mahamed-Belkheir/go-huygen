from socket import socket, AF_INET, SOCK_STREAM
import os

def validate_input(str_input: str) -> list: 
    split_input = str_input.strip().split(",")
    parsed_input = []
    for r in split_input:
        try:
            if r is "":
                continue;
            r = float(r)
            parsed_input.append(r)
        except:
            raise Exception(f"input \"{r}\" is an invalid number")
    return parsed_input    

byte_order = "little"
string_encoding = "utf-8"

interface = os.environ['HUYGENINFER_INTF'] or "127.0.0.1"
port = os.environ['HUYGENINFER_PORT'] or "8080"
addr = (interface, port)


def create_socket():
    return socket(AF_INET, SOCK_STREAM)

def serialize_and_send(sock: socket, str_input: str):
    length_bytes = len(str_input).to_bytes(4, byte_order)
    data = str_input.encode(string_encoding)
    payload = bytes([*length_bytes, *data])
    sock.sendall(payload)
    
def read_and_deserialize(sock: socket):
    length_bytes = sock.recv(4)
    length = int.from_bytes(length_bytes, byte_order)
    payload_bytes = sock.recv(length)
    return payload_bytes.decode(string_encoding)

