package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// DEFAULTDOMAINS Default domain list
var DEFAULTDOMAINS = []string{
	"twitter.com",
	"facebook.com",
	"youtube.com",
	"reddit.com"}

// LOCKFILE Location of the lockfile
var LOCKFILE = "/var/run/procrastistop.lock"

func main() {
	VERSION := "1.0"
	USAGE := fmt.Sprintf("USAGE: %s block|allow", os.Args[0])

	log.Printf("Running as %s", os.Args)

	if len(os.Args) != 2 {
		log.Println("Error: Unable to parse arguments")
		log.Fatal(USAGE)
	}

	mode := os.Args[1]

	switch mode {
	case "block":
		block()
	case "allow":
		allow()
	case "version":
		log.Printf("procrastistop v%s", VERSION)
		os.Exit(0)
	default:
		log.Println("Error: Unknown argument")
		log.Fatal(USAGE)
	}

	log.Println("Finished successfully")
	log.Println("--")
}

func block() {
	// Block domains

	if _, err := os.Stat(LOCKFILE); !os.IsNotExist(err) {
		log.Println("Lockfile present. Not doing anything")
		os.Exit(0)
	}

	err := cpFile("/etc/hosts", "/etc/hosts-bak")
	if err != nil {
		log.Println("Unable to backup file")
		log.Fatal(err)
	}

	lock, err := os.Create(LOCKFILE)
	if err != nil {
		log.Printf("Error: Unable to create lock file. Aborting")
		log.Fatal(err)
	}
	defer lock.Close()

	err = addDomains()
	if err != nil {
		log.Println("Unable to add domains to /etc/hosts")
		_ = cpFile("/etc/hosts-bak", "/etc/hosts")
		log.Fatal(err)
	}
}

func allow() {
	// Allow domains
	backup := "/etc/hosts-bak"

	if _, err := os.Stat(backup); os.IsNotExist(err) {
		log.Fatal("Backup file does not exist")
	}

	err := cpFile(backup, "/etc/hosts")
	if err != nil {
		log.Println("Unable to restore file")
		log.Fatal(err)
	}

	err = os.Remove(backup)
	if err != nil {
		log.Println("Unable to delete backup file after restoring")
		log.Fatal(err)
	}

	err = os.Remove(LOCKFILE)
	if err != nil {
		log.Printf("Error: Unable to remove lockfile %s", LOCKFILE)
		log.Fatal(err)
	}
}

func cpFile(input string, output string) error {
	// Copy to and from backup file
	file, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(output, file, 0644)
	return err

}

func addDomains() error {
	// Add desired domains to hosts file
	file, err := os.OpenFile("/etc/hosts", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	domains, err := readDomains()
	if err != nil {
		domains = DEFAULTDOMAINS
	}

	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}

	for i := 0; i < len(domains); i++ {
		dom := fmt.Sprintf("127.0.0.1	%s\n127.0.0.1	www.%s\n\n", domains[i], domains[i])
		_, err = file.WriteString(dom)
		if err != nil {
			return err
		}
	}

	return err

}

func readDomains() ([]string, error) {
	domainsFile := "/etc/procrastistop/domains.conf"
	var domains []string

	file, err := os.Open(domainsFile)
	if err != nil {
		log.Println("Unable to open domain config file")
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		log.Println("Unable to parse domain list")
		return nil, err
	}

	return domains, nil

}
