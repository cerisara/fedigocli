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

/*

Notes:
- les messages que l'on recoit lorsqu'on follow qqun ne sont pas stockes en local: ils sont envoyes instantannement aux followers, qui les traitent mais ne les sauvent pas

*/

var actor activityserve.Actor

// const actype = "Service"
const actype = "Person"

func myMsg() {
    fmt.Printf("mymsg %v\n",actor)
}

func cli() {
        reader := bufio.NewReader(os.Stdin)
        for ;; {
            fmt.Print("-> ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            if strings.HasPrefix(text,"post ") {
                post(text[5:])
            } else if strings.Compare(text,"l") == 0 { myMsg()
            } else if strings.Compare(text,"w") == 0 {
                fmt.Printf("whoami %v\n",actor.WhoAmI())
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
        var userdesc string = "I'm a bot"
        var port int = 8081
        for scanner.Scan() {
            s := scanner.Text()
            if strings.HasPrefix(s,"userAgent") {
                ss := strings.Split(s,"\"")
                userag = ss[len(ss)-2]
            } else if strings.HasPrefix(s,"userDesc") {
                ss := strings.Split(s,"\"")
                userdesc = ss[len(ss)-2]
            } else if strings.HasPrefix(s,"port") {
                ss := strings.Split(s," ")
                port,_ = strconv.Atoi(ss[len(ss)-1])
            }
        }
        fmt.Println("loaded userag "+userag+ " desc "+userdesc + " port "+strconv.Itoa(port))

	// This creates the actor if it doesn't exist.
	actor, _ = activityserve.GetActor(userag, userdesc, actype)

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

