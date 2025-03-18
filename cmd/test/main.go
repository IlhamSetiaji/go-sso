package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/scheduler"
	"app/go-sso/internal/rabbitmq"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)

	go rabbitmq.InitProducer(viperConfig, log)
	go rabbitmq.InitConsumer(viperConfig, log)

	// scheduler factory
	syncScheduler := scheduler.SyncMidsuitSchedulerFactory(viperConfig, log)
	authResp, err := syncScheduler.AuthOneStep()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	log.Printf("Auth response: %v", authResp)

	err = syncScheduler.SyncOrganizationType(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync organization type: %v", err)
	}
	log.Infof("Successfully synced organization type")

	err = syncScheduler.SyncOrganization(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync organization: %v", err)
	}
	log.Infof("Successfully synced organization")

	err = syncScheduler.SyncJobLevel(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync job level: %v", err)
	}
	log.Infof("Successfully synced job level")

	err = syncScheduler.SyncOrganizationLocation(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync organization location: %v", err)
	}
	log.Infof("Successfully synced organization location")

	err = syncScheduler.SyncOrganizationStructure(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync organization structure: %v", err)
	}
	log.Infof("Successfully synced organization structure")

	err = syncScheduler.SyncJob(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync job: %v", err)
	}
	log.Infof("Successfully synced job")

	err = syncScheduler.SyncEmployee(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync employee: %v", err)
	}
	log.Infof("Successfully synced employee")

	err = syncScheduler.SyncEmployeeJob(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync employee job: %v", err)
	}
	log.Infof("Successfully synced employee job")

	err = syncScheduler.SyncUserProfile(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync user profile: %v", err)
	}
	log.Infof("Successfully synced user profile")

	log.Printf("Successfully synced data")
}
