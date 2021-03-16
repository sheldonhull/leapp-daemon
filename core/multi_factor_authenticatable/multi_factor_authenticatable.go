package multi_factor_authenticatable

type MultiFactorAuthenticatable interface {
	isMfaRequired() (bool, error)
}
