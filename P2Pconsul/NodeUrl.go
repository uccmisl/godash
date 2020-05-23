package P2Pconsul

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	hlpr "github.com/uccmisl/godash/P2Pconsul/HelperFunctions"
	pb "github.com/uccmisl/godash/P2Pconsul/P2PService"

	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/uccmisl/godash/logging"
)

//NodeUrl - Collaborative Code - Start
type NodeUrl struct {
	//NodeUrl - node variables
	ClientName      string
	url             string
	previousUrl     map[string]string
	Addr            string
	ContentPort     string
	ContentLocation string
	Clients         map[string]string
	IP              net.IP
	Registered      bool

	update bool

	debug bool

	//consul variables
	SDAddress string
	SDKV      api.KV

	//debug details
	Debugfile string
	Debuglog  bool

	//Server for implementation
	pb.UnimplementedP2PServiceServer
}

// Initialisation -
func (n *NodeUrl) Initialisation() {
	//init required varibales
	n.Clients = make(map[string]string)
	n.Registered = false
	n.previousUrl = make(map[string]string)
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(63000) + 1023
	n.GetOutboundIP()
	n.Addr = n.IP.String() + ":" + strconv.Itoa(port)
	n.debug = false

	s := fmt.Sprintf("addr : %v\n", n.Addr)
	n.DebugPrint(s)
	s = fmt.Sprintf(" Content addr : %v\n", n.ContentPort)
	n.DebugPrint(s)
	s = fmt.Sprintf(" Content Location :%v\n", n.ContentLocation)
	n.DebugPrint(s)
	s = fmt.Sprintf("IP ADREESS:%v\n", n.IP)
	n.DebugPrint(s)
	//start server listening

	n.RegisterNode()
	//console input that takes urls for searchingfor in n.Clients
	/*fmt.Printf("Input URL for search:\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("line text : %v\n", line)
		fmt.Printf("Returned URL : %v\n", n.Search(line))
		fmt.Printf("Client list: \n")
		for key, client := range n.Clients {
			fmt.Printf("key: %v client: %v\n", key, client)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}*/

}

// StartListening Start the node server listening
func (n *NodeUrl) StartListening(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", n.Addr)
	s := fmt.Sprintf("GRPC Server Listening on %v\n", n.Addr)
	n.DebugPrint(s)
	if err != nil {
		log.Fatalf("failed to start listening %v", err)
	}

	_n := grpc.NewServer()

	pb.RegisterP2PServiceServer(_n, n)

	reflection.Register(_n)
	defer _n.Stop()
	if err := _n.Serve(lis); err != nil {
		// fmt.Println("server failed")
		// log.Fatalf("Server Failed to serve")
		n.DebugPrint("server failed")
	}
	wg.Done()
}

// RegisterNode First register node under a given URL
func (n *NodeUrl) RegisterNode() {
	config := api.DefaultConfig()
	config.Address = n.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Panicln("Unable to register with KV Service Discovery")
	}

	//create Key Value store on consul server
	//Key is (URL+NodeAddress : NodeAddress)
	kv := consul.KV()

	//Store KV for later use
	n.SDKV = *kv
	// fmt.Printf("Successfully registered : (%v : %v)\n", hlpr.HashSha(n.url)+n.Addr, []byte(n.Addr))
	s := fmt.Sprintf("Successfully registered : (%v : %v)\n", hlpr.HashSha(n.url)+n.Addr, []byte(n.Addr))
	n.DebugPrint(s)
}

