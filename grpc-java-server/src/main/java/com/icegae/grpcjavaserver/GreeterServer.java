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
        // 1. Unary
        @Override
        public void sayHello(Greet.HelloRequest request, StreamObserver<Greet.HelloReply> responseObserver) {
            String message = "Hello, " + request.getName();
            Greet.HelloReply reply = Greet.HelloReply.newBuilder().setMessage(message).build();
            responseObserver.onNext(reply);
            responseObserver.onCompleted();
        }

        // 2. Server Streaming
        @Override
        public void greetManyTimes(Greet.HelloRequest request, StreamObserver<Greet.HelloReply> responseObserver) {
            for (int i = 1; i <= 5; i++) {
                Greet.HelloReply reply = Greet.HelloReply.newBuilder()
                        .setMessage("Hello " + request.getName() + " - Message " + i)
                        .build();
                responseObserver.onNext(reply);
                try {
                    Thread.sleep(1000); // simulate delay
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
            responseObserver.onCompleted();
        }

        // 3. Client Streaming
        @Override
        public StreamObserver<Greet.HelloRequest> longGreet(StreamObserver<Greet.HelloReply> responseObserver) {
            return new StreamObserver<>() {
                StringBuilder names = new StringBuilder();

                @Override
                public void onNext(Greet.HelloRequest request) {
                    names.append(request.getName()).append(", ");
                }

                @Override
                public void onError(Throwable t) {
                    t.printStackTrace();
                }

                @Override
                public void onCompleted() {
                    String message = "Hello to: " + names.toString();
                    Greet.HelloReply reply = Greet.HelloReply.newBuilder().setMessage(message).build();
                    responseObserver.onNext(reply);
                    responseObserver.onCompleted();
                }
            };
        }

        // 4. Bidirectional Streaming
        @Override
        public StreamObserver<Greet.HelloRequest> greetEveryone(StreamObserver<Greet.HelloReply> responseObserver) {
            return new StreamObserver<>() {
                @Override
                public void onNext(Greet.HelloRequest request) {
                    String message = "Hello, " + request.getName();
                    Greet.HelloReply reply = Greet.HelloReply.newBuilder().setMessage(message).build();
                    responseObserver.onNext(reply);
                }

                @Override
                public void onError(Throwable t) {
                    t.printStackTrace();
                }

                @Override
                public void onCompleted() {
                    responseObserver.onCompleted();
                }
            };
        }
    }
}
