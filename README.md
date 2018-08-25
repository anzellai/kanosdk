# Kano Community SDK (Proof of Concept)

## Summary

This is a proof of concept using cross platform server/client architecture for bidirectional serial communication between Kano Kit(s) and Client connection.

I chose [Protocol Buffers](https://developers.google.com/protocol-buffers/) for serialising data for its platform & language-neutral mechanism, which means 1 definition file, user can generate their own server/client.

I also chose [Go](https://golang.org/) to write the sample server & client for its portability, you can simply generate a client in your preferred language and communicate to my prebuilt server. Of course you can simply run both prebuilt server & client and play around with them, no compilation required.

Linux, Darwin (MacOS) and Windows are prebuilt and stored at the [/bin folder](https://github.com/anzellai/kanosdk/tree/master/bin).
