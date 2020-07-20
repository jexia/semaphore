# API

This package contains API artifacts such as protobuf annotations.

## Protobuf Usage

1. Define your [gRPC](https://grpc.io/docs/) service using protocol buffers 

   `your_service.proto`:
   ```protobuf
   syntax = "proto3";
   package example;
   message StringMessage {
     string value = 1;
   }

   service YourService {
     rpc Echo(StringMessage) returns (StringMessage) {}
   }
   ```

2. Add a [`semaphore.api`](https://github.com/jexia/semaphore/blob/master/api/annotations.proto)
annotation to your .proto file

   `your_service.proto`:
   ```diff
    syntax = "proto3";
    package example;
   +
   +import "semaphore/api/annotations.proto";
   +
    message StringMessage {
      string value = 1;
    }

    service YourService {
   +  option (semaphore.api.service) = {
   +    host: "127.0.0.1:80"
   +    transport: "http"
   +    codec: "json"
   +  };
   +
   -  rpc Echo(StringMessage) returns (StringMessage) {}
   +  rpc Echo(StringMessage) returns (StringMessage) {
   +    option (semaphore.api.http) = {
   +      post: "/v1/example/echo"
   +      body: "*"
   +    };
   +  }
    }
   ```

   >You will need to provide the required third party protobuf files to the `protoc` compiler.
   >They are included in this repo under the `api` folder, and we recommend copying
   >them into your `protoc` generation file structure. If you've structured your protofiles according
   >to something like [the Buf style guide](https://buf.build/docs/style-guide#files-and-packages),
   >you could copy the files into a top-level `./semaphore` folder.

   If you do not want to modify the proto file for use with grpc-gateway you can
   alternatively use an external
   [Service Configuration](https://github.com/jexia/semaphore/tree/master/cmd/semaphore/config) file.
   [Check our documentation](https://jexia.gitbook.io/semaphore/getting-started/cli)
   for more information.

3. Write flow your definitions as usual
