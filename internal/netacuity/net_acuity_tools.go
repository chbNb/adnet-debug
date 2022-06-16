/*
* ***************************************************************************
 * File:           NetAcuityTools.go
 * Author:         Digital Envoy
 * Program Name:   NetAcuity API library
 * Version:        6.0.0.7
 * Date:           21-Feb-2017
 *
 * Copyright 2000-2017, Digital Envoy, Inc.  All rights reserved.
 *
 *  Description:
 *    Supporting functions for the Go implementation
 *    of the Digital Envoy NetAcuity API library to query
 *    for ip based location data.
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
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

/**
Structures for each feature code. Note that all floats are float64 because go's
string to float conversion defaults to a 64 bit float, even if 32 bit is specified
**/
type na_geo struct {
	ip                     string
	trans_id               string
	geo_country            string
	geo_region             string
	geo_city               string
	geo_conn_speed         string
	geo_country_conf       int
	geo_region_conf        int
	geo_city_conf          int
	geo_metro_code         int
	geo_latitude           float64
	geo_longitude          float64
	geo_country_code       int
	geo_region_code        int
	geo_city_code          int
	geo_continent_code     int
	geo_two_letter_country string
}

type na_edge struct {
	ip                      string
	trans_id                string
	edge_country            string
	edge_region             string
	edge_city               string
	edge_conn_speed         string
	edge_metro_code         int
	edge_latitude           float64
	edge_longitude          float64
	edge_postal_code        string
	edge_country_code       int
	edge_region_code        int
	edge_city_code          int
	edge_continent_code     int
	edge_two_letter_country string
	edge_internal_code      int
	edge_area_codes         string
	edge_country_conf       int
	edge_region_conf        int
	edge_city_conf          int
	edge_postal_conf        int
	edge_gmt_offset         int
	edge_in_dst             string
}

type na_sic struct {
	ip       string
	trans_id string
	sic_code int
}

type na_domain struct {
	ip          string
	trans_id    string
	domain_name string
}

type na_zip struct {
	ip            string
	trans_id      string
	area_code     int
	zip_code      int
	gmt_offset    int
	in_dst        string
	zip_code_text string
	zip_country   string
}

type na_isp struct {
	ip       string
	trans_id string
	isp_name string
}

type na_home_biz struct {
	ip           string
	trans_id     string
	homebiz_type string
}

type na_asn struct {
	ip       string
	trans_id string
	asn      int
	asn_name string
}

type na_language struct {
	ip             string
	trans_id       string
	primary_lang   string
	secondary_lang string
}

type na_proxy struct {
	ip                string
	trans_id          string
	proxy_type        string
	proxy_description string
}

type na_is_an_isp struct {
	ip        string
	trans_id  string
	is_an_isp string
}

type na_company struct {
	ip           string
	trans_id     string
	company_name string
}

type na_demographics struct {
	ip         string
	trans_id   string
	rank       int
	households int
	women      int
	w18_34     int
	w35_49     int
	men        int
	m18_34     int
	m35_49     int
	teens      int
	kids       int
}

type na_naics struct {
	ip         string
	trans_id   string
	naics_code int
}

type na_cbsa struct {
	ip         string
	trans_id   string
	cbsa_code  int
	cbsa_title string
	cbsa_type  string
	csa_code   int
	csa_title  string
	md_code    int
	md_title   string
}

type na_mobile_carrier struct {
	ip                  string
	trans_id            string
	mobile_carrier      string
	mcc                 string
	mnc                 string
	mobile_carrier_code string
}

type na_organization struct {
	ip                string
	trans_id          string
	organization_name string
}

type na_pulse struct {
	ip                       string
	trans_id                 string
	pulse_country            string
	pulse_region             string
	pulse_city               string
	pulse_conn_speed         string
	pulse_conn_type          string
	pulse_metro_code         int
	pulse_latitude           float64
	pulse_longitude          float64
	pulse_postal_code        string
	pulse_country_code       int
	pulse_region_code        int
	pulse_city_code          int
	pulse_continent_code     int
	pulse_two_letter_country string
	pulse_internal_code      int
	pulse_area_codes         string
	pulse_country_conf       int
	pulse_region_conf        int
	pulse_city_conf          int
	pulse_postal_conf        int
	pulse_gmt_offset         int
	pulse_in_dst             string
}

