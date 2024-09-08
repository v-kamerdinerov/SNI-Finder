package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/airnandez/tlsping"
	"github.com/sirupsen/logrus"
)

const (
	defaultAddress     = "0.0.0.0"
	defaultPort        = "443"
	defaultThreadCount = 128
	defaultTimeout     = 4
	outPutDef          = true
	outPutFileName     = "results.txt"
	domainsFileName    = "domains.txt"
	showFailDef        = false
	numIPsToCheck      = 10000
	tlsCount           = 3
	tlsHandshake       = false
	tlsVerify          = true
	numberBestServers  = 10
)

var log = logrus.New()
var zeroIP = net.ParseIP("0.0.0.0")
var maxIP = net.ParseIP("255.255.255.255")
var TlsDic = map[uint16]string{
	0x0301: "1.0",
	0x0302: "1.1",
	0x0303: "1.2",
	0x0304: "1.3",
}

type Scanner struct {
	addr           string
	port           string
	showFail       bool
	output         bool
	timeout        time.Duration
	wg             sync.WaitGroup
	numberOfThread int
	mu             sync.Mutex
	ip             net.IP
	logFile        *os.File
	domainFile     *os.File
	dialer         *net.Dialer
	logChan        chan string
}

func main() {
	addrPtr := flag.String("addr", defaultAddress, "Destination to start scan")
	portPtr := flag.String("port", defaultPort, "Port to scan")
	threadPtr := flag.Int("thread", defaultThreadCount, "Number of threads to scan in parallel")
	outPutFile := flag.Bool("o", outPutDef, "Is output to results.txt")
	timeOutPtr := flag.Int("timeOut", defaultTimeout, "Time out of a scan")
	showFailPtr := flag.Bool("showFail", showFailDef, "Is Show fail logs")

	flag.Parse()
	scanner := newScanner(*addrPtr, *portPtr, *threadPtr, *timeOutPtr, *outPutFile, *showFailPtr)

	defer scanner.logFile.Close()
	defer scanner.domainFile.Close()

	go scanner.logWriter()

	// Start the worker pool
	scanner.startWorkers()

	log.Info("Scan completed.")

	// Choice best servers
	findTopServers(outPutFileName)
}

func newScanner(addr, port string, threadCount, timeout int, output, showFail bool) *Scanner {
	scanner := &Scanner{
		addr:           addr,
		port:           port,
		showFail:       showFail,
		output:         output,
		timeout:        time.Duration(timeout) * time.Second,
		numberOfThread: threadCount,
		ip:             net.ParseIP(addr),
		dialer:         &net.Dialer{},
		logChan:        make(chan string, numIPsToCheck),
	}

	log.SetFormatter(&CustomTextFormatter{})
	log.SetLevel(logrus.InfoLevel)

	var err error
	scanner.logFile, err = os.OpenFile(outPutFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.WithError(err).Fatal("Failed to open log file")
	}

	scanner.domainFile, err = os.OpenFile(domainsFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.WithError(err).Fatal("Failed to open domains.txt file")
	}

	return scanner
}

func (s *Scanner) startWorkers() {
	ipChan := make(chan net.IP, numIPsToCheck)

	for i := 0; i < s.numberOfThread; i++ {
		go s.worker(ipChan)
	}

	for i := 0; i < numIPsToCheck; i++ {
		nextIP := s.nextIP(true)
		if nextIP != nil {
			s.wg.Add(1)
			ipChan <- nextIP
		}
	}

	close(ipChan)
	s.wg.Wait()
	close(s.logChan)
}

func (s *Scanner) logWriter() {
	for str := range s.logChan {
		log.Info(str)
		if s.output {
			_, err := s.logFile.WriteString(str + "\n")
			if err != nil {
				log.WithError(err).Error("Error writing into file")
			}
		}
	}
}

func (s *Scanner) worker(ipChan <-chan net.IP) {
	for ip := range ipChan {
		s.Scan(ip)
		s.wg.Done()
	}
}

func (s *Scanner) nextIP(increment bool) net.IP {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipb := big.NewInt(0).SetBytes(s.ip.To4())
	if increment {
		ipb.Add(ipb, big.NewInt(1))
	} else {
		ipb.Sub(ipb, big.NewInt(1))
	}

	b := ipb.Bytes()
	b = append(make([]byte, 4-len(b)), b...)
	nextIP := net.IP(b)

	if nextIP.Equal(zeroIP) || nextIP.Equal(maxIP) {
		return nil
	}

	s.ip = nextIP
	return s.ip
}

