// Package utils provides functions to extract device and location information from HTTP requests.
package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
)

// MetadataExtractor holds the pre-compiled regex patterns
type MetadataExtractor struct {
	parser *uaparser.Parser
}

// NewMetadataExtractor initializes the parser patterns
func NewMetadataExtractor() *MetadataExtractor {
	// uaparser.New() uses the embedded YAML patterns by default
	return &MetadataExtractor{
		parser: uaparser.NewFromSaved(),
	}
}

type RequestMetadata struct {
	IP        string
	OS        string
	Device    string
	UserAgent string
	Location  map[string]string
}

func (e *MetadataExtractor) GetMetadata(c *gin.Context) RequestMetadata {
	uaString := c.GetHeader("User-Agent")
	client := e.parser.Parse(uaString)

	// Format OS: e.g., "iOS 15.4" or "Windows 11"
	os := fmt.Sprintf("%s %s", client.Os.Family, client.Os.Major)
	if client.Os.Minor != "" {
		os = fmt.Sprintf("%s.%s", os, client.Os.Minor)
	}

	// Format Device: e.g., "iPhone", "Samsung SM-G900P", or "Generic Desktop"
	device := client.Device.Family
	if device == "Other" {
		device = "Desktop/PC"
	}

	ip := getClientIP(c)
	loc, err := GetLocationFromIP(ip)
	if err != nil {
		loc = nil
	}

	return RequestMetadata{
		IP:        ip,
		OS:        strings.TrimSpace(os),
		Device:    device,
		UserAgent: uaString,
		Location:  loc,
	}
}

func getClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	if xrip := c.GetHeader("X-Real-IP"); xrip != "" {
		return xrip
	}
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return ip
}

func GetLocationFromIP(ip string) (map[string]string, error) {
	// Note: In production, use a local .mmdb file for speed.
	// This is an external network call.
	var rtn = make(map[string]string)
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		rtn["status"] = "fail"
		return nil, err
	}
	defer resp.Body.Close()

	// 3. Use json.Unmarshal to decode the data into the map
	err = json.NewDecoder(resp.Body).Decode(&rtn)
	if err != nil {
		return nil, err
	}

	return rtn, nil
}
