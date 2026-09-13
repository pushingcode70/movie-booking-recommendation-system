package utils

import (
	"time"

	"movie-booking/config"

	"github.com/golang-jwt/jwt/v5"
)

func jwtSecret() []byte {
	return []byte(config.AppConfig.JWTSecret)
} //so what does the secret key do?To prove that the JWT was created by your server and hasn't been changed.

/*when your server creates a JWT, it uses the secret key to generate a digital signature based on:
the header,
the payload (claims),
and the secret key.*/

// genratetoken creates a jwt for the given user id
func GenerateToken(userID uint, role string) (string, error) {

	//data stored inside JWT(called claims)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), //valid for one day
	}

	//create a jwt object using the HS256 signing algorithm
	//"When you sign this token later, use the HS256 algorithm."
	// the algo will be used in signed string.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //build a jwt object using claims where we use signing method 256 one from algo family

	// sign the JWT object using the secret key and return the JWT string
	return token.SignedString(jwtSecret())
}

// validateToken verifies the JWT.
func ValidateToken(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { //all that decoding stuff happend inside library it just takes client token string and function which will return secret key and after that it compare signaturesa and other suff internally

		// ensure the token uses an HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		// return the secret key for signature verification
		return jwtSecret(), nil //if its signed with same algo only than return jewt secret key
	})
}

/*think of the secret key as a special stamp that only your server owns.
when a user logs in, the server first creates the header (which contains information like "alg": "HS256" and "typ": "JWT")
and the payload (which contains the data you want to store, such as "user_id": 5 and the expiration time).
 these two pieces of information are not secret—they are simply encoded and can be read by anyone.
 the server then takes the header, the payload, and its secret key and feeds all three into the HS256 (HMAC-SHA256) algorithm.
  the algorithm performs a series of mathematical operations and produces a seemingly random string called the signature.
  this signature is not the secret key; it is the result of using the secret key.
  it is similar to baking a cake: flour, sugar, eggs, and butter go into the oven, and a cake comes out. T
  he cake is made using those ingredients, but the cake itself is not flour or sugar.
   likewise, the signature is created using the secret key, but it is not the secret key itself, nor can you recover the secret key from it.

when the server sends the JWT to the client, it sends only three things: Header.Payload.Signature.
 notice that the secret key is never included in the JWT. It stays safely on the server.
  later, when the client sends the JWT back, the server receives exactly the same three parts.
   it already has the header, the payload, and the signature that came from the client,
    but it still has its own copy of the secret key stored in memory.
	the server then repeats the exact same mathematical process it used when the JWT was first created:
	 it combines the received header and payload with its own secret key and runs the HS256 algorithm again.
	 if nothing has changed, mathematics guarantees that the algorithm will produce exactly the same signature as before.
	 the server then compares the received signature (the one inside the JWT) with the newly computed signature.
	 if they are identical, the server knows that the header and payload have not been altered and that the JWT was
	 originally signed by someone who knew the secret key—namely, the server itself.

now imagine an attacker intercepts the JWT and changes the payload from "user_id": 5 to "user_id": 999.
the attacker can easily modify the payload because it is only encoded, not encrypted.
 however, the attacker does not know the server's secret key.
  therefore, they cannot generate a new valid signature for the modified payload.
  when the server later receives this tampered JWT, it extracts the modified header and payload,
   uses its own secret key to compute what the signature should be, and gets a completely different result.
   the signature inside the JWT is still the old one, so the comparison fails.
   the server immediately knows that someone changed the token after it was signed and rejects it.*/
