package utils

import (
	"net"
	"regexp"
	"strings"
)

type EmailDomain struct {
	Valid bool

	Domain string
	Addr   []string

	HasMX    bool
	HasSPF   bool
	HasDMARC bool

	SPFRecord   string
	DMARCRecord string
}

var (
	emailRegexp = regexp.MustCompile("(?i)" + // case insensitive
		"^[a-z0-9!#$%&'*+/=?^_`{|}~.-]+" + // local part
		"@" +
		"[a-z0-9-]+(\\.[a-z0-9-]+)+\\.?$") // domain part

	domainRegexp = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)
)

func NormalizeEmail(email string) (local, domain string, ok bool) {
	email = strings.TrimSpace(email)
	ok = validEmail(email)
	if !ok {
		return
	}

	local, domain, ok = splitEmail(email)
	if !ok {
		return
	}

	domain = strings.TrimRight(domain, ".")
	domain = strings.ToLower(domain)

	return local, domain, ok
}

func CheckEmailDomain(domain string) *EmailDomain {
	var emailDomain EmailDomain
	emailDomain.Domain = domain

	if !validDomain(domain) {
		return &emailDomain
	}

	addr, err := net.LookupHost(domain)
	if err != nil || len(addr) == 0 {
		ErrorLogger.Printf("Failed to lookup Host for the given domain - '%s'\n", domain)
		return &emailDomain
	} else {
		emailDomain.Addr = addr
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		ErrorLogger.Printf("Failed to lookup MX for the given domain - '%s'\n", domain)
	} else if len(mxRecords) > 0 {
		emailDomain.HasMX = true
		emailDomain.Valid = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		ErrorLogger.Printf("Failed to lookup TXT for the given domain - '%s'\n", domain)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				emailDomain.HasSPF = true
				emailDomain.SPFRecord = record
				break
			}
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		ErrorLogger.Printf("Failed to lookup DMARC TXT for the given domain - '%s'\n", domain)
	} else {
		for _, record := range dmarcRecords {
			if strings.HasPrefix(record, "v=DMARC1") {
				emailDomain.HasDMARC = true
				emailDomain.DMARCRecord = record
				break
			}
		}
	}

	return &emailDomain
}

func validEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	return emailRegexp.MatchString(email)
}

func validDomain(domain string) bool {
	return domainRegexp.MatchString(domain)
}

func splitEmail(email string) (local, domain string, ok bool) {
	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return
	}

	local = parts[0]
	domain = parts[1]

	if len(local) < 1 || len(domain) < len("x.xx") {
		return
	}

	return local, domain, true
}
