package controller

import (
	"errors"
	"fmt"
	"log"
	"nnyd-back/config"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type AuthController struct{}

func (ac *AuthController) GenerateTotpKeyController(conn *gorm.DB, msg *protosv1.GenerateTotpKeyRequest) (*protosv1.GenerateTotpKeyResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "firebase_id = ?", msg.FirebaseId).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if u.OtpSecret != "" || u.OtpUrl != "" {
		err := fmt.Errorf("already set up a Secret")
		log.Println(err)
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.APP_NAME,
		AccountName: u.DisplayId,
		SecretSize:  15,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	userDataToUpdate := db.Users{
		OtpSecret: key.Secret(),
		OtpUrl:    key.URL(),
	}

	conn.Model(&u).Updates(userDataToUpdate)

	log.Println("otp url: " + key.URL())
	log.Println("otp secret: " + key.Secret())

	log.Println("Success totp generate.")

	totpResponse := &protosv1.GenerateTotpKeyResponse{
		TotpKey: key.Secret(),
		TotpUrl: key.URL(),
	}

	return totpResponse, nil
}

func (ac *AuthController) VerifyTotpController(conn *gorm.DB, msg *protosv1.VerifyTotpRequest) (*protosv1.VerifyTotpResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "firebase_id = ?", msg.FirebaseId).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if u.OtpSecret == "" && u.OtpUrl == "" {
		err := fmt.Errorf("otp secret is not set")
		log.Println(err)
		return nil, err
	}

	valid, err := totp.ValidateCustom(msg.Token, u.OtpSecret, time.Now(), totp.ValidateOpts{
		Period:    30,                // トークンの有効期間（デフォルトは30秒）
		Digits:    otp.DigitsSix,     // トークンの桁数
		Algorithm: otp.AlgorithmSHA1, // 使用するハッシュアルゴリズム
		Skew:      0,                 // 許容する時間ステップの範囲（前後のステップを許容しない）
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !valid {
		err := fmt.Errorf("failed verify otp token")
		log.Println(err)
		return nil, err
	}

	userDataToUpdate := db.Users{
		OtpEnabled:  true,
		OtpVerified: true,
	}

	conn.Model(&u).Updates(userDataToUpdate)

	log.Println("Success totp verify.")

	totpResponse := &protosv1.VerifyTotpResponse{
		Status: true,
	}

	return totpResponse, nil
}

func (ac *AuthController) ValidateTotpController(conn *gorm.DB, msg *protosv1.ValidateTotpRequest) (*protosv1.ValidateTotpResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "firebase_id = ?", msg.FirebaseId).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if u.OtpSecret == "" && u.OtpUrl == "" {
		err := fmt.Errorf("otp secret is not set")
		log.Println(err)
		return nil, err
	}

	valid, err := totp.ValidateCustom(msg.Token, u.OtpSecret, time.Now(), totp.ValidateOpts{
		Period:    30,                // トークンの有効期間（デフォルトは30秒）
		Digits:    otp.DigitsSix,     // トークンの桁数
		Algorithm: otp.AlgorithmSHA1, // 使用するハッシュアルゴリズム
		Skew:      0,                 // 許容する時間ステップの範囲（前後のステップを許容しない）
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !valid {
		err := fmt.Errorf("failed verify otp token")
		log.Println(err)
		return nil, err
	}

	log.Println("Success totp verify.")

	totpResponse := &protosv1.ValidateTotpResponse{
		Status: true,
	}

	return totpResponse, nil
}

func (ac *AuthController) CheckVerifyTotp(conn *gorm.DB, firebase_id string) (string, error) {
	u := db.Users{}

	if err := conn.First(&u, "firebase_id = ?", firebase_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user having this firebase_id is not found")
		} else {
			log.Println(err)
			return "", err
		}
	}

	if !u.OtpVerified {
		return "", fmt.Errorf("totp check failed")
	}

	return u.ID.String(), nil
}

func (ac *AuthController) SignOut(conn *gorm.DB, user_id string) (*protosv1.SignOutResponse, error) {
	u := db.Users{}
	if err := conn.First(&u, "id = ?", user_id).Error; err != nil {
		return nil, err
	}

	u.OtpVerified = false

	if err := conn.Save(u).Error; err != nil {
		return nil, err
	}

	return &protosv1.SignOutResponse{
		Status: true,
	}, nil
}
