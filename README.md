# Outbound Privacy Test
This app is designed to simulate a heavy load on your privacy API endpoint similar to how Outbound processes broadcast campaigns. While processing broadcasts, your privacy endpoint will be required to load a lot of users information and return it to Outbound very fast.

Note that this is not an exact simulation of an Outbound broadcast campaign. In terms of time, your broadcast can take a lot longer to run overall due to several factors that can not be replicated in this test.

Tests require a file containing a list of URLs to make requests to.

## Test Mode
## Unrestricted
To test your endpoint under the most extreme circumstances, run in unrestricted mode. This triggers all connections as fast as it can. The only limit to the number of connections made is the limitations of the system you run the test on.

## Concurrent
By specifying a number of concurrent connections to use you can ensure no more than this number of requests will be made at any given time. This will reduce the load on your own system but will slow down the entire process.

## Usage

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