# Outbound Privacy Test
This app is designed to simulate a heavy load on your privacy API endpoint similar to how Outbound processes broadcast campaigns. While processing broadcasts, your privacy endpoint will be required to load a lot of users information and return it to Outbound very fast.

Note that this is not an exact simulation of an Outbound broadcast campaign. In terms of time, your broadcast can take a lot longer to run overall due to several factors that can not be replicated in this test.

Tests require a file containing a list of URLs to make requests to.

## Installation
There are 2 installation options. The first is downloading one of the prebuilt binaries for your system.

| Platform | Link |
|:----|:-------:|
| Mac OS X  | [Link]()  |
| Window/amd6 | [Link]() |
| Linux/amd64 | [Link]() |

Alternatively, if you have [Go 1.6+]() installed on your system, you can simply run `go get` to download and install the library.

    go get -u github.com/outboundio/privacy-api-test

## Test Mode
## Unrestricted
To test your endpoint under the most extreme circumstances, run in unrestricted mode. This triggers all connections as fast as it can. The only limit to the number of connections made is the limitations of the system you run the test on.

## Concurrent
By specifying a number of concurrent connections to use you can ensure no more than this number of requests will be made at any given time. This will reduce the load on your own system but will slow down the entire process.

## Usage
Example usage:

    privacy-api-test --file urls.log --concurrent 100 -v

Complete usage options:

    Usage of privacy-api-test:
        -concurrent int
            number of concurrent requests to make if not unrestricted
        -file string
            file of urls to process
        -timeout duration
            http request timeout (default 5s)
        -unrestricted
            make all requests with no restrictions on concurrency
        -v	verbose logging
