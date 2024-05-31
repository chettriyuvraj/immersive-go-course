# README

To read the original objective of this exercise, refer to [this](https://github.com/chettriyuvraj/immersive-go-course/tree/main/projects/output-and-error-handling) readme.


I will instead be using this to demonstrate writing a simple client for a provided server in Go + unit-tests for the same.

Given below is a small excerpt from the original exercise explaining what we are doing

## Exercise

In this project, we have been supplied with a server. Our server code lives in the server subdirectory of this project. Run it by cding into that directory, and running _go run_. 

### Server Behaviour
The server is an HTTP server, which listens on port 8080 and responds in a few different ways:

- If you make an HTTP GET request to it, it will respond with the current weather. When this happens, you should display it to the user on the terminal.

- Sometimes this server will overload and respond with a status code 429. When this happens, the client should:

    - Wait the amount of time indicated in the Retry-After response header, and
attempt the request again.

- Sometimes, this server will drop a connection before responding. When this happens, you should assume the server is non-responsive. Making more requests to it could make things worse. The client should:
    - Give up its request.
    - Tell the user something irrecoverable went wrong.

Have a read of the server code. Make sure you understand what it's doing, and what kinds of responses you may need to handle.

## Notes:

- Have written TC for only a single handler, since others will mostly be the same.
- The written code is mostly part of the _fetcher_ directory.
- In my opinion, the main takeaway from this exercise is how to mock a server to test your client using the _RoundTripper_ interface.