//Search search network for a given url
func (n *NodeUrl) Search(url string) string {
	start := time.Now()
	n.DebugPrint("in consul search url :" + url)
	notFound := true
	l := strings.Split(url, "/")
	location := l[len(l)-1]
	n.update = true

	key := hlpr.HashSha(url)

	n.DebugPrint("Start of consul search Location:" + location)

	//if desired content is not in current clients
	//search clients of clients
	if len(n.Clients[key]) != 0 {
		//if current client is known to have correct content from previous requests
		contentServer, err := n.GetContentServerAddress(n.Clients[key])
		if err != nil {
			return url
		}
		url = "http://" + contentServer + "/" + location + "::localclient"
		notFound = false
	}

	//loop over all known nodes
	for _, client := range n.Clients {
		n.DebugPrint("Looping client check")
		//establish connection to client and check for content
		conn, err := grpc.Dial(client, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect to server : %v", err)
			break
		}

		defer conn.Close()
		c := pb.NewP2PServiceClient(conn)

		//GRPC call to check clients for content
		downloadAddress, err := c.CheckClients(context.Background(), &pb.CheckRequest{Address: n.Addr, Target: key})
		//fmt.Printf("check %v\n",err)
		//fmt.Printf("downloadAddress client loop: %v\n", downloadAddress.Addr)

		if err != nil {
			n.DebugPrint("rpc error client check")

			continue
		}

		if downloadAddress.Addr != "nil" {
			//add relevant client to client list
			n.Clients[key] = downloadAddress.Addr

			//get content server address for Url for download
			contentServer, err := n.GetContentServerAddress(downloadAddress.Addr)
			if err != nil {
				n.DebugPrint("rpc error 2 client check")
				break
			}
			url = "http://" + contentServer + "/" + location + "::clients"

			notFound = false
			break
		}
	}

	//in Case not currently known Locally consult Consul Server7
	if notFound {
		n.DebugPrint("checking consul")
		kvpairs, _, err := n.SDKV.List(key, nil)
		/* (len(kvpairs)>4){
			n.update = false
		}else{
			n.update = true
		}*/
		n.DebugPrint("checking consul too")
		if err != nil {
			n.DebugPrint("consul error")
			return url
		}
		//Loop Key Value pair matches query
		//randomly shuffle key value pairs
		n.DebugPrint("checking consul keys")
		for i := 1; i < len(kvpairs); i++ {
			r := rand.Intn(i + 1)
			if i != r {
				kvpairs[r], kvpairs[i] = kvpairs[i], kvpairs[r]
			}
		}
		for _, kventry := range kvpairs {
			//Check key isnt this node
			n.DebugPrint("Looping consul entries")
			if kventry.Key[0:len(key)] == key && kventry.Key != key+n.Addr {

				//Add random pick for which node to download from

				//add relevant node to clients
				n.Clients[kventry.Key[0:len(key)]] = string(kventry.Value)

				contentServer, err := n.GetContentServerAddress(string(kventry.Value))
				if err != nil {
					log.Fatalf("Error Kventry : %v", err)
					fmt.Println("KV ERROR")
					break
				} else {
					url = "http://" + contentServer + "/" + location + "::consul"
					fmt.Println("location" + location)
					//fmt.Printf("download content from consul\n")
					notFound = false
					break
				}
			}
		}
	}
	return url + "::" + time.Since(start).String()
}

//UpdateConsul consul reference to this node
// updates nodes URL references also
func (n *NodeUrl) UpdateConsul(url string) {
	//add new consul entry

	n.DebugPrint(fmt.Sprintf("consul Update : %v\n", url+n.Addr))
	p := &api.KVPair{Key: url + n.Addr, Value: []byte(n.Addr)}
	_, err := n.SDKV.Put(p, nil)
	fmt.Println("updating consul ###############################################")
	n.DebugPrint(fmt.Sprintf("error update consul %v\n", err))
	if err != nil {
		n.DebugPrint("error update consul")
		fmt.Println("issue")
	}
	//fmt.Printf("new consul entry created\n")
	//update nodes url references
	n.previousUrl[n.url] = n.url
	n.url = url

	//fmt.Printf("Node Url references updated\n")
}

func (n *NodeUrl) getContentAddr(address string) (serverAddr string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server : %v\n", err)
	}
	defer conn.Close()
	c := pb.NewP2PServiceClient(conn)

	if err != nil {
		log.Fatalf("Did not connect to server : %v\n", err)
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	downloadAddress, err := c.GetServerAddr(ctx, &pb.ServerRequest{Address: n.Addr})
	return downloadAddress.Addr

}

