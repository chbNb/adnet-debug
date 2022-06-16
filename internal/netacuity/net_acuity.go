/*
* ***************************************************************************
 * File:           NetAcuity.go
 * Author:         Digital Envoy
 * Program Name:   NetAcuity API library
 * Version:        6.0.0.7
 * Date:           21-Feb-2017
 *
 * Copyright 2000-2017, Digital Envoy, Inc.  All rights reserved.
 *
 *  Description:
 *    Go implementation of the Digital Envoy NetAcuity API library
 *    to query for ip based location data.
 *
 *
 *
 * This library is provided as an access method to the NetAcuity software
 * provided to you under License Agreement from Digital Envoy Inc.
 * You may NOT redistribute it and/or modify it in any way without express
 * written consent of Digital Envoy, Inc.
 *
 * Address bug reports and comments to:  tech-support@digitalenvoy.net
 *
 *
 * **************************************************************************
*/

package netacuity

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

const IMPLEMENTATION_PROTOCOL_VERSION = "5"
const GO_API_ID = "11"
const NETACUITY_PORT = "5400"
const UDP6 = "udp6"
const UDP4 = "udp4"
const MAX_RESPONSE_SIZE = 1500
const TRANSACTION_LENGTH = 10

/** Simulated enum for feature codes **/
var FeatureCodeEnums map[int]string = map[int]string{
	3:  "geo",
	4:  "edge",
	5:  "sic",
	6:  "domain",
	7:  "zip",
	8:  "isp",
	9:  "home_biz",
	10: "asn",
	11: "language",
	12: "proxy",
	14: "is_an_isp",
	15: "company",
	17: "demographics",
	18: "naics",
	19: "cbsa",
	24: "mobile_carrier",
	25: "organization",
	26: "pulse",
	30: "pulseplus",
}

/**
Purpose : Queries a netacuity server using IPv4 or IPv6
Params :
	params [5]interface{} - An array of 5 parameters (Feature Code - int
													  API ID - int
													  IP Address - string
													  Netacuity Server IP - string
													  Timeout Delay(ms) - int)
Returns :
	interface{} - Using the empty interface in order to return different structs for different feature codes
	error - An error if one occurs
**/
func Query(params [5]interface{}) (interface{}, error) {
	// Check parameters
	featureCode := params[0].(int)
	apiId := params[1].(int)
	ipAddress := params[2].(string)
	netAcuityIp := params[3].(string)
	timeoutDelay := params[4].(int)

	var srvAddrType int

	if checkFeatureCode(featureCode) {
		if checkApiId(apiId) {
			if checkIpAddress(ipAddress) != -1 {
				srvAddrType = checkIpAddress(netAcuityIp)
				if srvAddrType != -1 {
					if checkTimeoutDelay(timeoutDelay) {

					} else {
						return nil, errors.New("Error : Invalid delay!")
					}
				} else {
					return nil, errors.New("Error : Invalid server IP address!")
				}
			} else {
				return nil, errors.New("Error : Invalid target IP address!")
			}
		} else {
			return nil, errors.New("Error : Invalid api ID!")
		}
	} else {
		return nil, errors.New("Error : Invalid feature code!")
	}

	// Generate a transaction id, set up buffer, declare vars for later use
	transactionId := generateTransactionId(TRANSACTION_LENGTH)
	var buffer [MAX_RESPONSE_SIZE]byte
	var serverIpString string
	var protocol string

	// Set up the server ip string and udp type to be used to create the udp address object
	if srvAddrType == 6 {
		serverIpString = fmt.Sprintf("[%s]:%s", netAcuityIp, NETACUITY_PORT)
		protocol = UDP6
	} else {
		serverIpString = fmt.Sprintf("%s:%s", netAcuityIp, NETACUITY_PORT)
		protocol = UDP4
	}

	// Create UDP address object and connection
	serverAddress, _ := net.ResolveUDPAddr(protocol, serverIpString)
	connection, _ := net.DialUDP(protocol, nil, serverAddress)

	// Set the timeout in milliseconds
	connection.SetDeadline(time.Now().Add(time.Duration(timeoutDelay) * time.Millisecond))

	// Close the connection when this method resolves
	defer connection.Close()

	// Build the message, send the message, and copy the response into a buffer
	udpMessage := fmt.Sprintf("%v;%v;%s;%s;%s;%s", featureCode, apiId, ipAddress, IMPLEMENTATION_PROTOCOL_VERSION,
		GO_API_ID, transactionId)
	connection.Write([]byte(udpMessage))
	connection.Read(buffer[:])

	return parseResponse(string(buffer[:]), featureCode), nil
}

