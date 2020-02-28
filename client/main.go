package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "github.com/kravcs/go-grpc/customer"
)

const (
	address = "localhost:50051"
)

func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {
	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatal("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatal("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    101,
		Name:  "some name",
		Email: "some@email.com",
		Phone: "123456789",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Some 1.1 Street",
				City:              "Some 1.1 City",
				State:             "Some 1.1 State",
				Zip:               "12345",
				IsShippingAddress: false,
			},
			&pb.CustomerRequest_Address{
				Street:            "Some 1.2 Street",
				City:              "Some 1.2 City",
				State:             "Some 1.2 State",
				Zip:               "67890",
				IsShippingAddress: false,
			},
		},
	}

	createCustomer(client, customer)

	customer = &pb.CustomerRequest{
		Id:    102,
		Name:  "some 2 name",
		Email: "some2@email.com",
		Phone: "987654321",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Some 2 Street",
				City:              "Some 2 City",
				State:             "Some 2 State",
				Zip:               "54321",
				IsShippingAddress: false,
			},
		},
	}

	createCustomer(client, customer)

	filter := &pb.CustomerFilter{Keyword: ""}
	getCustomers(client, filter)
}
