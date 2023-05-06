package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Alg string   `json:"alg"` //Signing algorithm
	Kty string   `json:"kty"` //Key type
	Use string   `json:"use"` //Use case ex. sig = signature verification
	N   string   `json:"n"`   //moduluos for a standard pem?
	E   string   `json:"e"`   //exponent for a standard pem
	Kid string   `json:"kid"` //Key id
	X5t string   `json:"x5t"` //Thumbprint of the x.509 cert?
	X5c []string `json:"x5c"` //The "x5c" (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates?
}

// AuthMiddleware validates RSA encrypted Bearer tokens passed in the authorization headers.
// Awesome blog post! https://auth0.com/blog/navigating-rs256-and-jwks/#TL-DR
func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve token
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			http.Error(w, "Authorization header missing.", http.StatusUnauthorized)
			return
		}

		// Validate different properties of the token and retrieve a public key from auth0
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			issuer := "https://dev-klg37k4khu3qm746.us.auth0.com/"
			audience := "http://localhost:5000/"

			// Validate headers
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				log.Println("Expected alg RS256, received: " + token.Method.Alg())
				http.Error(w, "Not authorized.", http.StatusUnauthorized)
			}

			// Validate claims
			if token.Claims.(jwt.MapClaims)["iss"] != issuer {
				log.Println("Expected issuer: ", issuer, " Got: ", token.Claims.(jwt.MapClaims)["iss"])
				http.Error(w, "Not authorized.", http.StatusUnauthorized)
			}

			var verifiedAudience bool
			for _, v := range token.Claims.(jwt.MapClaims)["aud"].([]interface{}) {
				if v == audience {
					verifiedAudience = true
				}
			}
			if !verifiedAudience {
				log.Println("Expected audience: ", audience, " Got: ", token.Claims.(jwt.MapClaims)["aud"])
				http.Error(w, "Not authorized.", http.StatusUnauthorized)
			}

			// Find matching key
			publicKey, err := getPublicRSAKey("https://dev-klg37k4khu3qm746.us.auth0.com/.well-known/jwks.json", token.Header["kid"].(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			return publicKey, nil
		})

		if err != nil {
			log.Println("Token not authorized: ", err.Error())
			http.Error(w, "Not authorized.", http.StatusUnauthorized)
		}

		if token.Valid {
			log.Println("Token authorized!")
		}

		// Authorized! Serve handler function
		next.ServeHTTP(w, r)
	})
}

// getPublicRSAKey sends get request to the url paramater to retrieve a JSON Web Key Set.
// It then compares the keys with the provided key id(kid) and returns an it as an RSA Public key if a match is found.
func getPublicRSAKey(url string, tokenKid string) (interface{}, error) {
	// Get keySet
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve JSON key set. %v", err)
	}
	defer res.Body.Close()

	var keySet JWKS
	err = json.NewDecoder(res.Body).Decode(&keySet)
	if err != nil {
		return nil, fmt.Errorf("Failed to Decode JSON key set. %v", err)
	}

	// Compare keys to token kid
	var jwk JWK
	for _, v := range keySet.Keys {
		if v.Kid == tokenKid {
			jwk = v
			break
		}
	}
	// Check if match was found or fields are missing
	if jwk.E == "" || jwk.N == "" {
		return nil, fmt.Errorf("Matched key not found or properties missing\n\tkey: %v", jwk)
	}

	// Encode and return matched key
	mod, err := base64.RawURLEncoding.DecodeString(jwk.N)
	exp, err := base64.RawStdEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode modulus or exponent. %v", err)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(mod),
		E: int(new(big.Int).SetBytes(exp).Int64()),
	}, nil
}
