// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/wwfirebase"

	"github.com/dgrijalva/jwt-go"
	"github.com/op/go-logging"
)

// LOGGER for logging
var LOGGER = logging.MustGetLogger("middlewares")

// HasIP : check if the request ip address is included in token ip array
func HasIP(claims jwt.MapClaims, remoteAddrs []string) bool {

	// split address into ip and port
	//remoteAddrs[index], _, _ = net.SplitHostPort(ip)

	//ip = "192.1.1.1"
	//check with ip specified in token is present in the caller ip url
	ipToken := claims["ips"].([]interface{})
	ipString := make([]string, len(ipToken))
	for _, remoteAddr := range remoteAddrs {
		for i, v := range ipToken {
			ipString[i] = v.(string)
			if strings.Contains(remoteAddr, ipString[i]) == true {
				return true
			}
		}
	}
	return false
}

// HasEndpoint : check if requested endpoint is included in endpoint array
func HasEndpoint(claims jwt.MapClaims, enpTarget string) bool {

	// TODO: from chase to Nakul - need to check if the endpoint trying
	// to be accessed is allows per the permissions.json. I started this
	// check in the commented code below but the parsing and checking
	// still needs to be implemented.

	// // check if the endpoint is listed as an endpoint approved for JWT authentication:
	// // get list of endpoints allowed per
	// filename := "./auth-service/permissions.json"
	// text, _ := ioutil.ReadFile(filename)
	// var data interface{}
	// err := json.Unmarshal(text, &data)
	// fmt.Printf(data)

	// mapJwtEndpoints := authconstants.EndpointJwtPermissions()

	// // TODO: Nakul, This is an LCP problem and can be optimized to an O(n) instead of KMP + O(n).
	// for k := range mapJwtEndpoints {
	// 	if strings.Contains(enpTarget, k) {
	// 		isEndpointValid = true
	// 	}
	// }

	// check if the claims exist on the JWT token required to access the endpoint
	enpToken := claims["enp"].([]interface{})
	enpString := make([]string, len(enpToken))
	for i, v := range enpToken {
		enpString[i] = v.(string)
		if strings.Contains(enpTarget, enpString[i]) == true {
			return true
		}
	}
	return false

}

// IsValid : check token against db, isOnCount and
// func IsValid(jti string, count float64) (bool, error) {

// 	// Indeed, storing all issued JWT IDs undermines the
// 	// stateless nature of using JWTs. However, the purpose
// 	// of JWT IDs is to be able to revoke previously-issued
// 	// JWTs. This can most easily be achieved by blacklisting
// 	// instead of whitelisting. If you've included the "exp"
// 	// claim (you should), then you can eventually clean up
// 	// blacklisted JWTs as they expire naturally. Of course
// 	// you can implement other revocation options alongside
// 	// (e.g. revoke all tokens of one client based on a combination
// 	// of "iat" and "aud").
// 	isMatchingSubStr := false
// 	if jti == encodedToken[len(encodedToken)-8:] {
// 		isMatchingSubStr = true
// 	}

// 	// check if the token count is the same as the db count
// 	isOnCount := false
// 	if data["n"].(float64) == count {
// 		isOnCount = true
// 	}

// 	// check if all validatations pass
// 	if isOnCount &&
// 		isMatchingJTI {
// 		return true, nil
// 	}

// 	return false, errors.New("token is not valid")

// }

// IsForParticipant : if service is not a singleton then check if the token participant matches the participant id per the env var
func IsForParticipant(claims jwt.MapClaims) bool {

	// check if env is set for participantId
	// get the participant id from env var
	participantID, exists := os.LookupEnv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	if exists {
		// skip participant id checking if its a global service identified by "ww"; tentatively code it here.
		if participantID == "ww" {
			return true
		}
		// get audience string array
		aud, _ := claims["aud"].(string)

		// check if the token is associated with the participants id
		if aud == participantID {
			// matches
			return true
		}

		// does not match
		return false
	}

	// return true if participantId is not present
	return true

}

