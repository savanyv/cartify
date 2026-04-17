package helpers

import "golang.org/x/crypto/bcrypt"

type BcryptService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}

type bcryptService struct {
	cost int
}

func NewBcryptService() BcryptService {
	return &bcryptService{
		cost: bcrypt.DefaultCost,
	}
}

func (b *bcryptService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (b *bcryptService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}