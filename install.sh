#!/bin/bash

compile_code() {
    # Compiles code using go
    export GOPATH=$SCRIPT_DIR
    /usr/local/go/bin/go build
}

setup_cron() {
    # Enables cronjob to execute script every n minutes
    croncmd="/usr/local/bin/procrastistop block > /var/log/procrastistop.log 2>&1"
    cronjob="*/15 * * * *  $croncmd"

    crontab -l | grep -v -F "$croncmd" ; echo "$cronjob" | crontab -
}

move_binary() {
    # Moves binary to bin location
    mv procrastistop /usr/local/bin
    chmod +x /usr/local/bin
}

systemd_handle () {
    # Adds systemd config to handle suspend

    cat <<'EOF' > /lib/systemd/system-sleep/procrastistop
#!/bin/sh
set -e

case $1 in
pre)
    /usr/local/bin/procrastistop allow > /var/log/procrastistop.log 2>&1
    ;;
esac
EOF

    cat <<'EOF' > /lib/systemd/system-shutdown/procrastistop
#!/bin/sh
set -e

case $1 in
halt)
    /usr/local/bin/procrastistop allow > /var/log/procrastistop.log 2>&1
    ;;
poweroff)
    /usr/local/bin/procrastistop allow > /var/log/procrastistop.log 2>&1
    ;;
reboot)
    /usr/local/bin/procrastistop allow > /var/log/procrastistop.log 2>&1
    ;;
esac
EOF
}


main() {
    SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

    echo "Installing procrastistop..."
    cd "$SCRIPT_DIR" || exit

    compile_code
    move_binary
    setup_cron
    systemd_handle
    echo "Done"
}

main