// ExtractJWTClaims : parses (decodes) jwt token using secret and returns claims if successful
func ExtractJWTClaims(tokenStr string, r *http.Request) (jwt.MapClaims, bool) {

	// init data var
	// get secure info from database
	var jwtSecure map[string]interface{}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect...
		// someone could inject an easier alg to solve
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			LOGGER.Debugf("Unexpected signing alg: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// println(string("alg: " + string(token.Header["alg"].(string))))

		if string(token.Header["alg"].(string)) != "HS256" {
			// expected alg does not match alg used HS256
			LOGGER.Debugf("Unexpected signing alg: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing alg: %v", token.Header["alg"])
		}

		// fmt.Println(token.Header["alg"])
		fmt.Println(token.Header["kid"])

		// get key id for looking up secret from firebase db
		dbKey := strings.Split(string(token.Header["kid"].(string)), ".")[0]

		// fmt.Println(dbKey)
		// get the key id for looking up the secret from env vars
		pepperKey := strings.Split(string(token.Header["kid"].(string)), ".")[1]

		if err := wwfirebase.FbRef.Child("/jwt_secure/"+dbKey).
			Get(wwfirebase.AppContext, &jwtSecure); err != nil {
			// return false, errors.New("token authentication system is down")
			LOGGER.Debugf("token authentication system is down, firebase database connection failed and cannot retrieve secret key necessary to decrypt JWT token")
			return nil, fmt.Errorf("token authentication system is down, firebase database connection failed and cannot retrieve secret key necessary to decrypt JWT token")
		}

		if len(jwtSecure) == 0 {
			LOGGER.Debugf("Unable to get jwtSecure from DB, check DB credentials")
			return nil, fmt.Errorf("Unable to get jwtSecure from DB, check DB credentials")
		}

		type Request struct {
			O int8              `json:"o"`
			C int8              `json:"c"`
			V map[string]string `json:"v"`
		}

		sBase64 := os.Getenv(global_environment.ENV_KEY_WW_JWT_PEPPER_OBJ)
		s, err := base64.StdEncoding.DecodeString(sBase64)
		if err != nil {
			return nil, fmt.Errorf("Error Decoding base64 from Pepper Object")
		}

		data := Request{}
		json.Unmarshal(s, &data)
		// fmt.Println(data.C)

		decryptionKey := jwtSecure["s"].(string) + data.V[pepperKey]

		// println(decryptionKey)

		return []byte(decryptionKey), nil
	})

	// error parsing token
	if err != nil {
		LOGGER.Debugf("attempt to parse JWT token failed")
		return nil, false
	}

	// run authorization checks on jwt token claims:
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		ipString := r.Header.Get("X-Forwarded-For") // capitalisation

		// exclude all private IP from the x-forwarded-for header
		var privateIPBlocks []*net.IPNet
		for _, cidr := range []string{
			"127.0.0.0/8",    // IPv4 loopback
			"10.0.0.0/8",     // RFC1918
			"172.16.0.0/12",  // RFC1918
			"192.168.0.0/16", // RFC1918
			"::1/128",        // IPv6 loopback
			"fe80::/10",      // IPv6 link-local
			"fc00::/7",       // IPv6 unique local addr
		} {
			_, block, _ := net.ParseCIDR(cidr)
			privateIPBlocks = append(privateIPBlocks, block)
		}

		// parse the incoming ip list as array
		ips := strings.Split(ipString, ",")
		var filteredIps []string
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			IPAddress := net.ParseIP(ip)
			isPrivateIP := false
			for _, block := range privateIPBlocks {
				if block.Contains(IPAddress) {
					isPrivateIP = true
					break
				}
			}
			if isPrivateIP {
				continue
			}
			filteredIps = append(filteredIps, ip)
		}
		LOGGER.Infof("JWT Validation: Receiving request from %+v", filteredIps)
		// check if the token contains the requested ip
		hasIP := HasIP(claims, filteredIps)

		// check if the requested endpoint is included in token
		hasEndpoint := HasEndpoint(claims, r.RequestURI)

		// // Indeed, storing all issued JWT IDs undermines the
		// // stateless nature of using JWTs. However, the purpose
		// // of JWT IDs is to be able to revoke previously-issued
		// // JWTs. This can most easily be achieved by blacklisting
		// // instead of whitelisting. If you've included the "exp"
		// // claim (you should), then you can eventually clean up
		// // blacklisted JWTs as they expire naturally. Of course
		// // you can implement other revocation options alongside
		// // (e.g. revoke all tokens of one client based on a combination
		// // of "iat" and "aud").
		isMatchingJTI := false
		if jwtSecure["i"] == claims["jti"] {
			isMatchingJTI = true
		}

		//Check if env set in the token matches with current runtime micro-service env
		isMatchingENV := false
		env := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)

		if claims["env"] == env {
			isMatchingENV = true
		}

		// check if the token count is the same as the db count
		isOnCount := false
		if jwtSecure["n"].(float64) == claims["n"].(float64) {
			isOnCount = true
		}

		// check if the token is for this participant's id
		isForParticipant := IsForParticipant(claims)

		// check if all validatations pass
		if isOnCount == true &&
			isMatchingJTI == true &&
			hasEndpoint == true &&
			isForParticipant == true &&
			isMatchingENV == true &&
			hasIP == true {
			// token is valid return true

			return claims, true
		}
		LOGGER.Debugf("JWT token is not valid because either isOnCount:%v, isMatchingJTI:%v, hasEndpoint:%v, isForParticipant:%v, or hasIp:%v check failed: ip =%v, isMatchingENV:%v, currentENV:%v, Tokenenv:%v", isOnCount,
			isMatchingJTI, hasEndpoint, isForParticipant, hasIP, filteredIps[0], isMatchingENV, env, claims["env"])
	}

	// default return invalid
	LOGGER.Debugf("JWT token is not valid because either isOnCount, isMatchingJTI, hasEndpoint, isForParticipant, or hasIp check failed")
	return nil, false
}
