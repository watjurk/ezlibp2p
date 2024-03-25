# ezlibp2p - easy libp2p

[![PkgGoDev](https://pkg.go.dev/badge/github.com/watjurk/ezlibp2p)](https://pkg.go.dev/github.com/watjurk/ezlibp2p)

Quality of life and common functionality for [go-libp2p](https://github.com/libp2p/go-libp2p).

Have you ever felt that for every single project that uses [libp2p](https://libp2p.io/) you had to reimplement the same functionality over and over again? Well, me too... ezlibp2p implements this common functionality.

Most notably:

- bootstrapping
- auto autorelay finding
- persistent identity
- writing/reading protobuf messages
- mdns

## Project Status

This library must NOT be considered stable. I will definitely stabilize it one day but for that to happen I must gain more experience with how to handle go-libp2p.

I am using ezlibp2 in all of my projects that use libp2p, every now and then I find a better solution than the one present in ezlibp2p, and then I slightly change the API of ezlibp2p to account for that.

## Contributing

All contributions are welcome! I have nothing more to say than to respect others and use GitHub features.

If you would like to submit a pull request, before doing any work create an issue and I will reach out to you so we can discuss it.

Thanks for sticking around, stay curious! Best,  
Wiktor
