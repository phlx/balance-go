package exchangerates

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"balance/internal/redis"
	"balance/internal/utils"
)

type Currency string

const (
	RedisCachedKeyPattern   = "currency:date=%s"
	RedisCacheRatesDuration = 24 * time.Hour

	AUD Currency = "AUD"
	BGN Currency = "BGN"
	BRL Currency = "BRL"
	CAD Currency = "CAD"
	CHF Currency = "CHF"
	CNY Currency = "CNY"
	CZK Currency = "CZK"
	DKK Currency = "DKK"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	HKD Currency = "HKD"
	HRK Currency = "HRK"
	HUF Currency = "HUF"
	IDR Currency = "IDR"
	ILS Currency = "ILS"
	INR Currency = "INR"
	ISK Currency = "ISK"
	JPY Currency = "JPY"
	KRW Currency = "KRW"
	MXN Currency = "MXN"
	MYR Currency = "MYR"
	NOK Currency = "NOK"
	NZD Currency = "NZD"
	PHP Currency = "PHP"
	PLN Currency = "PLN"
	RON Currency = "RON"
	RUB Currency = "RUB"
	SEK Currency = "SEK"
	SGD Currency = "SGD"
	THB Currency = "THB"
	TRY Currency = "TRY"
	USD Currency = "USD"
	ZAR Currency = "ZAR"
)

var (
	Currencies = []Currency{
		AUD, BGN, BRL, CAD, CHF, CNY, CZK, DKK, EUR, GBP, HKD,
		HRK, HUF, IDR, ILS, INR, ISK, JPY, KRW, MXN, MYR, NOK,
		NZD, PHP, PLN, RON, RUB, SEK, SGD, THB, TRY, USD, ZAR,
	}
)

func StringToCurrency(currency string) (*Currency, error) {
	for _, existed := range Currencies {
		upped := strings.ToUpper(currency)
		if upped == string(existed) {
			result := Currency(upped)
			return &result, nil
		}
	}
	return nil, errors.Errorf("Invalid currency: %s", currency)
}

type Service struct {
	client *Client
	redis  redis.Client
}

func NewService(client *Client, redis redis.Client) *Service {
	return &Service{
		client: client,
		redis:  redis,
	}
}

func (s *Service) getCachedRate(ctx context.Context, target Currency) (float64, error) {
	cached, err := s.GetCachedRates(ctx)
	if err != nil {
		return 0, err
	}
	rate, ok := (*cached)[target]
	if !ok {
		return 0, errors.Errorf("Unable to find cached rate for target currency %s", target)
	}
	return rate, nil
}

func (s *Service) GetRate(ctx context.Context, target Currency) (float64, error) {
	cached, err := s.getCachedRate(ctx, target)
	if err != nil && err != goredis.Nil {
		return 0, err
	}
	if err == nil {
		return cached, nil
	}

	response, err := s.client.Latest(ctx)
	if err != nil {
		return 0, err
	}
	rate, ok := response.Rates[target]
	if !ok {
		return 0, errors.Errorf("Unable to find rate for target currency %+v", target)
	}
	err = s.SetCachedRates(ctx, response.Rates)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

func (s *Service) GetCachedRates(ctx context.Context) (*map[Currency]float64, error) {
	cached, err := s.redis.Get(ctx, getRedisKey())
	if err != nil {
		return nil, err
	}
	var rates map[Currency]float64
	err = json.Unmarshal([]byte(cached), &rates)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON from string '%s'", cached)
	}
	return &rates, nil
}

func (s *Service) SetCachedRates(ctx context.Context, rates map[Currency]float64) error {
	ratesJson, err := json.Marshal(rates)
	if err != nil {
		return err
	}
	err = s.redis.Set(ctx, getRedisKey(), ratesJson, RedisCacheRatesDuration)
	if err != nil {
		return err
	}
	return nil
}

func currentDate() string {
	return utils.Now().Format("2006-01-02")
}

func getRedisKey() string {
	return fmt.Sprintf(RedisCachedKeyPattern, currentDate())
}
