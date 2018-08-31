#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import grpc
import kanosdk_pb2 as kano
import kanosdk_pb2_grpc as kanorpc
import os
import uuid


def on_proximity(device, uid):
    yield kano.StreamRequest(
        name=device,
        request=kano.Request(
            type="rpc-request",
            id=uid,
            method="set-mode",
            params=[kano.Param(mode="proximity")],
        ),
    )

def on_gesture(device, uid):
    yield kano.StreamRequest(
        name=device,
        request=kano.Request(
            type="rpc-request",
            id=uid,
            method="set-mode",
            params=[kano.Param(mode="gesture")],
        ),
    )

def communicate(stub, action, device, uid):
    return stub.Communicate(action(device, uid))


def connect(address, device):
    channel = grpc.insecure_channel(address)
    stub = kanorpc.ConnectorStub(channel)
    uuidV4 = uuid.uuid4().hex
    stream = None

    while True:
        if stream is not None:
            for response in stream:
                print(response)
        action = input("1) proximity; 2) gesture; (1 / 2)? ")
        if action == "2":
            stream = communicate(stub, on_gesture, device, uuidV4)
        else:
            stream = communicate(stub, on_proximity, device, uuidV4)


if __name__ == "__main__":
    address = os.getenv("ADDRESS", "localhost:55555")
    device = os.getenv("DEVICE", "/dev/tty.usbmodem14331")
    connect(address, device)
