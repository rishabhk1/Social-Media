package main

import (
	// "fmt"
	// "log"
	// "net/http"
	// "os"
	// "rest-api/web"
	// "time"

	//"github.com/gorilla/mux"
	//"github.com/hyperledger/fabric-gateway/pkg/client"
	//"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	//"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"fmt"
	"log"
	"rest-api/web"
)

// type User struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// }

// type Community struct {
// 	ID          string    `json:"id"`
// 	Name        string    `json:"name"`
// 	Description string    `json:"description"`
// 	Creator     string    `json:"creator"`
// 	CreatedAt   time.Time `json:"createdAt"`
// 	Moderators  []string  `json:"moderators"`
// 	Posts       []Post    `json:"posts"`
// }

// type Post struct {
// 	ID        string    `json:"id"`
// 	Title     string    `json:"title"`
// 	Content   string    `json:"content"`
// 	Author    string    `json:"author"`
// 	Score     int       `json:"score"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	Comments  []Comment `json:"comments"`
// }

// type Comment struct {
// 	ID        string    `json:"id"`
// 	Content   string    `json:"content"`
// 	Author    string    `json:"author"`
// 	Score     int       `json:"score"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	Replies   []Comment `json:"replies"`
// }

// type OrgSetup struct {
// 	OrgName      string
// 	MSPID        string
// 	CryptoPath   string
// 	CertPath     string
// 	KeyPath      string
// 	TLSCertPath  string
// 	PeerEndpoint string
// 	GatewayPeer  string
// 	Gateway      client.Gateway
// }

// var err error
// var network *gateway.Network
// var contract *gateway.Contract

// func handleRequests() {
// 	r := mux.NewRouter().StrictSlash(true)
// 	r.HandleFunc("/community", getCommunity).Queries("name", "{name}").Methods("GET")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }

// func getCommunity(w http.ResponseWriter, r *http.Request) {
// 	log.Println("--> Endpoint hit: getCommunity")
// 	v := r.URL.Query()
// 	name := v.Get("name")
// 	log.Printf("%v", name)
// 	log.Printf("requesting community %v", name)
// 	communityJson, err := contract.EvaluateTransaction("GetCommunity", name)
// 	var community Community
// 	if err != nil {
// 		log.Printf("failed to evaluate chaincode: %v", err)
// 		w.WriteHeader(http.StatusNotFound)
// 	}
// 	err = json.Unmarshal(communityJson, &community)
// 	if err != nil {
// 		log.Printf("error marshaling/unmarshaling profile: %v", err)
// 	}
// 	communityJson, err = json.Marshal(community)
// 	if err != nil {
// 		log.Printf("error marshaling/unmarshaling profile: %v", err)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(communityJson)
// }

func main() {
	log.Println("============ Application starts ============")
	// network, err = configNet()
	// contract = network.GetContract("basic")
	// // _, err = contract.SubmitTransaction("InitLedger")
	// // if err != nil {
	// // 	log.Fatalf("error submitting chaincode transaction: %v", err)
	// // }
	// log.Println("============ End cn ============")
	// handleRequests()
	cryptoPath := "../../test-network/organizations/peerOrganizations/org1.example.com"
	orgConfig := web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
		KeyPath:      cryptoPath + "/users/User1@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}
	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		fmt.Println("Error initializing setup for Org1: ", err)
	}
	web.Serve(web.OrgSetup(*orgSetup))
}

// 	orgSetup, err := web.Initialize(orgConfig)
// 	if err != nil {
// 		fmt.Println("Error initializing setup for Org1: ", err)
// 	}
// 	web.Serve(web.OrgSetup(*orgSetup))
// }

// func configNet() (*gateway.Network, error) {
// 	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
// 	if err != nil {
// 		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
// 	}

// 	wallet, err := gateway.NewFileSystemWallet("wallet")
// 	if err != nil {
// 		log.Fatalf("Failed to create wallet: %v", err)
// 	}

// 	if !wallet.Exists("appUser") {
// 		err = populateWallet(wallet)
// 		if err != nil {
// 			log.Fatalf("Failed to populate wallet contents: %v", err)
// 		}
// 	}

// 	ccpPath := filepath.Join(
// 		"..",
// 		"..",
// 		"test-network",
// 		"organizations",
// 		"peerOrganizations",
// 		"org1.example.com",
// 		"connection-org1.yaml",
// 	)

// 	gw, err := gateway.Connect(
// 		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
// 		gateway.WithIdentity(wallet, "appUser"),
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to gateway: %v", err)
// 	}
// 	defer gw.Close()

// 	network, err := gw.GetNetwork("mychannel")
// 	if err != nil {
// 		log.Fatalf("Failed to get network: %v", err)
// 	}
// 	log.Println("============ End cn ============")
// 	return network, nil
// }

// // populateWallet populates a wallet in case it is not configuered already.
// func populateWallet(wallet *gateway.Wallet) error {
// 	log.Println("============ Populating wallet ============")
// 	credPath := filepath.Join(
// 		"..",
// 		"..",
// 		"test-network",
// 		"organizations",
// 		"peerOrganizations",
// 		"org1.example.com",
// 		"users",
// 		"User1@org1.example.com",
// 		"msp",
// 	)

// 	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
// 	// read the certificate pem
// 	cert, err := os.ReadFile(filepath.Clean(certPath))
// 	if err != nil {
// 		return err
// 	}

// 	keyDir := filepath.Join(credPath, "keystore")
// 	// there's a single file in this dir containing the private key
// 	files, err := os.ReadDir(keyDir)
// 	if err != nil {
// 		return err
// 	}
// 	if len(files) != 1 {
// 		return fmt.Errorf("keystore folder should have contain one file")
// 	}
// 	keyPath := filepath.Join(keyDir, files[0].Name())
// 	key, err := os.ReadFile(filepath.Clean(keyPath))
// 	if err != nil {
// 		return err
// 	}

// 	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
// 	log.Println("============ Populating wallet ============")
// 	return wallet.Put("appUser", identity)
// }

// // // hashID hashes the data and returns the hash as ID
// // func hashTxn(data string) string {
// // 	id_bytes := sha256.Sum256([]byte(data))
// // 	ID := hex.EncodeToString(id_bytes[:])
// // 	return ID
// // }

// // // verifyTxn verifies connection between a Transaction on chain and its data in database
// // func verifyTxn(txid string, contract *gateway.Contract) {
// // 	log.Printf("--> Evaluate Transaction!")
// // 	result, err := contract.EvaluateTransaction("ReadPost", txid)
// // 	if err != nil {
// // 		log.Fatalf("Failed to evaluate transaction: %v", err)
// // 	}
// // 	log.Printf("Transaction %s, is verfied!\n", string(result))
// // }
