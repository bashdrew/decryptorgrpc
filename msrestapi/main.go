package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	pb "bashdrew/bsscodingassignment/pbdecryptor"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
)

const (
	port             = ":8083"
	decryptorAddress = "localhost:50051"
	maxMemory        = 1 * 1024 * 1024
)

const (
	decryptID     = "type"
	decryptedFile = "../data/decrypted.txt"
)

func getDecryptorConn() (conn *grpc.ClientConn, client pb.PBDecryptorClient, err error) {
	conn, err = grpc.Dial(decryptorAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client = pb.NewPBDecryptorClient(conn)

	return
}

func getFileContent(w http.ResponseWriter, r *http.Request) (result string) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			buf, _ := ioutil.ReadAll(file)
			result = result + string(buf)
		}
	}

	return result
}

func saveDecrytedFile(outFile string, dataFile []byte) {
	err := ioutil.WriteFile(outFile, dataFile, 0644)
	if err != nil {
		log.Fatalf("File write error: %v", err)
	}
}

// PostDecryptEndpoint ... decrypt a text
func PostDecryptEndpoint(w http.ResponseWriter, req *http.Request) {
	var decReply *pb.DecryptResponse
	var decRequest *pb.DecryptRequest

	// Set up a connection to the server.
	conn, c, err := getDecryptorConn()
	defer conn.Close()
	if err == nil {
		params := mux.Vars(req)
		decRequest = new(pb.DecryptRequest)
		decRequest.EncText = getFileContent(w, req)
		id64, _ := strconv.ParseInt(params[decryptID], 10, 64)
		decReply, err = c.Decrypt(context.Background(),
			&pb.DecryptRequest{
				Id:      id64,
				EncText: decRequest.EncText,
			})
		if err != nil {
			log.Fatalf("could not decrypt text: %v", err)
		} else {
			saveDecrytedFile(decryptedFile, []byte(decReply.PlnText))

			http.ServeFile(w, req, decryptedFile)
		}
	}

	json.NewEncoder(w).Encode(decReply)
}

func setupRestAPIConn() (conn *grpc.ClientConn, err error) {
	// Set up a connection to the server
	conn, err = grpc.Dial(decryptorAddress, grpc.WithInsecure())

	return
}

func init() {
	fmt.Println("REST API Microservice started...")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/decrypt", PostDecryptEndpoint).Methods("POST")
	router.HandleFunc("/decrypt/{"+decryptID+"}", PostDecryptEndpoint).Methods("POST")

	log.Fatal(http.ListenAndServe(port, router))
}
