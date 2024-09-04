package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
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
	workerPoolSize     = 100
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

type CustomTextFormatter struct {
	logrus.TextFormatter
}

type Scanner struct {
	addr           string
	port           string
	showFail       bool
	output         bool
	timeout        time.Duration
	wg             *sync.WaitGroup
	numberOfThread int
	mu             sync.Mutex
	ip             net.IP
	logFile        *os.File
	domainFile     *os.File
	dialer         *net.Dialer
	logChan        chan string
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	msg := entry.Message

	formattedEntry := timestamp + msg + "\n\n"

	return []byte(formattedEntry), nil
}

func (s *Scanner) Print(outStr string) {
	parts := strings.Split(outStr, " ")
	ipAddress := parts[0]
	rest := strings.Join(parts[1:], " ")

	maxIPLength := len("255.255.255.255")
	formattedIP := fmt.Sprintf("%-*s", maxIPLength-8, ipAddress)

	logEntry := formattedIP + rest

	domain := extractDomain(logEntry)

	saveDomain(domain, s.domainFile)

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

func main() {
	addrPtr := flag.String("addr", defaultAddress, "Destination to start scan")
	portPtr := flag.String("port", defaultPort, "Port to scan")
	threadPtr := flag.Int("thread", defaultThreadCount, "Number of threads to scan in parallel")
	outPutFile := flag.Bool("o", outPutDef, "Is output to results.txt")
	timeOutPtr := flag.Int("timeOut", defaultTimeout, "Time out of a scan")
	showFailPtr := flag.Bool("showFail", showFailDef, "Is Show fail logs")

	flag.Parse()
	s := &Scanner{
		addr:           *addrPtr,
		port:           *portPtr,
		showFail:       *showFailPtr,
		output:         *outPutFile,
		timeout:        time.Duration(*timeOutPtr) * time.Second,
		wg:             &sync.WaitGroup{},
		numberOfThread: *threadPtr,
		ip:             net.ParseIP(*addrPtr),
		dialer:         &net.Dialer{},
		logChan:        make(chan string, numIPsToCheck),
	}

	log.SetFormatter(&CustomTextFormatter{})
	log.SetLevel(logrus.InfoLevel)

	var err error
	s.logFile, err = os.OpenFile(outPutFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.WithError(err).Error("Failed to open log file")
		return
	}
	defer s.logFile.Close()

	s.domainFile, err = os.OpenFile(domainsFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.WithError(err).Error("Failed to open domains.txt file")
		return
	}
	defer s.domainFile.Close()

	go s.logWriter()

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
	log.Info("Scan completed.")
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

	if ip.To4() == nil {
		str = "[" + str + "]"
	}

	pinger, err := ping.NewPinger(str)
	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprintf("%s - Ping failed: %v", str, err))
		}
		return
	}
	pinger.Count = 1
	pinger.Timeout = s.timeout

	err = pinger.Run()
	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprintf("%s - Ping run failed: %v", str, err))
		}
		return
	}
	stats := pinger.Statistics()
	rtt := stats.AvgRtt

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	conn, err := s.dialer.DialContext(ctx, "tcp", str+":"+s.port)
	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprintf("%s - Dial failed: %v", str, err))
		}
		return
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	remoteIP := remoteAddr.IP.String()
	port := remoteAddr.Port
	line := fmt.Sprintf("%s:%d", remoteIP, port) + "\t"
	conn.SetDeadline(time.Now().Add(s.timeout))
	c := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	})
	err = c.Handshake()

	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprintf("%s - TLS handshake failed: %v", line, err))
		}
		return
	}
	defer c.Close()

	state := c.ConnectionState()
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

		s.Print(fmt.Sprintf(" %s ---- TLS v%s    ALPN: %s ----    %s:%s ---- Ping RTT: %v", line, TlsDic[state.Version], alpn, certSubject, s.port, rtt))
	}
}
