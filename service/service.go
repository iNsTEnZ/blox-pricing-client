package service

import (
	"blox-client-service/cache"
	"blox-client-service/config"
	pb "blox-client-service/proto"
	"blox-client-service/util"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type Service struct {
	config   config.ServiceConfig
	cache    cache.ICache
	ticker   *time.Ticker
	client   pb.CryptoPricingClient
	conn     *grpc.ClientConn
	quitChan chan struct{}
}

func New() *Service {
	cfg := config.GetServiceConfig()

	return &Service{
		config:   cfg,
		cache:    cache.NewLRUCache(cfg.MaxCacheSize),
		ticker:   time.NewTicker(cfg.FetchInterval),
		quitChan: make(chan struct{}),
	}
}

func (s *Service) Start() {
	conn, err := grpc.Dial(s.config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	s.conn = conn
	s.client = pb.NewCryptoPricingClient(conn)

	s.waitForConnection()

	s.getTokenPrice()
	log.Printf("starting ticker for token fetching")

	for {
		select {
		case <-s.ticker.C:
			s.waitForConnection()
			s.getTokenPrice()
		case <-s.quitChan:
			s.ticker.Stop()
			s.conn.Close()
			return
		}
	}
}

func (s *Service) Stop() {

}

func (s *Service) getTokenPrice() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	for _, token := range s.config.Tokens {
		response, err := s.client.GetPrice(ctx, &pb.PriceRequest{
			Currency: s.config.Currency,
			Symbols:  []string{token},
		})

		if err != nil {
			log.Printf("could not get price for token %s: %v", token, err)
			continue
		}

		s.printStats(response)
		s.saveToCache(response)
	}
}

func (s *Service) printStats(newData *pb.PriceResponse) {
	for _, symbol := range newData.Symbols {
		if val, exist := s.cache.Get(util.CacheKey(symbol)); exist {
			previousData := val.(*pb.Symbol)
			diff := ((symbol.Price - previousData.Price) / previousData.Price) * 100
			log.Printf("%s: %f %f %.02f%%", symbol.Name, previousData.Price, symbol.Price, diff)
		}
	}
}

func (s *Service) saveToCache(newData *pb.PriceResponse) {
	for _, symbol := range newData.Symbols {
		s.cache.Set(util.CacheKey(symbol), symbol)
	}
}

func (s *Service) waitForConnection() {
	for {
		if s.conn.GetState() == connectivity.Ready {
			break
		}

		s.conn.Connect()
		log.Printf("Waiting for connection to be ready...")
		time.Sleep(2 * time.Second)
	}
}