func (s *Scanner) Scan(ip net.IP) {
	str := ip.String()
	ping := time.Duration(0)
	if ip.To4() == nil {
		str = "[" + str + "]"
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	conn, err := s.dialer.DialContext(ctx, "tcp", net.JoinHostPort(str, s.port))
	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprintf("Dial failed: %v", err), ping)
		}
		return
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	remoteIP := remoteAddr.IP.String()
	port := remoteAddr.Port
	line := fmt.Sprintf("%s:%d", remoteIP, port)

	if err := conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
		log.WithError(err).Error("Error setting deadline")
		return
	}

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	})
	err = tlsConn.Handshake()

	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprintf("%s - TLS handshake failed: %v", line, err), ping)
		}
		return
	}

	defer tlsConn.Close()

	state := tlsConn.ConnectionState()
	alpn := state.NegotiatedProtocol

	if alpn == "" {
		alpn = "  "
	}

	if s.showFail || (state.Version == 0x0304 && alpn == "h2") {
		certSubject := ""
		if len(state.PeerCertificates) > 0 {
			certSubject = state.PeerCertificates[0].Subject.CommonName
		}

		numPeriods := strings.Count(certSubject, ".")

		if strings.HasPrefix(certSubject, "*") || certSubject == "localhost" || numPeriods != 1 || certSubject == "invalid2.invalid" || certSubject == "OPNsense.localdomain" {
			return
		}

		// Config for tlsping
		config := tlsping.Config{
			Count:              tlsCount,
			AvoidTLSHandshake:  tlsHandshake,
			InsecureSkipVerify: tlsVerify,
		}

		result, err := tlsping.Ping(certSubject+":443", &config)
		avgDuration := time.Duration(result.Avg * float64(time.Second))
		ping := avgDuration.Truncate(time.Microsecond)
		if err != nil {
			s.Print(fmt.Sprintf("%s ---- TLS v%s    ALPN: %s ----    %s:%s ---- TCP ping failed: %v", line, TlsDic[state.Version], alpn, certSubject, s.port, err), ping)
			return
		}

		s.Print(fmt.Sprintf("%s ---- TLS v%s    ALPN: %s ----    %s:%s", line, TlsDic[state.Version], alpn, certSubject, s.port), ping)
	}
}

func (s *Scanner) Print(outStr string, ping time.Duration) {
	// Split the output string into IP address and the rest
	parts := strings.Split(outStr, " ")
	ipAddress := parts[0]
	rest := strings.Join(parts[1:], " ")

	// Format the IP address with a fixed width
	formattedIP := fmt.Sprintf("%-22s", ipAddress)

	// Extract and format TLS and ALPN
	restParts := strings.Split(rest, "----")
	var tlsAndAlpn string
	if len(restParts) > 1 {
		tlsAndAlpn = strings.TrimSpace(restParts[1])
	} else {
		tlsAndAlpn = "Unknown TLS/ALPN"
	}
	formattedTLS := fmt.Sprintf("%-22s", tlsAndAlpn)

	// Extract domain from the log entry
	domain := extractDomain(outStr)

	// Correctly format domain
	var formattedDomain string
	if domain != "" && domain != ipAddress {
		formattedDomain = fmt.Sprintf("%-22s", domain)
	} else {
		formattedDomain = fmt.Sprintf("%-22s", "") // If no domain, leave it empty
	}

	// Format ping duration
	var formattedPing string
	if ping == 0 {
		formattedPing = ""
	} else {
		formattedPing = fmt.Sprintf("Ping: %-30s", ping)
	}

	// Create the final log entry with alignment
	logEntry := fmt.Sprintf("%s%s", formattedIP, formattedTLS)

	if formattedDomain != "" {
		logEntry += formattedDomain
	}

	if formattedPing != "" {
		logEntry += formattedPing
	}

	// Save the domain to domains.txt if needed
	if domain != "" && domain != ipAddress {
		saveDomain(domain, s.domainFile)
	}

	// Send the log entry to the log channel
	s.logChan <- logEntry
}

func extractDomain(logEntry string) string {
	parts := strings.Fields(logEntry)
	for i, part := range parts {
		if strings.Contains(part, ".") && !strings.HasPrefix(part, "v") && i > 0 {
			domainParts := strings.Split(part, ":")
			return domainParts[0]
		}
	}
	return ""
}

func saveDomain(domain string, file *os.File) {
	if domain != "" {
		_, err := file.WriteString(domain + "\n")
		if err != nil {
			log.WithError(err).Error("Error writing domain into file")
		}
	}
}

type CustomTextFormatter struct {
	logrus.TextFormatter
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	msg := entry.Message
	formattedEntry := timestamp + " " + msg + "\n\n"
	return []byte(formattedEntry), nil
}

func findTopServers(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.WithError(err).Fatal("Failed to open results.txt file for reading")
	}
	defer file.Close()

	type Server struct {
		Line string
		Ping time.Duration
	}

	var servers []Server

	// Regex to extract Ping value
	pingRegex := regexp.MustCompile(`Ping:\s*([0-9]+(?:\.[0-9]+)?[a-z]+)`)
	// Regex to filter lines with ALPN: h2 and a domain name
	alpnH2Regex := regexp.MustCompile(`ALPN:\s*h2\s+([a-zA-Z0-9\.\-]+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line matches the ALPN: h2 filter
		if alpnH2Regex.MatchString(line) {
			matches := pingRegex.FindStringSubmatch(line)

			// Extract the Ping value if present
			if len(matches) > 1 {
				pingStr := matches[1]
				ping, err := time.ParseDuration(pingStr)
				if err == nil {
					// Add the server line and parsed ping duration to the slice
					servers = append(servers, Server{Line: line, Ping: ping})
				} else {
					log.WithError(err).Errorf("Failed to parse ping duration from: %s", pingStr)
				}
			}
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		log.WithError(err).Fatal("Error reading from results.txt file")
	}

	// Sort servers by Ping value
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Ping < servers[j].Ping
	})

	// Determine the number of servers to display
	topCount := numberBestServers
	if len(servers) < topCount {
		topCount = len(servers)
	}

	// Display top servers, keeping original lines from the results file
	fmt.Println("Top servers by TLS Ping:")
	for i := 0; i < topCount; i++ {
		// Print the original line from the results file
		fmt.Printf("%d: %s\n", i+1, servers[i].Line)
	}
}
