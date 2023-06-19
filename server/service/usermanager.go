package service

import (
	"fmt"
	"rabbitnft/server/config"
	"rabbitnft/server/message"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

var tokenKey = "malimalihome"

type UserManager struct {
	accSvc  *AccountSvc
	userMap map[string]UserInfo
	mu      sync.RWMutex
	startAt time.Time
}

func NewUserManager(accSvc *AccountSvc) *UserManager {
	return &UserManager{
		accSvc:  accSvc,
		userMap: map[string]UserInfo{},
		mu:      sync.RWMutex{},
		startAt: time.Now(),
	}
}

func (u *UserManager) UserLogin(account *message.Account) (string, error) {
	dbAccount, err := u.accSvc.getDBAccount(account.Account, account.Password)
	if err != nil {
		return "", err
	}
	accountInfo := message.AccountInfo{
		Account:  dbAccount.Account,
		Nickname: dbAccount.Username,
		Address:  dbAccount.Address,
	}
	user := &UserInfo{
		Private: *dbAccount,
		Basic:   accountInfo,
	}
	u.mu.Lock()
	u.userMap[account.Account] = *user
	u.mu.Unlock()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account.Account,
		"address": dbAccount.Address,
		"created": time.Now().Format(time.RFC3339),
	})
	tokenStr, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", fmt.Errorf("token sign failed: %s", err)
	}
	return tokenStr, err
}

func (u *UserManager) UserLogout(account string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, ok := u.userMap[account]
	if !ok {
		return fmt.Errorf("User not exist")
	}
	delete(u.userMap, account)
	return nil
}

func (u *UserManager) ValidateToken(tokenStr string) error {
	claims, err := u.readTokenClaims(tokenStr)
	if err != nil {
		return err
	}
	createdStr, ok := claims["created"].(string)
	if !ok {
		return fmt.Errorf("Read token created time failed.")
	}
	created, err := time.Parse(time.RFC3339, createdStr)
	if err != nil {
		return fmt.Errorf("parse token created time failed.")
	}
	live := config.C.GetDuration("token.live")
	if limit := created.Add(live * time.Minute); limit.Before(time.Now()) {
		return fmt.Errorf("Token expired.")
	}
	if created.Before(u.startAt) {
		return fmt.Errorf("Token expired.")
	}
	return nil
}

func (u *UserManager) GetUserInfoByToken(tokenStr string) (*UserInfo, error) {
	claims, err := u.readTokenClaims(tokenStr)
	if err != nil {
		return nil, err
	}
	account, ok := claims["account"]
	if !ok {
		return nil, fmt.Errorf("Account not exist")
	}
	accountStr, ok := account.(string)
	if !ok {
		return nil, fmt.Errorf("Account assertion failed")
	}

	userInfo, ok := u.userMap[accountStr]
	if !ok {
		return nil, fmt.Errorf("user info not exist")
	}
	balance, err := u.accSvc.ethSDK.GetBalance(userInfo.Basic.Address)
	if err != nil {
		return nil, err
	}
	userInfo.Basic.Balance = balance
	return &userInfo, nil
}

// private methods

func (u *UserManager) readTokenClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token not valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Token claim failed")
	}
	return claims, nil
}
