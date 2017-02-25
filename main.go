package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/matobet/shabang/db"
	"github.com/matobet/shabang/model"
)

var bitlen = flag.Uint("bitlen", 32, "collision prefix bit length")
var seed = flag.String("seed", "foobarbaz", "seed string")
var dbPath = flag.String("dbpath", "hash.db", "path to the hash store")

type PersistRequest struct {
	Hash, Value model.HashBytes
}

var totalHashes int64 = 0

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *bitlen > 8*sha256.Size {
		log.Fatal("Prefix bit length cannot be longer than the whole hash size!")
	}

	hashDb, err := db.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	hashDb = db.Bloomify(hashDb, 1000000, 3)
	defer hashDb.Close()

	persistRequests := make(chan PersistRequest, 10000000)
	quit := make(chan struct{})

	startTime := time.Now()

	go hasher(hashDb, persistRequests, quit)
	go persister(hashDb, persistRequests, quit)

	<-quit

	elapsed := time.Since(startTime)
	fmt.Printf("Checked %d hashes in %s @ %f H/s\n", totalHashes, elapsed, float64(totalHashes)/elapsed.Seconds())
}

func hasher(db db.HashDB, persistRequests chan PersistRequest, quit chan struct{}) {
	ctx := sha256.New()

	first := model.HashBytes(*seed)
	var next model.HashBytes

	count := 0

	for {
		first.Trim(*bitlen)
		next = first.Sum(ctx)
		next.Trim(*bitlen)
		if preImage, err := db.Check(next); err != nil || preImage != nil {
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Collision detected!!!!! %x hashes to the same '%x' as %x\n", first, next, preImage)
			close(quit)
			break
		}
		persistRequests <- PersistRequest{Hash: next, Value: first}
		first = next
		count++
	}

	totalHashes = int64(count)
}

func persister(db db.HashDB, requests <-chan PersistRequest, quit chan struct{}) {
	for {
		select {
		case req := <-requests:
			err := db.Write(req.Hash, req.Value)
			if err != nil {
				log.Fatal(err)
			}
		case <-quit:
			return
		}
	}
}
