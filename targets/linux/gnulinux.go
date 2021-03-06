package linux

import "fmt"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/tools"

type Linux struct {
    name string
    wol string
}

func Register() {
    targets.RegisterTarget("linux", Create)
}

func Create(name string, params map[string]string) (targets.Target, error) {
    l := &Linux{name: name}
    if wol, ok := params["wol"]; ok {
        l.wol = wol
    }
    return l, nil
}

func (d *Linux) Commands() map[string]*targets.Command {
    return nil
}
func (d *Linux) Stop() error {
    return nil
}

func (l *Linux) SendCommand(cmd string, args ...string) error {
    switch cmd {
    case "poweron":
        clog.Debug("Power on %s", l.name)
        return l.powerOn()
    }
    return fmt.Errorf("could not send command `%s` on `%s`", cmd, l.name)
}

func (l *Linux) powerOn() error {
    if l.wol != "" {
        ok := tools.Wol(l.wol)
        if !ok { return fmt.Errorf("can not power on %s", l.name) }
    }
    return fmt.Errorf("do not know how to power on %s", l.name)
}
