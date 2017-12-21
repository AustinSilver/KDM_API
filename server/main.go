package main

import (
	"fmt"
	"net"
	"strings"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	prsn "github.com/austinsilver/kdm_api/person"
)

const (
	port = ":3333"
)

// Person is used to implement prsn.PersonServer.

type Person struct {
	savedPersons []*prsn.Person.Requests
}

// CreatePerson creates a new Persson

func (p *Person) CreatePerson(c context.Context, input *prsn.PersonRequest) (*prsn.PersonResponses, error){
	p.savedPersons = append(p.savedPersons, input)

	return &prsn.PersonResponse{Id: input.Id, Success: true}, nil
}

// GetPersons returns all persons using filter

func (p *Person) GetPersons(fltr *prsn.PerosnFilter, stream prsn.Person_GetPersonsServer) error {

	for _, person := range p.savedPersons {

		if fltr.Keyword != "" {
			if !string.Contains(person.Name, fltr.Keyword) {
				continue
			}
		}

		err := stream.Send(person)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to listen: ", err)
		return
	}
	// Create a new gRPC person
	s := grpc.NewServer()
	prsn.RegisterPersonServer(s, &Person{})
	s.Serve(lis)
}