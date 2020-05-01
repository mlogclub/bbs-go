<img src="gh_logo.png" />

[![build status](https://img.shields.io/travis/kataras/neffos/master.svg?style=for-the-badge)](https://travis-ci.org/kataras/neffos) [![report card](https://img.shields.io/badge/report%20card-a%2B-ff3333.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/kataras/neffos)<!--[![godocs](https://img.shields.io/badge/go-%20docs-488AC7.svg?style=for-the-badge)](https://godoc.org/github.com/kataras/neffos)--> [![view examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=for-the-badge)](https://github.com/kataras/neffos/tree/master/_examples) [![chat](https://img.shields.io/gitter/room/neffos-framework/community.svg?color=blue&logo=gitter&style=for-the-badge)](https://gitter.im/neffos-framework/community) [![frontend pkg](https://img.shields.io/badge/JS%20-client-BDB76B.svg?style=for-the-badge)](https://github.com/kataras/neffos.js)

## About neffos

Neffos is a cross-platform real-time framework with expressive, elegant API written in [Go](https://golang.org). Neffos takes the pain out of development by easing common tasks used in real-time backend and frontend applications such as:

- Scale-out using redis or nats[*](_examples/scale-out)
- Adaptive request upgradation and server dialing
- Acknowledgements
- Namespaces
- Rooms
- Broadcast
- Event-Driven architecture
- Request-Response architecture
- Error Awareness
- Asynchronous Broadcast
- Timeouts
- Encoding
- Reconnection
- Modern neffos API client for Browsers, Nodejs[*](https://github.com/kataras/neffos.js) and Go

## Learning neffos

<details>
<summary>Qick View</summary>

## Server

```go
import (
    // [...]
    "github.com/kataras/neffos"
    "github.com/kataras/neffos/gorilla"
)

func runServer() {
    events := make(neffos.Namespaces)
    events.On("/v1", "workday", func(ns *neffos.NSConn, msg neffos.Message) error {
        date := string(msg.Body)

        t, err := time.Parse("01-02-2006", date)
        if err != nil {
            if n := ns.Conn.Increment("tries"); n >= 3 && n%3 == 0 {
                // Return custom error text to the client.
                return fmt.Errorf("Why not try this one? 06-24-2019")
            } else if n >= 6 && n%2 == 0 {
                // Fire the "notify" client event.
                ns.Emit("notify", []byte("What are you doing?"))
            }
            // Return the parse error back to the client.
            return err
        }

        weekday := t.Weekday()

        if weekday == time.Saturday || weekday == time.Sunday {
            return neffos.Reply([]byte("day off"))
        }

        // Reply back to the client.
        responseText := fmt.Sprintf("it's %s, do your job.", weekday)
        return neffos.Reply([]byte(responseText))
    })

    websocketServer := neffos.New(gorilla.DefaultUpgrader, events)

    // Fire the "/v1:notify" event to all clients after server's 1 minute.
    time.AfterFunc(1*time.Minute, func() {
        websocketServer.Broadcast(nil, neffos.Message{
            Namespace: "/v1",
            Event:     "notify",
            Body:      []byte("server is up and running for 1 minute"),
        })
    })

    router := http.NewServeMux()
    router.Handle("/", websocketServer)

    log.Println("Serving websockets on localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Go Client

```go
func runClient() {
    ctx := context.TODO()
    events := make(neffos.Namespaces)
    events.On("/v1", "notify", func(c *neffos.NSConn, msg neffos.Message) error {
        log.Printf("Server says: %s\n", string(msg.Body))
        return nil
    })

    // Connect to the server.
    client, err := neffos.Dial(ctx,
        gorilla.DefaultDialer,
        "ws://localhost:8080",
        events)
    if err != nil {
        panic(err)
    }

    // Connect to a namespace.
    c, err := client.Connect(ctx, "/v1")
    if err != nil {
        panic(err)
    }

    fmt.Println("Please specify a date of format: mm-dd-yyyy")

    for {
        fmt.Print(">> ")
        var date string
        fmt.Scanln(&date)

        // Send to the server and wait reply to this message.
        response, err := c.Ask(ctx, "workday", []byte(date))
        if err != nil {
            if neffos.IsCloseError(err) {
                // Check if the error is a close signal,
                // or make use of the `<- client.NotifyClose`
                // read-only channel instead.
                break
            }

            // >> 13-29-2019
            // error received: parsing time "13-29-2019": month out of range
            fmt.Printf("error received: %v\n", err)
            continue
        }

        // >> 06-29-2019
        // it's a day off!
        //
        // >> 06-24-2019
        // it's Monday, do your job.
        fmt.Println(string(response.Body))
    }
}
```

## Javascript Client

Navigate to: <https://github.com/kataras/neffos.js>

</details>

Neffos contains extensive and thorough **[wiki](https://github.com/kataras/neffos/wiki)** making it easy to get started with the framework.

For a more detailed technical documentation you can head over to our [godocs](https://godoc.org/github.com/kataras/neffos). And for executable code you can always visit the [_examples](_examples/) repository's subdirectory.

### Do you like to read while traveling?

You can [request](https://bit.ly/neffos-req-book) a PDF version of the **E-Book** today and be participated in the development of neffos.

[![https://iris-go.com/images/neffos-book-overview.png](https://iris-go.com/images/neffos-book-overview.png)](https://bit.ly/neffos-req-book)

## Contributing

We'd love to see your contribution to the neffos real-time framework! For more information about contributing to the neffos project please check the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## Security Vulnerabilities

If you discover a security vulnerability within neffos, please send an e-mail to [neffos-go@outlook.com](mailto:neffos-go@outlook.com). All security vulnerabilities will be promptly addressed.

## License

The word "neffos" has a greek origin and it is translated to "cloud" in English dictionary.

The neffos real-time framework is open-source software licensed under the [MIT license](https://opensource.org/licenses/MIT).
<!-- [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkataras%2Fneffos.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkataras%2Fneffos?ref=badge_large) -->
