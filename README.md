Packages used:

    gorilla/mux

Source files:

    decryptor/decryptor.go - Cipher-related functions
    freqanalysis/freqanalysis.go - Letter-frequency functions
    msdecryptor/main.go - Decryptor service
    msrestapi/main.go - REST API service

Commands used:

    protoc -I pbdecryptor/ pbdecryptor/pbdecryptor.proto --go_out=plugins=grpc:pbdecryptor
 
Test commands:

    curl -i -H "Content-Type: multipart/form-data" -X POST http://localhost:8083/decrypt/1 -F "data=@encrypted.txt"

    curl -i -H "Content-Type: multipart/form-data" -X POST http://localhost:8083/decrypt/2 -F "data=@encrypted_hard.txt"