type na_pulseplus struct {
	ip                           string
	trans_id                     string
	pulseplus_country            string
	pulseplus_region             string
	pulseplus_city               string
	pulseplus_conn_speed         string
	pulseplus_conn_type          string
	pulseplus_metro_code         int
	pulseplus_latitude           float64
	pulseplus_longitude          float64
	pulseplus_postal_code        string
	pulseplus_postal_code_ext    string
	pulseplus_country_code       int
	pulseplus_region_code        int
	pulseplus_city_code          int
	pulseplus_continent_code     int
	pulseplus_two_letter_country string
	pulseplus_internal_code      int
	pulseplus_area_codes         string
	pulseplus_country_conf       int
	pulseplus_region_conf        int
	pulseplus_city_conf          int
	pulseplus_postal_conf        int
	pulseplus_gmt_offset         int
	pulseplus_in_dst             string
}

/**
Purpose : Parses a non-xml response and returns a corresponding struct for the feature code
Params :
	rawResponse string - A string representation of the response sent from the netacuity server
	featureCode int - The feature code queried
Returns :
	interface{} - Using the empty interface in order to return different structs for different feature codes,
		a string error of the error field in the netacuity response
**/
func parseResponse(rawResponse string, featureCode int) interface{} {
	sliceArray := strings.Split(rawResponse, ";")
	if len(sliceArray) < 4 {
		return fmt.Sprintf("Error in packet received data [%v] len less than 4.", sliceArray)
	}
	var error string = sliceArray[3]
	if error != "" {
		return fmt.Sprintf("Error in packet received : %v", error)
	}

	switch featureCode {
	case 3:
		geo := na_geo{ip: sliceArray[1],
			trans_id:               sliceArray[2],
			geo_country:            sliceArray[4],
			geo_region:             sliceArray[5],
			geo_city:               sliceArray[6],
			geo_conn_speed:         sliceArray[7],
			geo_country_conf:       wrap(strconv.Atoi(sliceArray[8]))[0].(int),
			geo_region_conf:        wrap(strconv.Atoi(sliceArray[9]))[0].(int),
			geo_city_conf:          wrap(strconv.Atoi(sliceArray[10]))[0].(int),
			geo_metro_code:         wrap(strconv.Atoi(sliceArray[11]))[0].(int),
			geo_latitude:           wrap(strconv.ParseFloat(sliceArray[12], 64))[0].(float64),
			geo_longitude:          wrap(strconv.ParseFloat(sliceArray[13], 64))[0].(float64),
			geo_country_code:       wrap(strconv.Atoi(sliceArray[14]))[0].(int),
			geo_region_code:        wrap(strconv.Atoi(sliceArray[15]))[0].(int),
			geo_city_code:          wrap(strconv.Atoi(sliceArray[16]))[0].(int),
			geo_continent_code:     wrap(strconv.Atoi(sliceArray[17]))[0].(int),
			geo_two_letter_country: sliceArray[18]}
		return geo
	case 4:
		edge := na_edge{ip: sliceArray[1],
			trans_id:                sliceArray[2],
			edge_country:            sliceArray[4],
			edge_region:             sliceArray[5],
			edge_city:               sliceArray[6],
			edge_conn_speed:         sliceArray[7],
			edge_metro_code:         wrap(strconv.Atoi(sliceArray[8]))[0].(int),
			edge_latitude:           wrap(strconv.ParseFloat(sliceArray[9], 64))[0].(float64),
			edge_longitude:          wrap(strconv.ParseFloat(sliceArray[10], 64))[0].(float64),
			edge_postal_code:        sliceArray[11],
			edge_country_code:       wrap(strconv.Atoi(sliceArray[12]))[0].(int),
			edge_region_code:        wrap(strconv.Atoi(sliceArray[13]))[0].(int),
			edge_city_code:          wrap(strconv.Atoi(sliceArray[14]))[0].(int),
			edge_continent_code:     wrap(strconv.Atoi(sliceArray[15]))[0].(int),
			edge_two_letter_country: sliceArray[16],
			edge_internal_code:      wrap(strconv.Atoi(sliceArray[17]))[0].(int),
			edge_area_codes:         sliceArray[18],
			edge_country_conf:       wrap(strconv.Atoi(sliceArray[19]))[0].(int),
			edge_region_conf:        wrap(strconv.Atoi(sliceArray[20]))[0].(int),
			edge_city_conf:          wrap(strconv.Atoi(sliceArray[21]))[0].(int),
			edge_postal_conf:        wrap(strconv.Atoi(sliceArray[22]))[0].(int),
			edge_gmt_offset:         wrap(strconv.Atoi(sliceArray[23]))[0].(int),
			edge_in_dst:             sliceArray[24]}
		return edge
	case 5:
		sic := na_sic{ip: sliceArray[1],
			trans_id: sliceArray[2],
			sic_code: wrap(strconv.Atoi(sliceArray[4]))[0].(int)}
		return sic
	case 6:
		domain := na_domain{ip: sliceArray[1],
			trans_id:    sliceArray[2],
			domain_name: sliceArray[4]}
		return domain
	case 7:
		zip := na_zip{ip: sliceArray[1],
			trans_id:      sliceArray[2],
			area_code:     wrap(strconv.Atoi(sliceArray[4]))[0].(int),
			zip_code:      wrap(strconv.Atoi(sliceArray[5]))[0].(int),
			gmt_offset:    wrap(strconv.Atoi(sliceArray[6]))[0].(int),
			in_dst:        sliceArray[7],
			zip_code_text: sliceArray[8],
			zip_country:   sliceArray[9]}
		return zip
	case 8:
		isp := na_isp{ip: sliceArray[1],
			trans_id: sliceArray[2],
			isp_name: sliceArray[4]}
		return isp
	case 9:
		homebiz := na_home_biz{ip: sliceArray[1],
			trans_id:     sliceArray[2],
			homebiz_type: sliceArray[4]}
		return homebiz
	case 10:
		asn := na_asn{ip: sliceArray[1],
			trans_id: sliceArray[2],
			asn:      wrap(strconv.Atoi(sliceArray[4]))[0].(int),
			asn_name: sliceArray[5]}
		return asn
	case 11:
		language := na_language{ip: sliceArray[1],
			trans_id:       sliceArray[2],
			primary_lang:   sliceArray[4],
			secondary_lang: sliceArray[5]}
		return language
	case 12:
		proxy := na_proxy{ip: sliceArray[1],
			trans_id:          sliceArray[2],
			proxy_type:        sliceArray[4],
			proxy_description: sliceArray[5]}
		return proxy
	case 14:
		is_an_isp := na_is_an_isp{ip: sliceArray[1],
			trans_id:  sliceArray[2],
			is_an_isp: sliceArray[4]}
		return is_an_isp
	case 15:
		company := na_company{ip: sliceArray[1],
			trans_id:     sliceArray[2],
			company_name: sliceArray[4]}
		return company
	case 17:
		demographics := na_demographics{ip: sliceArray[1],
			trans_id:   sliceArray[2],
			rank:       wrap(strconv.Atoi(sliceArray[4]))[0].(int),
			households: wrap(strconv.Atoi(sliceArray[5]))[0].(int),
			women:      wrap(strconv.Atoi(sliceArray[6]))[0].(int),
			w18_34:     wrap(strconv.Atoi(sliceArray[7]))[0].(int),
			w35_49:     wrap(strconv.Atoi(sliceArray[8]))[0].(int),
			men:        wrap(strconv.Atoi(sliceArray[9]))[0].(int),
			m18_34:     wrap(strconv.Atoi(sliceArray[10]))[0].(int),
			m35_49:     wrap(strconv.Atoi(sliceArray[11]))[0].(int),
			teens:      wrap(strconv.Atoi(sliceArray[12]))[0].(int),
			kids:       wrap(strconv.Atoi(sliceArray[13]))[0].(int)}
		return demographics
	case 18:
		naics := na_naics{ip: sliceArray[1],
			trans_id:   sliceArray[2],
			naics_code: wrap(strconv.Atoi(sliceArray[4]))[0].(int)}
		return naics
	case 19:
		cbsa := na_cbsa{ip: sliceArray[1],
			trans_id:   sliceArray[2],
			cbsa_code:  wrap(strconv.Atoi(sliceArray[4]))[0].(int),
			cbsa_title: sliceArray[5],
			cbsa_type:  sliceArray[6],
			csa_code:   wrap(strconv.Atoi(sliceArray[7]))[0].(int),
			csa_title:  sliceArray[8],
			md_code:    wrap(strconv.Atoi(sliceArray[9]))[0].(int),
			md_title:   sliceArray[10]}
		return cbsa
	case 24:
		mobile_carrier := na_mobile_carrier{ip: sliceArray[1],
			trans_id:            sliceArray[2],
			mobile_carrier:      sliceArray[4],
			mcc:                 sliceArray[5],
			mnc:                 sliceArray[6],
			mobile_carrier_code: sliceArray[7]}
		return mobile_carrier
	case 25:
		organization := na_organization{ip: sliceArray[1],
			trans_id:          sliceArray[2],
			organization_name: sliceArray[4]}
		return organization
	case 26:
		pulse := na_pulse{ip: sliceArray[1],
			trans_id:                 sliceArray[2],
			pulse_country:            sliceArray[4],
			pulse_region:             sliceArray[5],
			pulse_city:               sliceArray[6],
			pulse_conn_speed:         sliceArray[7],
			pulse_conn_type:          sliceArray[8],
			pulse_metro_code:         wrap(strconv.Atoi(sliceArray[9]))[0].(int),
			pulse_latitude:           wrap(strconv.ParseFloat(sliceArray[10], 64))[0].(float64),
			pulse_longitude:          wrap(strconv.ParseFloat(sliceArray[11], 64))[0].(float64),
			pulse_postal_code:        sliceArray[12],
			pulse_country_code:       wrap(strconv.Atoi(sliceArray[13]))[0].(int),
			pulse_region_code:        wrap(strconv.Atoi(sliceArray[14]))[0].(int),
			pulse_city_code:          wrap(strconv.Atoi(sliceArray[15]))[0].(int),
			pulse_continent_code:     wrap(strconv.Atoi(sliceArray[16]))[0].(int),
			pulse_two_letter_country: sliceArray[17],
			pulse_internal_code:      wrap(strconv.Atoi(sliceArray[18]))[0].(int),
			pulse_area_codes:         sliceArray[19],
			pulse_country_conf:       wrap(strconv.Atoi(sliceArray[20]))[0].(int),
			pulse_region_conf:        wrap(strconv.Atoi(sliceArray[21]))[0].(int),
			pulse_city_conf:          wrap(strconv.Atoi(sliceArray[22]))[0].(int),
			pulse_postal_conf:        wrap(strconv.Atoi(sliceArray[23]))[0].(int),
			pulse_gmt_offset:         wrap(strconv.Atoi(sliceArray[24]))[0].(int),
			pulse_in_dst:             sliceArray[25]}
		return pulse
	case 30:
		pulseplus := na_pulseplus{ip: sliceArray[1],
			trans_id:                     sliceArray[2],
			pulseplus_country:            sliceArray[4],
			pulseplus_region:             sliceArray[5],
			pulseplus_city:               sliceArray[6],
			pulseplus_conn_speed:         sliceArray[7],
			pulseplus_conn_type:          sliceArray[8],
			pulseplus_metro_code:         wrap(strconv.Atoi(sliceArray[9]))[0].(int),
			pulseplus_latitude:           wrap(strconv.ParseFloat(sliceArray[10], 64))[0].(float64),
			pulseplus_longitude:          wrap(strconv.ParseFloat(sliceArray[11], 64))[0].(float64),
			pulseplus_postal_code:        sliceArray[12],
			pulseplus_postal_code_ext:    sliceArray[13],
			pulseplus_country_code:       wrap(strconv.Atoi(sliceArray[14]))[0].(int),
			pulseplus_region_code:        wrap(strconv.Atoi(sliceArray[15]))[0].(int),
			pulseplus_city_code:          wrap(strconv.Atoi(sliceArray[16]))[0].(int),
			pulseplus_continent_code:     wrap(strconv.Atoi(sliceArray[17]))[0].(int),
			pulseplus_two_letter_country: sliceArray[18],
			pulseplus_internal_code:      wrap(strconv.Atoi(sliceArray[19]))[0].(int),
			pulseplus_area_codes:         sliceArray[20],
			pulseplus_country_conf:       wrap(strconv.Atoi(sliceArray[21]))[0].(int),
			pulseplus_region_conf:        wrap(strconv.Atoi(sliceArray[22]))[0].(int),
			pulseplus_city_conf:          wrap(strconv.Atoi(sliceArray[23]))[0].(int),
			pulseplus_postal_conf:        wrap(strconv.Atoi(sliceArray[24]))[0].(int),
			pulseplus_gmt_offset:         wrap(strconv.Atoi(sliceArray[25]))[0].(int),
			pulseplus_in_dst:             sliceArray[26]}
		return pulseplus

	default:
		return nil
	}
}

