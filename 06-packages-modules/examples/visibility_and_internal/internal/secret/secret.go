package secret

// Unexported variable (scoped to package secret)
var privateSalt = "super-secret-salt"

// Exported function (callable within visibility_and_internal parent tree)
func EncryptToken(raw string) string {
	return "[" + raw + ":" + privateSalt + "]"
}
