package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/scheduler"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)

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

	log.Printf("Successfully synced data")
}
