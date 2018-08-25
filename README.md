# Kano Community SDK (Proof of Concept)

## Summary

This is a proof of concept using cross platform server/client architecture for bidirectional serial communication between Kano Kit(s) and Client connection.

I chose [Protocol Buffers](https://developers.google.com/protocol-buffers/) for serialising data for its platform & language-neutral mechanism, which means 1 definition file, user can generate their own server/client.

I also chose [Go](https://golang.org/) to write the sample server & client for its portability, you can simply generate a client in your preferred language and communicate to my prebuilt server. Of course you can simply run both prebuilt server & client and play around with them, no compilation required.

Linux, Darwin (MacOS) and Windows are prebuilt and stored at the [/bin folder](https://github.com/anzellai/kanosdk/tree/master/bin).


### Example Usage

Just run the server & client with your host platform under bin folder, the server and client should run on separate processes.

Server accepts environment variable for running "port" (default "55555");
Client accepts environment vairable for connecting "address" for server (default "localhost:55555").

To list the available devices, look at either:

Linux/Darwin: `dmesg |grep tty`, or `ls -al /dev/tty.*`;
Windows: `mode`

Input your device serial address, and send data to gRPC server to communicate between serial and your server+client in bidirectional manner.


### TODO

Plenty, like full implementation of device protocols, communication message format / device set-mode, helper functions to list connected devices etc., testing and more robust definitions and more samples in different programming languages.

Contribution welcome, but not for this repository, please submit your issues/ideas/feature requests to [Kano Community SDK](https://github.com/KanoComputing/community-sdk).