//CheckClients functions takes requests from other node
// to check their clients for desired content
func (n *NodeUrl) CheckClients(ctx context.Context, in *pb.CheckRequest) (*pb.CheckReply, error) {
	fmt.Printf("\nCheck Client Target: %v\n", n.Clients[in.Target])
	response := pb.CheckReply{Addr: "nil"}
	if len(n.previousUrl[in.Target]) > 0 {
		response = pb.CheckReply{Addr: n.Addr}
		fmt.Printf("N addr: %v\n", n.Addr)
		return &response, nil
	}
	if len(n.Clients[in.Target]) > 0 {
		response = pb.CheckReply{Addr: n.Clients[in.Target]}
	}
	return &response, nil
	//add second check client here

}

//SecondCheckLoop -
func (n *NodeUrl) SecondCheckLoop(url string) (addr string) {
	for _, client := range n.Clients {
		conn, err := grpc.Dial(client, grpc.WithInsecure())
		if err != nil {
			continue
			log.Fatalf("Did not connect to server : %v", err)
		}
		defer conn.Close()
		c := pb.NewP2PServiceClient(conn)

		response, err := c.SecondCheckClient(context.Background(), &pb.SecondCheckRequest{Url: url})
		if err != nil {
			fmt.Printf("Error in second check client%v\n", err)
			return "nil"
		}
		if response.Addr != "nil" {
			return response.Addr
		}
	}
	return "nil"
}

//SecondCheckClient -
func (n *NodeUrl) SecondCheckClient(ctx context.Context, in *pb.SecondCheckRequest) (*pb.SecondCheckReply, error) {
	if len(n.previousUrl[in.Url]) > 0 || n.url == in.Url {
		response := pb.SecondCheckReply{Addr: n.Addr}
		return &response, nil
	}
	response := pb.SecondCheckReply{Addr: "nil"}
	return &response, nil
	//figure out broken flow for checking clients of clients
}

//GetServerAddr -
func (n *NodeUrl) GetServerAddr(ctx context.Context, in *pb.ServerRequest) (*pb.ServerRequestReply, error) {
	response := pb.ServerRequestReply{Addr: n.IP.String() + n.ContentPort}
	return &response, nil
}

//GetContentServerAddress -
func (n *NodeUrl) GetContentServerAddress(address string) (string, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server : %v\n", err)
		fmt.Println(err)
	}

	s := pb.NewP2PServiceClient(conn)

	if err != nil {
		log.Fatalf("Did not connect to server : %v\n", err)
		fmt.Println(err)
	}

	downloadAddress, err := s.GetServerAddr(context.Background(), &pb.ServerRequest{Address: address})
	fmt.Printf("download address %v\n", address)
	if err != nil {
		fmt.Println(err)
		return "nil", err
		//log.Fatalf("Error in check clients\nerr : %v\n", err)
	}
	url := downloadAddress.Addr
	return url, nil
}

//ContentServerStart -
func (n *NodeUrl) ContentServerStart(location string, port string, wg *sync.WaitGroup) {
	server := http.NewServeMux()

	//handlers that serve the home html file when called
	fs := http.FileServer(http.Dir(location))

	//os := http.FileServer(http.Dir("./"))

	//handles paths by serving correct files
	//there will be if statements down here that check if someone has won or not soon
	server.Handle("/", fs)
	//server.Handle("/"+n.ClientName, os)

	//logs that server is Listening
	s := fmt.Sprintf("Listening... %v\n", location)
	n.DebugPrint(s)
	//starts server
	http.ListenAndServe(n.IP.String()+port, server)
	wg.Done()
}

//GetOutboundIP -
func (n *NodeUrl) GetOutboundIP() {
	conn, err := net.Dial("udp", "10.0.0.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	n.IP = conn.LocalAddr().(*net.UDPAddr).IP

}

//SetDebug -
func (n *NodeUrl) SetDebug(DebugFile string, DebugLog bool) {
	n.Debugfile = DebugFile
	n.Debuglog = DebugLog
	n.debug = true
}

//DebugPrint -
func (n *NodeUrl) DebugPrint(s string) {
	if n.debug {
		logging.DebugPrint(n.Debugfile, n.Debuglog, "\nDEBUG: ", s)
	}
}

/*func main() {
	noden := NodeUrl{SDAddress: "127.0.0.1:8500", Clients: nil} // noden is for opeartional purposes
	noden.Initialisation()
	//start server for downloading content from
}*/
// Collaborative Code - End
