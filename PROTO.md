Proto files are stored locally in both `/socketio_client` and `/websocket_client`
Copy the spec you intend to extend into your proto dir, be sure to asjust the package.

You can do one of two things:
 1. - delete the `*.pb.go` file and break the imports
    - then modify the client with your own package.

 - replace the `*.pb.go` file with an extended version using a `go_package` of `socketioClient` or `wssClient` respectivley.
 

NOTE: If you use the extended spec in a binary thats staticly linked to a client in this, you should do the latter.

To modify a client with your own package, add your pb import like so:
```
import (
    pb "example.com/protos/pb"
)
```

To fix the broken imports, add `pb.` before each type. 


If golangs build system cannot retrive `example.com/protos/pb` package:
 - replace the import with its relative path in `go.mod`
 - `replace example.com/protos/pb => ../path/to/pb`

If replacing `/socketio_client` pb, both `client.go` and `decoder.go` depend on the spec