package lircsocket

import "net"
import "io"
import "strings"
import "bufio"
import "strconv"
import "time"

import "github.com/cnf/go-claw/commandstream"
import "github.com/cnf/go-claw/clog"

type LircSocketListener struct {
    Path string
    // conn net.Conn
    reader *bufio.Reader
}

func Register() {
    commandstream.RegisterListener("lircsocket", Create)
}

func Create(ptype string, params map[string]string) (l commandstream.Listener, ok bool) {
    // var sl listeners.Listener
    sl := &LircSocketListener{}
    if val, ok := params["path"]; ok {
        clog.Warn("Incorrect parameters")
        // sl := &LircSocketListener{Path: val}
        sl.Path = val
    } else {
        return nil, false
    }
    return sl, false
}

func (self *LircSocketListener) Setup(cs *commandstream.CommandStream) bool {
    clog.Debug("Opening socket: %s", self.Path)
    c, err := net.Dial("unix", self.Path)
    // If there is no socket to bind to during setup, we fail.
    if err != nil {
        clog.Warn("Socket setup failed for %s", self.Path)
        cs.ChErr <- err
        return false
    }
    self.reader = bufio.NewReader(c)
    return true
}

func (self *LircSocketListener) RunListener(cs *commandstream.CommandStream) {
    // var err error
    // self.conn, err = net.Dial("unix", self.Path)
    // if err != nil {
        // cs.ChErr <- err
        // return
    // }
    if (!self.Setup(cs)) {
        cs.Fatal = true
        return
    }
    for {
        now := time.Now()
        str, err := self.reader.ReadString('\n')
        if err != nil {
            if err != io.EOF {
                // Remote end closed socket
                clog.Error("Unknown error occured: %s", err.Error())
            } else {
                clog.Error("Socket closed by remote host: %s", err.Error())
                time.Sleep(1000 * time.Millisecond)
                // var err error
                // self.conn, err = net.Dial("unix", self.Path)
                if (!self.Setup(cs)) {
                    time.Sleep(3000 * time.Millisecond)
                    continue
                }
                continue
            }
            // cs.ChErr <- err
            // return
            time.Sleep(1000 * time.Millisecond)
            continue
        }

        out := strings.Split(strings.TrimSpace(str), " ")
        if (len(out) != 4) {
            clog.Error("Length of split '%v' is not 4!", str)
            continue
        }
        rpt, err := strconv.ParseInt(out[1], 16, 0)
        if (err != nil) {
            clog.Error("Could not parse %v, not a number?", out[1])
            continue
        }
        cs.Ch <- &commandstream.RemoteCommand{ Code: out[0], Repeat: int(rpt), Key: out[2], Source: out[3], Time: now }

    }
}
