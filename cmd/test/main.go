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

	err = syncScheduler.SyncOrganization(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync organization: %v", err)
	}

	err = syncScheduler.SyncJobLevel(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync job level: %v", err)
	}

	err = syncScheduler.SyncOrganizationLocation(authResp.Token)
	if err != nil {
		log.Fatalf("Failed to sync organization location: %v", err)
	}

	log.Printf("Successfully synced data")
}