/**
Purpose : Wraps multiple return values from a function so that each individual value
	can be accessed without having to use the assignment operator for multi-value returns. The compiler
	will take all the arguments and automatically return them in a slice for easy access.
	This feature is unavailable in GO as of August 2016, on GO v1.6
Params :
	x ...interface{} - Any number of any type
Returns :
	[]interface{} - A slice of the return values provided
**/
func wrap(x ...interface{}) []interface{} {
	return x
}

/**
Purpose : Generates an alphanumeric transaction ID of specified length
Params :
	length int - The length of the transaction ID
Returns :
	string - The random transaction ID
**/
func generateTransactionId(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

/**
Purpose : Checks to see if a feature code is legitimate
Params :
	code int - The feature code to be checked
Returns :
	bool - True or false if the feature code is legitimate
**/
func checkFeatureCode(code int) bool {
	if strings.Compare("", FeatureCodeEnums[code]) == 0 {
		return false
	}
	return true
}

/**
Purpose : Checks to see if a comma delimited feature code string contains legitimate feature codes
Params :
	codes string - A comma delimited string of feature codes to be checked
Returns :
	bool - True or false if the feature codes are acceptable
	[]string - A string slice of the feature codes
**/
func checkXmlFeatureCodes(codes string) (bool, []string) {
	sliceArray := strings.Split(codes, ",")
	for _, code := range sliceArray {
		if !checkFeatureCode(wrap(strconv.Atoi(code))[0].(int)) {
			return false, nil
		}
	}
	return true, sliceArray
}

/**
Purpose : Checks to see if an API ID is legitimate
Params :
	id int - The API ID to be checked
Returns :
	bool - True or false if the id is acceptable
**/
func checkApiId(id int) bool {
	return (id >= 0 && id <= 127)
}

/**
Purpose : Checks to see if an IP string is a legitimate IP address
Params :
	ip string - The string form of the IP to be checked
Returns :
	int - 4 (IPv4 type)
		  6 (IPv6 type)
		 -1 (invalid)
**/
func checkIpAddress(ip string) int {
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return -1
	} else {
		if strings.Contains(ip, ":") {
			return 6
		} else {
			return 4
		}
	}
}

/**
Purpose : Checks to see if an int delay is legitimate
Params :
	delay int - The delay to be checked
Returns :
	bool - True or false if the delay is acceptable
**/
func checkTimeoutDelay(delay int) bool {
	if delay >= 0 {
		return true
	}
	return false
}
