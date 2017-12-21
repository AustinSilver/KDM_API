package main

import (
	"io"
	"fmt"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	prsn "github.com/austinsilver/kdm_api/person"
)

const (
	address = "localhost:3333"
)

// createPerson calls the RPC method CreatePerson of PersonServer

func createPerson(client prsn.PersonClient, person *prsn.PersonRequest){

	resp, err := client.CreatePerson(context.Background(), person)

	if err != nil {
		fmt.Println("Could not create Person: ", err)
		return
	}
	if resp.Success {
		fmt.Println("A new Person has been added with id :", resp.Id)

	}
}
// getPersons calls that RPC method GetPersons from PersonServer
func getPersons(client prsn.PersonClient, filter *prsn.PersonFilter) {
	// Call the streaming API
	stream, err := client.GetPersons(context.Background(), filter)
	if err != nil {
		fmt.Println("Error on get persons", err)
		return
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("%v.GetPersons(_) = _, %v", client, err)
		}
		fmt.Println("Person: ", person)
	}
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect: ", err)
		return
	}
	defer conn.Close()
	// Create a new PersonClient
	client := prsn.NewPersonClient(conn)
	person := &prsn.PersonRequest {
		Id: 1001,
		Name: "Reddy",
		Email: "reddv@xyz.com",
		Phone: "00000",
		Addresses: []*prsn.PersonRequest_Address{
			&prsn.PersonRequest_Address{
				Street: "Tripilcane",
				City: "Tirupati",
				State: "AP",
				Zip: "20002",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new Person
	createPerson(client, person)
	person = &prsn.PersonRequest{
		Id: 1002,
		Name: "mmm",
		Email: "reddv@xyz.com",
		Phone: "888888",
		Addresses: []*prsn.PersonRequest_Address{
			&prsn.PersonRequest_Address{
				Street: "Tripilcane",
				City: "Tirupati",
				State: "AP",
				Zip: "20002",
				IsShippingAddress: true,
			},
		},
	}

	//Create a new person
	createPerson(client, person)

	filter := &prsn.PersonFilter{Keyword: ""}
	getPersons(client, filter)
}