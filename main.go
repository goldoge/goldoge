package main

import (
	"encoding/json"
	"github.com/goldoge/goldoge/gd"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Star struct {
	Address   string
	Message   string
	Signature string
	Star      string
}

type Wallet struct {
	Address string
}

func main() {
	router := mux.NewRouter()

	blockchain := gd.BlockChain{
		Chain:  make([]gd.Block, 0),
		Height: -1,
	}
	blockchain.InitializeChain()

	// getBlockByHeight
	router.HandleFunc("/blocks/height/{height}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			vars := mux.Vars(r)
			height, ok := vars["height"]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
			}
			heightInt64, err := strconv.ParseInt(height, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
			block, errors := blockchain.GetBlockByHeight(heightInt64)
			if errors != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			blockJson, errors := json.Marshal(block)
			if errors != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(blockJson)
		}
	})

	// requestOwnership
	router.HandleFunc("/requestOwnership", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var newWallet Wallet
			err := json.NewDecoder(r.Body).Decode(&newWallet)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			message := blockchain.RequestMessageOwnershipVerification(newWallet.Address)
			w.Header().Set("Content-Type", "application/json")
			messageJson, errors := json.Marshal(message)
			if errors != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(messageJson)
			return
		}
	})

	// submitstar
	router.HandleFunc("/submitstar", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var newStar Star
			err := json.NewDecoder(r.Body).Decode(&newStar)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			block, errors := blockchain.SubmitStar(newStar.Address, newStar.Message, newStar.Signature, newStar.Star)
			if errors != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			blockJson, errors := json.Marshal(block)
			if errors != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(blockJson)
			return
		}
	})

	// getBlockByHash
	router.HandleFunc("/blocks/hash/{hash}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			vars := mux.Vars(r)
			hash, ok := vars["hash"]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
			}
			block, errors := blockchain.GetBlockByHash(hash)
			if errors != nil {
				w.WriteHeader(http.StatusNotFound)
			}
			blockJson, errors := json.Marshal(block)
			if errors != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(blockJson)
		}
	})

	// getStarsByOwner (address)
	router.HandleFunc("/blocks/{address}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			vars := mux.Vars(r)
			address, ok := vars["address"]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
			}
			stars := blockchain.GetStarByWalletAddress(address)
			starsJson, errors := json.Marshal(stars)
			if errors != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(starsJson)
		}
	})

	// validateChain
	router.HandleFunc("/validateChain", func(w http.ResponseWriter, r *http.Request) {
		type Result struct {
			block gd.Block
			error bool
		}
		var errorBlocks []Result
		for _, block := range blockchain.Chain {
			if !block.Validate() {
				errorBlocks = append(errorBlocks, Result{
					block,
					false,
				})
			}
		}

		errorBlocksJSON, errors := json.Marshal(errorBlocks)
		if errors != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(errorBlocksJSON)
	})

	addr := ":5000"
	log.Println("Listen on", addr)
	log.Fatal(http.ListenAndServe(addr, router))

}
