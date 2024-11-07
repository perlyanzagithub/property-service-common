package services

import "github.com/perlyanzagithub/property-service-common/utils"

type JWTService struct {
	tokenUtil *utils.TokenUtil
}

func NewJWTService(tokenUtil *utils.TokenUtil) *JWTService {
	return &JWTService{tokenUtil: tokenUtil}
}

func (s *JWTService) CreateToken(userID map[string]interface{}) (string, error) {
	return s.tokenUtil.GenerateToken(userID)
}

func (s *JWTService) ParseToken(tokenStr string) (map[string]interface{}, error) {
	claims, err := s.tokenUtil.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
