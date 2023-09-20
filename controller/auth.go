package controller

import (
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

type TotpController struct{}

func (uc *TotpController) GenerateTotpKeyController(conn *gorm.DB, msg *protosv1.GenerateTotpKeyRequest) (*protosv1.GenerateTotpKeyResponse, error) {
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

func (uc *TotpController) VerifyTotpController(conn *gorm.DB, msg *protosv1.VerifyTotpRequest) (*protosv1.VerifyTotpResponse, error) {
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

func (uc *TotpController) ValidateTotpController(conn *gorm.DB, msg *protosv1.ValidateTotpRequest) (*protosv1.ValidateTotpResponse, error) {
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
