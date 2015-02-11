package deploys

import (
	"github.com/remind101/empire/apps"
	"github.com/remind101/empire/configs"
	"github.com/remind101/empire/releases"
	"github.com/remind101/empire/slugs"
)

// Deploy represents a deployment to the platform.
type Deploy struct {
	ID      string
	Status  string
	Release *releases.Release
}

type Service struct {
	AppsService     *apps.Service
	ConfigsService  *configs.Service
	SlugsService    *slugs.Service
	ReleasesService *releases.Service
}

// Deploy deploys an Image to the platform.
func (s *Service) Deploy(image *slugs.Image) (*Deploy, error) {
	app, err := s.AppsService.FindOrCreateByRepo(image.Repo)
	if err != nil {
		return nil, err
	}

	// Grab the latest config.
	config, err := s.ConfigsService.Head(app.ID)
	if err != nil {
		return nil, err
	}

	// Create a new slug for the docker image.
	//
	// TODO This is actually going to be pretty slow, so
	// we'll need to do
	// some polling or events/webhooks here.
	slug, err := s.SlugsService.CreateByImage(image)
	if err != nil {
		return nil, err
	}

	// Create a new release for the Config
	// and Slug.
	release, err := s.ReleasesService.Create(app, config, slug)
	if err != nil {
		return nil, err
	}

	// We're deployed! ...
	// hopefully.
	return &Deploy{
		ID:      "1",
		Release: release,
	}, nil
}
