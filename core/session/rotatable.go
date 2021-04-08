package session

type Rotatable interface {
  Rotate(rotateConfiguration *RotateConfiguration) error
}
