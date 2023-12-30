package service

import (
	"database/sql"
	"log"
	"time"

	"github.com/Vainsberg/discounts-telegram-bot/internal/bottg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	cronhandler "github.com/Vainsberg/discounts-telegram-bot/internal/cron_handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	RepositoryQuerys     *repository.RepositorySubs
	SubsRepository       *repository.RepositorySubs
	Logger               *zap.Logger
	DiscountsPlatiClient client.PlatiClient
	DiscountsRepository  repository.Repository
	RequestCron          bottg.RequestCron
}

func NewService(db *sql.DB) *Service {
	RepoSubs := repository.NewRepositorySubs(db)
	if RepoSubs == nil {
		log.Fatal("Failed to initialize RepositorySubs")
	}

	return &Service{
		RepositoryQuerys:     RepoSubs,
		SubsRepository:       RepoSubs,
		Logger:               &zap.Logger{},
		DiscountsPlatiClient: client.PlatiClient{},
	}
}

func (s *Service) ProcessQueryAndFetchGoods() {

	query, err := s.SubsRepository.GetQuerys()
	if err != nil {
		s.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
		return
	}

	for _, v := range query.Items {
		Query := pkg.ReplaceSpaceUrl(v.Query)
		s.Logger.Info("Sleep 5 second...", zap.String("query", Query))
		time.Sleep(5 * time.Second)
		goods, err := s.DiscountsPlatiClient.GetGoodsClient(Query)
		if err != nil {
			s.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
			return
		}

		for _, el := range goods.Items {
			s.DiscountsRepository.SaveGood(el.Name, float64(el.Price_rur), el.Url, el.Image, Query)
			//pastPrice := h.DiscountsRepository.SearchAveragePrice(float64(el.Price_rur), el.Url)
			if true {
				//rebate := pkg.CalculatePercentageDifference(pastPrice, float64(el.Price_rur))
				if true {
					s.Logger.Info("Finding the updated price", zap.Int("Price_rur", el.Price_rur))
					chat := s.SubsRepository.SearchChatID(Query)
					goodDiscounts := cronhandler.ProductDiscounts(el.Name, float64(el.Price_rur), el.Url, el.Image)
					s.RequestCron.HandleRequestCron(&goodDiscounts, chat)
				}
			}
		}
	}

}
