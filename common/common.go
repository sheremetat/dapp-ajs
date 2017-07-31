package common

import (
	"context"
	"fmt"
	cid "gx/ipfs/QmTprEaAA2A9bst5XH7exuyi5KzNMK3SEDNN8rBDnKWcUS/go-cid"

	"github.com/ipfs/go-ipfs/core"
	merkledag "github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/path"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
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
		Repo: r}
	node, err := core.NewNode(ctx, cfg)

	if err != nil {
		return nil, err
	}
	return node, nil
}

func GetStrings(node *core.IpfsNode, path path.Path) ([]string, error) {
	var tweetArray = resolveAllInOrder(node, path)
	return tweetArray, nil
}

func GetDAG(node *core.IpfsNode, id string) (path.Path, error) {
	pointsTo, err := node.Namesys.Resolve(node.Context(), id)
	return pointsTo, err
}

func resolveAllInOrder(nd *core.IpfsNode, path path.Path) []string {
	var stringArr []string
	//var node *merkledag.RawNode

	// node, err := nd.DAG.Get(nd.Context(), path)
	// fmt.Printf("the node is", node)
	// if err != nil {
	// 	fmt.Println(err)
	// 	//return
	// }
	// fmt.Printf("bout to crash")
	// fmt.Printf("%s ", string(node.RawData()[:]))
	// fmt.Println("not crashed ")

	// for {
	// 	var err error

	// 	if len(node.Links()) == 0 {
	// 		break
	// 	}

	// 	node, err := node.Links()[0].GetNode(nd.Context(), nd.DAG)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	fmt.Printf("%s ", string(node.RawData()[:]))
	// 	stringArr = append(stringArr, string(node.RawData()[:]))
	// }

	// fmt.Printf("\n")

	return stringArr

}

func AddString(node *core.IpfsNode, inputString string) (*cid.Cid, error) {
	return nil, fmt.Errorf("Not supported yet!")

	//pointsTo, err := node.Namesys.Resolve(node.Context(), node.Identity.Pretty())
	//If there is an error, user is new and hasn't yet created a DAG.
	// if err != nil {
	// 	//[3] Initialize a MerkleDAG node and key
	// 	var NewNode *merkledag.RawNode
	// 	//[4] Fill the node with user input
	// 	NewNode = makeStringNode(inputString)
	// 	//[5] Add the node to IPFS
	// 	nodeCid, _ := node.DAG.Add(NewNode.Copy())
	// 	// //publish to IPNS
	// 	output, err := commands.Publish(node, node.PrivateKey, nodeCid.B58String())
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println("You published to IPNS. Your peer ID is ", output.Name)
	// 	}

	// 	return nodeCid, nil

	// } else {
	// 	//[7] Initialize a new MerkleDAG node and key
	// 	//var NewNode *merkledag.RawNode
	// 	//[8] Fill the node with user input
	// 	NewNode := makeStringNode(inputString)
	// 	//[10] Convert it into a key
	// 	Key, _ := cid.Parse(pointsTo.String())
	// 	//[11] Get the Old MerkleDAG node and key
	// 	//var OldNode *merkledag.RawNode
	// 	OldNode, _ := node.DAG.Get(node.Context(), Key)
	// 	//[12]Add a link to the old node
	// 	NewNode = NewNode.AddNodeLink("next", OldNode)

	// 	//[13] Add thew new node to IPFS
	// 	Key2, _ := node.DAG.Add(NewNode)
	// 	// //publish to IPNS
	// 	output, err := commands.Publish(node, node.PrivateKey, Key2.B58String())
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println("You published to IPNS. Your peer ID is ", output.Name)
	// 	}
	// 	return Key2, nil
	// }

}

func makeStringNode(s string) *merkledag.RawNode {
	data := []byte(s)
	return merkledag.NewRawNode(data)
}
