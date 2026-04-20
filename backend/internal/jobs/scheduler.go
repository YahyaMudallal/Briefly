package jobs

import (
	"context"
	"log"

	"github.com/YahyaMudallal/newsWebSite/internal/services"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
    cron           *cron.Cron
    articleService *services.ArticleService
}

func NewScheduler(articleService *services.ArticleService) *Scheduler {
    c := cron.New(cron.WithSeconds()) 
    return &Scheduler{
        cron:           c,
        articleService: articleService,
    }
}

func (s *Scheduler) Start() {
	// planify the job to run every day at 2:00 AM
    _, err := s.cron.AddFunc("0 0 2 * * *", func() {
        log.Println("Start the daily articles synchronization...")
        
        ctx := context.Background()
        err := s.articleService.SyncDailyArticles(ctx)
        
        if err != nil {
            log.Printf("Error during synchronization : %v\n", err)
        } else {
            log.Println("Synchronization completed successfully !")
        }
    })

    if err != nil {
        log.Fatalf("Error during CRON configuration : %v", err)
    }

    s.cron.Start()
}

func (s *Scheduler) Stop() {
    s.cron.Stop()
}