/**
Purpose : Queries a netacuity server using IPv4 or IPv6 for a multifeature XML response
Params :
	params [5]interface{} - An array of 5 parameters (Feature Code(s) - comma delimited string
													  API ID - int
													  IP Address - string
													  Netacuity Server IP - string
													  Timeout Delay(ms) - int)
Returns :
	string - A string xml response
	error - An error if one occurs
**/
func QueryXml(params [5]interface{}) (string, error) {
	// Check parameters
	featureCodes := params[0].(string)
	apiId := params[1].(int)
	ipAddress := params[2].(string)
	netAcuityIp := params[3].(string)
	timeoutDelay := params[4].(int)

	isValidFc, codeSlice := checkXmlFeatureCodes(featureCodes)
	var srvAddrType int
	result := ""

	if isValidFc {
		if checkApiId(apiId) {
			if checkIpAddress(ipAddress) != -1 {
				srvAddrType = checkIpAddress(netAcuityIp)
				if srvAddrType != -1 {
					if checkTimeoutDelay(timeoutDelay) {

					} else {
						return result, errors.New("Error : Invalid delay!")
					}
				} else {
					return result, errors.New("Error : Invalid server IP address!")
				}
			} else {
				return result, errors.New("Error : Invalid target IP address!")
			}
		} else {
			return result, errors.New("Error : Invalid api ID!")
		}
	} else {
		return result, errors.New("Error : Invalid feature codes!")
	}

	// Generate a transaction id, set up buffer, declare vars for later use
	transactionId := generateTransactionId(TRANSACTION_LENGTH)
	tempBuffer := make([]byte, MAX_RESPONSE_SIZE)
	var serverIpString, protocol string

	// Set up the server ip string and udp type to be used to create the udp address object
	if srvAddrType == 6 {
		serverIpString = fmt.Sprintf("[%s]:%s", netAcuityIp, NETACUITY_PORT)
		protocol = UDP6
	} else {
		serverIpString = fmt.Sprintf("%s:%s", netAcuityIp, NETACUITY_PORT)
		protocol = UDP4
	}

	// Create UDP address object and connection
	serverAddress, _ := net.ResolveUDPAddr(protocol, serverIpString)
	connection, _ := net.DialUDP(protocol, nil, serverAddress)

	// Set the timeout in milliseconds
	connection.SetDeadline(time.Now().Add(time.Duration(timeoutDelay) * time.Millisecond))

	// Close the connection when this method resolves
	defer connection.Close()

	// Build the xml request message and then send it
	udpMessage := fmt.Sprintf("<request trans-id=\"%v\" ip=\"%v\" api-id=\"%v\">", transactionId, ipAddress, apiId)
	for _, code := range codeSlice {
		udpMessage += fmt.Sprintf(" <query db=\"%v\" />", code)
	}
	udpMessage += " </request>"
	connection.Write([]byte(udpMessage))

	// Catch all parts of a multi-packet response
	size, _ := connection.Read(tempBuffer[:])
	currentPacket, _ := strconv.Atoi(string(tempBuffer[0:2]))
	lastPacket, _ := strconv.Atoi(string(tempBuffer[2:4]))
	result += string(tempBuffer[4 : size-1])

	for currentPacket < lastPacket {
		tempBuffer = nil
		tempBuffer = make([]byte, MAX_RESPONSE_SIZE)
		size, _ := connection.Read(tempBuffer[:])
		result += string(tempBuffer[4 : size-1])
		currentPacket, _ = strconv.Atoi(string(tempBuffer[0:2]))
		lastPacket, _ = strconv.Atoi(string(tempBuffer[2:4]))
	}

	return result, nil

}
