package common

import (
	"context"
	"fmt"
	cid "gx/ipfs/QmTprEaAA2A9bst5XH7exuyi5KzNMK3SEDNN8rBDnKWcUS/go-cid"

	"github.com/ipfs/go-ipfs/core"
	merkledag "github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/path"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	key "github.com/ipfs/go-key"
)

func StartNode() (*core.IpfsNode, error) {
	// Make our 'master' context and defer cancelling it
	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	r, err := fsrepo.Open("~/.ipfs")
	if err != nil {
		return nil, err
	}

	//[1] init IPFS node
	cfg := &core.BuildCfg{Online: true,
		Repo: r,
	}
	node, err := core.NewNode(ctx, cfg)

	if err != nil {
		return nil, err
	}
	//fmt.Printf("the node is", node.DAG)
	return node, nil
}

func GetStrings(node *core.IpfsNode, userId string) ([]string, error) {
	fmt.Println("the userId is", userId)
	key := key.B58KeyDecode(userId)
	fmt.Println("the key is", key)
	var tweetArray = []string{} //resolveAllInOrder(node, key)
	return tweetArray, nil
}

func GetDAG(node *core.IpfsNode, id string) (path.Path, error) {
	pointsTo, err := node.Namesys.Resolve(node.Context(), id)
	return pointsTo, err
}

func resolveAllInOrder(nd *core.IpfsNode, key key.Key) []string {
	var stringArr []string
	//var node *merkledag.RawNode
	fmt.Println("the key is", key.B58String())
	cid, _ := cid.Parse(key.B58String())
	fmt.Printf("the cid is", cid.String())
	node, err := nd.DAG.Get(nd.Context(), cid)
	fmt.Printf("the node is", node)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Printf("bout to crash")
	fmt.Printf("%s ", string(node.RawData()[:]))
	fmt.Println("not crashed ")

	for {
		var err error

		if len(node.Links()) == 0 {
			break
		}

		node, err := node.Links()[0].GetNode(nd.Context(), nd.DAG)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s ", string(node.RawData()[:]))
		stringArr = append(stringArr, string(node.RawData()[:]))
	}

	fmt.Printf("\n")

	return stringArr

}

func AddString(node *core.IpfsNode, inputString string) (*cid.Cid, error) {
	pointsTo, err := node.Namesys.Resolve(node.Context(), node.Identity.Pretty())
	fmt.Println("pointsTo on start add string: ", pointsTo, err)
	//If there is an error, user is new and hasn't yet created a DAG.
	if err != nil {
		//[3] Initialize a MerkleDAG node and key
		var NewNode *merkledag.ProtoNode
		//[4] Fill the node with user input
		NewNode = makeStringNode(inputString)
		//[5] Add the node to IPFS
		nodeCid, _ := node.DAG.Add(NewNode)
		fmt.Println("Node DAG add with CID", nodeCid)

		//publish to IPNS
		err := node.Namesys.Publish(node.Context(), node.PrivateKey, pointsTo)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("You published to IPNS. Your peer ID is ")
		}

		return nodeCid, nil

	} else {
		//[7] Initialize a new MerkleDAG node and key
		var NewNode *merkledag.ProtoNode
		//[8] Fill the node with user input
		NewNode = makeStringNode(inputString)
		//[10] Convert it into a key
		cids, _ := cid.Parse(pointsTo.String())
		//[11] Get the Old MerkleDAG node and key
		//var OldNode *merkledag.RawNode
		OldNode, err := node.DAG.Get(node.Context(), cids)
		//[12]Add a link to the old node
		NewNode.AddNodeLink("next", OldNode)
		//[13] Add thew new node to IPFS
		Key2, _ := node.DAG.Add(NewNode)
		// //publish to IPNS
		err = node.Namesys.Publish(node.Context(), node.PrivateKey, pointsTo)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("You published to IPNS. Your peer ID is ", Key2)
		}
		return Key2, nil
	}

}

func makeStringNode(s string) *merkledag.ProtoNode {
	data := []byte(s)
	return merkledag.NodeWithData(data)
}
