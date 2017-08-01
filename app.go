package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	cid "gx/ipfs/QmTprEaAA2A9bst5XH7exuyi5KzNMK3SEDNN8rBDnKWcUS/go-cid"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"

	"github.com/ipfs/go-ipfs/core"
	"github.com/julienschmidt/httprouter"
	"github.com/sheremetat/dapp-ajs/common"
)

//struct to handle the IPFS core node
type IPFSHandler struct {
	node *core.IpfsNode
}

//struct to pass vars to front-end
type DemoPage struct {
	Title   string
	Author  string
	Tweet   []string
	isMine  bool
	Balance float64
}

//struct for list of all peers in IPFS
type PeerList struct {
	Allpeers []string
	Balance  float64
}

func main() {

	node, err := common.StartNode() //ipfs.StartNode()
	if err != nil {
		panic(err)
	}

	//[2] Define routes
	router := httprouter.New()

	//Route 1 Home (profile)
	router.GET("/", TextInput(node))
	//Route 2 Discover page
	router.GET("/discover", displayUsers(node))
	//Route 3 Other user profiles
	router.GET("/profile/:name", TextInput(node))
	//Route 4 Add text to IPFS
	router.POST("/textsubmitted", addTexttoIPFS(node))

	//[3] link resources
	router.ServeFiles("/resources/*filepath", http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.Handle("/", router)

	//[4] Start server
	fmt.Println("serving at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

//Called on the discovery page. It will create the list of
//all IPFS peers currently online
func displayUsers(node *core.IpfsNode) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		//get peers
		peers := node.Peerstore.Peers()
		data := make([]string, len(peers))
		for i := range data {
			data[i] = peer.IDB58Encode(peers[i])
		}

		//send the peer list to the front end template
		demoheader := PeerList{Allpeers: data, Balance: 0}
		fp := path.Join("templates", "discover.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, demoheader); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

//Called on all profile pages. Fills the profile page with tweets for the relevant user.
func TextInput(node *core.IpfsNode) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var userID = ps.ByName("name")

		//get your current balance

		//[1] If its your home profile page
		if userID == "" {
			pointsTo, err := common.GetDAG(node, node.Identity.Pretty())
			cid, _ := cid.Parse(pointsTo.String())
			tweetArray, err := common.GetStrings(node, cid.String())

			if err != nil {
				fmt.Println("WHOOPS", err)
			}
			fmt.Println("the tweet array is %s", tweetArray)
			//[1A] If no tweets, send nil
			if tweetArray == nil {
				fmt.Println("tweetarray is nil")
				demoheader := DemoPage{"Decentralized Twitter", "SR", nil, true, 0}
				fp := path.Join("templates", "index.html")
				tmpl, err := template.ParseFiles(fp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := tmpl.Execute(w, demoheader); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				fmt.Println("tweetarray is not nil")
				//[1B] If tweets, send tweet array
				demoheader := DemoPage{"Decentralized Twitter", "SR", tweetArray, true, 0}
				fp := path.Join("templates", "index.html")
				tmpl, err := template.ParseFiles(fp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := tmpl.Execute(w, demoheader); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

		} else {

			//[2] If its another profile
			//Pull from IPNS
			pointsTo, err := common.GetDAG(node, userID)

			//[2A] If nil, send nil
			if err != nil {
				fmt.Println("ERROR")
				fmt.Println("tweetarray is nil")
				demoheader := DemoPage{"Decentralized Twitter", "SR", nil, false, 0}
				fp := path.Join("templates", "index.html")
				tmpl, err := template.ParseFiles(fp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := tmpl.Execute(w, demoheader); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				//[2B] Else send tweetarray
				fmt.Println("RESOLVED")
				//else pull it from the URL
				tweetArray, err := common.GetStrings(node, pointsTo.String())
				if err != nil {
					panic(err)
				}
				demoheader := DemoPage{"Decentralized Twitter", "Siraj", tweetArray, false, 0}
				fp := path.Join("templates", "index.html")
				tmpl, err := template.ParseFiles(fp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := tmpl.Execute(w, demoheader); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

		}

	}

}

//Called when user submits text to IPFS.
func addTexttoIPFS(node *core.IpfsNode) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r.ParseForm()
		fmt.Println("input text is:", r.Form["sometext"])
		var userInput = r.Form["sometext"]
		Key, err := common.AddString(node, userInput[0])
		if err != nil {
			panic(err)
		}

		fmt.Println("the key", Key)

	}
}
