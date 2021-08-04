/* ProtocolBuffer
 */

package main

import (
	"encoding/json"
	"eventsgitn/contractsp"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8092", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Error de dial")
	}
	defer conn.Close()
	c := contractsp.NewEventsClient(conn)

	r, err := c.SearchAll(context.Background(), &contractsp.Empty{})
	if err != nil {
		log.Fatal("Error de call grpc ", err)
	}
	s, _ := json.MarshalIndent(&r, "", "    ")
	fmt.Println(string(s))

	q, err := c.SearchId(context.Background(), &contractsp.Id{Id: "60d28ba36efb600001647b75"})
	if err != nil {
		log.Fatal("Error de call grpc ", err)
	}
	s, _ = json.MarshalIndent(&q, "", "    ")
	fmt.Println(string(s))

	x, err := c.SearchName(context.Background(), &contractsp.Name{Name: "opera juana"})
	if err != nil {
		log.Fatal("Error de call grpc ", err)
	}
	s, _ = json.MarshalIndent(&x, "", "    ")
	fmt.Println(string(s))
}
