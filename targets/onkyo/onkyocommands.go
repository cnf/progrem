package onkyo

import "fmt"
import "strconv"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

func (d *OnkyoReceiver) Commands() map[string]*targets.Command {
    cmds := map[string]*targets.Command {
        "PowerOn"     : targets.NewCommand("PowerOn", "Powers on the receiver"),
        "PowerOff"    : targets.NewCommand("PowerOff", "Powers off the receiver"),
        "PowerToggle" : targets.NewCommand("PowerOff", "Powers off the receiver"),
        "Power"       : targets.NewCommand("PowerOff", "Powers off the receiver",
                targets.NewParameter("powerstate", "The power state", false).SetList("on", "off", "toggle"),
                ),
        "MuteOn"      : targets.NewCommand("MuteOn", "Mutes the sound"),
        "MuteOff"     : targets.NewCommand("MuteOff", "Unmutes the sound"),
        "MuteToggle"  : targets.NewCommand("MuteToggle", "Toggles the muting of the sound"),
        "Mute"        : targets.NewCommand("Mute", "Controls the Mute state",
                targets.NewParameter("mutestate", "The mute state", false).SetList("on", "off", "toggle"),
                ),
        "VolumeUp"    : targets.NewCommand("VolumeUp", "Turns up the volume"),
        "VolumeDown"  : targets.NewCommand("VolumeDown", "Turns down the volume"),
        "Volume"      : targets.NewCommand("Volume", "Sets the volume",
                targets.NewParameter("volumelevel", "The volume level", false).SetRange(0, 77),
                ),
    }
    if (true) {
        // Fix the volume range for specific models
        cmds["Volume"].Parameters[0].SetRange(0,77)
    }
    return cmds
}

func (r *OnkyoReceiver) Mute(state string) (string, error) {
    var rv string
    var err error
    switch state {
    case "on":
        rv, err = r.sendCmd("AMT01", 0)
    case "off":
        rv, err = r.sendCmd("AMT00", 0)
    case "toggle":
        rv, err = r.sendCmd("AMTTG", 0)
    }
    return rv, err
}

func (r *OnkyoReceiver) Power(state string) (string, error) {
    var rv string
    var err error

    switch state {
    case "on":
        rv, err = r.sendCmd("PWR01", -1)
    case "off":
        rv, err = r.sendCmd("PWR00", -1)
    case "toggle":
        rv, err = r.sendCmd("PWRQSTN", -1)
        if err != nil {
            clog.Error("ERROR: %s", err.Error())
            return "", err
        }
        clog.Debug("Power state query: '%s', %d", rv, len(rv))
        if rv == "PWR00" {
            clog.Debug("Sending PWR01")
            r.sendCmd("PWR01", -1)
        } else {
            clog.Debug("Sending PWR00")
            r.sendCmd("PWR00", -1)
        }
    }
    return rv, err
}
func (r *OnkyoReceiver) onkyoCommand(cmd string, args []string) error {
    var err error
    switch cmd {
    case "PowerOn":
        _, err = r.Power("on")
    case "PowerOff":
        _, err = r.Power("off")
    case "PowerToggle":
        _, err = r.Power("toggle")
    case "Power":
        _, err = r.Power(args[0])
    case "MuteOn":
        _, err = r.Mute("on")
    case "MuteOff":
        _, err = r.Mute("off")
    case "MuteToggle":
        _, err = r.Mute("toggle")
    case "Mute":
        _, err = r.Mute(args[0])
    case "VolumeUp":
        _, err = r.sendCmd("MVLUP",0)
    case "VolumeDown":
        _, err = r.sendCmd("MVLDOWN",0)
    case "Volume":
        ml, _ := strconv.Atoi(args[0])
        // TODO: Most models require hex volume level, some require decimal!
        _, err = r.sendCmd(fmt.Sprintf("MVL%02X", ml), 0)
    }
    return err
}