package main

import (
	"flag"
	"fmt"
        "strings"
        "strconv"
        "os"
        "bufio"

	"github.com/gologme/log"
	"github.com/cerisara/activityserve"
)

var actor activityserve.Actor

func cli() {
        reader := bufio.NewReader(os.Stdin)
        for ;; {
            fmt.Print("-> ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            if strings.HasPrefix(text,"post ") {
                post(text[5:])
            } else if strings.Compare(text,"following") == 0 {
                fl := actor.Following()
                fmt.Printf("following: %v\n",fl)
            } else if strings.Compare(text,"followers") == 0 {
                fl := actor.Followers()
                fmt.Printf("followers: %v\n",fl)
            } else if strings.HasPrefix(text,"follow ") {
                follow(text[7:])
            } else if strings.Compare("quit", text) == 0 {
                break
            }
        }
}

func post(s string) {
        fmt.Println("posting "+s)
	actor.CreateNote(s, "")
}

func follow(u string) {
	// actor.Follow("https://mastodon.etalab.gouv.fr/@cerisara")
        fmt.Println("following "+u)
	actor.Follow(u)
}

func gotmsg(o map[string]interface{}) {
    fmt.Printf("INCOMING MSG: FROM %v\n",o["attributedTo"])
    fmt.Printf("%v\n",o["content"])
    /*
gokey attributedTo http://actpub.duckdns.org/detson
gokey cc http://actpub.duckdns.org/detson/followers
gokey content toto est beau
gokey id http://actpub.duckdns.org/detson/item/h2V5X80ZLmy7rUYZ
gokey published 2020-02-23T11:44:20+01:00
gokey to https://www.w3.org/ns/activitystreams#Public
gokey type Note
gokey url http://actpub.duckdns.org/detson/item/h2V5X80ZLmy7rUYZ
    */
}

func main() {

	debugFlag := flag.Bool("debug", false, "set to true to get debugging information in the console")
	flag.Parse()

	if *debugFlag == true {
		log.EnableLevel("error")
	} else {
		log.DisableLevel("info")
	}

	activityserve.Setup("config.ini", *debugFlag)

        // get the port and actor name also here
        file, err := os.Open("config.ini")
        if err != nil { log.Fatal(err) }
        defer file.Close()
        scanner := bufio.NewScanner(file)
        var userag string = "newact"
        var port int = 8081
        for scanner.Scan() {
            s := scanner.Text()
            if strings.HasPrefix(s,"userAgent") {
                ss := strings.Split(s,"\"")
                userag = ss[len(ss)-2]
            } else if strings.HasPrefix(s,"port") {
                ss := strings.Split(s," ")
                port,_ = strconv.Atoi(ss[len(ss)-1])
            }
        }
        fmt.Println("loaded userag "+userag+ " port "+strconv.Itoa(port))

	// This creates the actor if it doesn't exist.
	actor, _ = activityserve.GetActor(userag, "This is polson father", "Service")

        fmt.Printf("actor created %T\n",actor)

	// actor.Follow("https://pleroma.site/users/qwazix")
	// actor.CreateNote("Hello World!", "")
	// let's boost @tzo's fox
	// actor.Announce("https://cybre.space/@tzo/102564367759300737")

	// this can be run any subsequent time
	// actor, _ := activityserve.LoadActor("activityserve_test_actor_2")

	// available actor events at this point are .OnReceiveContent and .OnFollow
	actor.OnReceiveContent = func(activity map[string]interface{}) {
		object := activity["object"].(map[string]interface{})
                gotmsg(object)
	}

        go func() {
            activityserve.ServeSingleActor(actor,port)
        }()
        fmt.Println("starting cli")
        cli()
}

