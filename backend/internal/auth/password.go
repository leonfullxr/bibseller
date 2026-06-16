package auth

import "github.com/alexedwards/argon2id"

// hashParams is the OWASP Password Storage Cheat Sheet's first recommended
// argon2id configuration: m=19 MiB, t=2, p=1.
//
// Why argon2id: it is *memory-hard* - each guess costs 19 MiB of RAM, which
// is what actually slows down GPU/ASIC cracking rigs (their cores are many,
// but their fast memory is not). The "id" variant mixes the data-independent
// pass of argon2i (side-channel resistance) with the data-dependent pass of
// argon2d (best cracking resistance).
//
// These parameters can be raised later without breaking existing hashes:
// every hash is stored as a PHC string - "$argon2id$v=19$m=19456,t=2,p=1$
// <salt>$<key>" - that embeds its own parameters and salt, and verification
// reads them from the string, not from this struct.
var hashParams = &argon2id.Params{
	Memory:      19 * 1024, // KiB -> 19 MiB per hash attempt
	Iterations:  2,         // passes over that memory
	Parallelism: 1,         // lanes; raise only with dedicated cores
	SaltLength:  16,        // 128-bit salt, random per hash (thwarts rainbow tables)
	KeyLength:   32,        // 256-bit derived key
}

func hashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, hashParams)
}

// verifyPassword re-derives the key using the salt+params embedded in the
// stored PHC string and compares in constant time (inside the library).
func verifyPassword(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

// dummyHash is verified against when login hits an unknown email, so both
// failure paths cost one argon2id computation. Without it, "unknown email"
// returns in microseconds while "wrong password" takes tens of milliseconds -
// a timing oracle that lets an attacker enumerate which emails have accounts.
// Computed once at startup so it always matches the current hashParams.
var dummyHash = func() string {
	h, err := argon2id.CreateHash("dummy-password-for-timing-equalization", hashParams)
	if err != nil {
		panic(err) // CSPRNG failure at startup: nothing sensible to do
	}
	return h
}()
