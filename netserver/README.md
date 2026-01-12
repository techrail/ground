# Net Server

This is a HTTP server structure which is based on the `net/http` golang package and thus does not depend on any external package.

## Usage

To use the netserver, you need to first create an instance of the server like this:

```golang
// Set port and blockOnStart params according to your need
s := NewSever(8080, true) 
```

Then create a controller like this:

```golang
// Create a new Controller cl
```
