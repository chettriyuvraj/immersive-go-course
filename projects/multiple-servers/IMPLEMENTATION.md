# Implementation

Setting up a scaled-down version of this exercise, so as to focus on the important aspects

## Modules and packages

- Create a directory skeleton, ignore extras such as flags and environment configs

## Static server

- Again, ignore extras but get the server serving a file system using http.FileServer()

## API server

- Get it running, using in-memory data instead of connecting to databse

## Query API server from Static

- Query API server from static server by modifying fetchImages() request, this will be a cross-origin request which will fail

- Understand what cross-origin requests are and fix this issue by modifying API headers

## Configure nginx

- Configure nginx without upstream (~load-balancing to multiple servers) since we haven't configured flags for ports

- (Was trying to launch nginx on the same port as the server - couldn't figure why it wasn't working)



