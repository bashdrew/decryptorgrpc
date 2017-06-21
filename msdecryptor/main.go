/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	dec "bashdrew/bsscodingassignment/decryptor"
	fa "bashdrew/bsscodingassignment/freqanalysis"
	pb "bashdrew/bsscodingassignment/pbdecryptor"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var database *sql.DB
var bssDecryptOrder []dec.DecryptInfoList

// server is used to implement pbdecryptor Methods.
type server struct{}

func getFileContent(inFile string) (result string) {
	b, err := ioutil.ReadFile(inFile) // just pass the file name
	if err != nil {
		log.Fatalf("File read error: %v", err)
	}

	result = string(b)

	return result
}

// Decrypt implements PBDecryptor.Decrypt
func (s *server) Decrypt(ctx context.Context, in *pb.DecryptRequest) (*pb.DecryptResponse, error) {
	decryptResp := new(pb.DecryptResponse)

	encTxt := in.EncText
	dicTxt := getFileContent("../data/plain.txt")
	plnTxt, cipherKey := bssDecryptOrder[in.Id-1].Decrypt(encTxt, dicTxt)

	decryptResp.Id = in.Id
	decryptResp.PlnText = string(plnTxt)
	decryptResp.CipherKey = cipherKey

	return decryptResp, nil
}

func init() {
	bssDecryptOrder = []dec.DecryptInfoList{
		{
			{FreqFunc: fa.GetWordFreq, LtrCnt: 1, EntTop: 3},
			{FreqFunc: fa.GetWordFreq, LtrCnt: 2, EntTop: 3},
			{FreqFunc: fa.GetLetterFreqMulti, LtrCnt: 2, EntTop: 2},
			{FreqFunc: fa.GetLetterFreqMulti, LtrCnt: 1, EntTop: -1},
		},
		{
			{FreqFunc: fa.GetWordFreq, LtrCnt: 2, EntTop: 10},
			{FreqFunc: fa.GetWordFreq, LtrCnt: 1, EntTop: 3},
			{FreqFunc: fa.GetLetterFreqMulti, LtrCnt: 1, EntTop: -1},
		},
	}

	fmt.Println("Decryptor Microservice started...")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPBDecryptorServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
