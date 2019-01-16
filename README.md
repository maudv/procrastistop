# procrastistop

Simple tool to prevent procrastination by blocking domains using hosts file

## Requirements
1. sudo access 
2. go

## Usage

Create domain list on **/etc/procrastistop/domains.conf**  
Run the tool as root: `sudo procrastistop block`  
To revert changes and allow procrastination again: `sudo procrastistop allow`  

## Automated installer
Simply run the installer with sudo `sudo ./installer.sh`  

## Compile
Set the gopath to the script location: `GOPATH=*source-location*`
Compile with `go build`