package com.icegae.grpcjavaserver;

import greet.Greet;
import greet.GreeterGrpc;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;

import java.io.IOException;

public class GreeterServer {

    public static void main(String[] args) throws IOException, InterruptedException {
        Server server = ServerBuilder
                .forPort(50051)
                .addService(new GreeterServiceImpl())
                .build();

        System.out.println("Starting gRPC server on port 50051...");
        server.start();
        System.out.println("Server listening...");

        server.awaitTermination();
    }

    static class GreeterServiceImpl extends GreeterGrpc.GreeterImplBase {
        @Override
        public void sayHello(Greet.HelloRequest request, StreamObserver<Greet.HelloReply> responseObserver) {
            String name = request.getName();
            String greeting = "Hello, " + name + "!";

            Greet.HelloReply reply = Greet.HelloReply.newBuilder()
                    .setMessage(greeting)
                    .build();

            responseObserver.onNext(reply);
            responseObserver.onCompleted();
        }
    }
}
