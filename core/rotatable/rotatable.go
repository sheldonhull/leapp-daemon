package rotatable

type Rotatable interface {
	Rotate(*string) error
	IsRotationIntervalExpired() (bool, error)
}
