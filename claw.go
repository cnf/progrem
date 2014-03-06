package main

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/commandstream"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"
import "os"
import "os/signal"

// import "os"

func main() {
    defer clog.Stop()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        clog.Stop()
        os.Exit(1)
    }()

    cfg.Setup()
    cfg.ReadConfigfile()

    // if Verbose {
    //     clog.SetLogLevel(clog.DEBUG)
    // } else {
    //     clog.SetLogLevel(clog.WARN)
    // }

    cs := commandstream.NewCommandStream()
    defer cs.Close()
    var out commandstream.RemoteCommand

    cs.AddListener(&listeners.LircSocketListener{Path: "/var/run/lirc/lircd"})
    cs.AddListener(&listeners.LircSocketListener{Path: "/tmp/echo.sock"})

    dispatcher.Setup(cfg.System.Targets)

    for cs.Next(&out) {
        if cs.HasError() {
            clog.Warn("An error occured somewhere: %v", cs.GetError())
            cs.ClearError()
        }
        dispatcher.Dispatch(&out)
        clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
    }
}
