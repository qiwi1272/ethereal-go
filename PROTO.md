To extend the proto file, replace each import of: `pb "github.com/qiwi1272/ethereal-go/_pb"`
With your protobuf package: `example.com/protos/pb`

You should include the `/_pb/*.proto` files in `example.com/protos/pb`, or at least keep all the messages you intend to use.

If golangs build system cannot retrive `example.com/protos/pb` package,
replace the import with its relative path in `go.mod`: `replace example.com/protos/pb => ../path/to/pb`
