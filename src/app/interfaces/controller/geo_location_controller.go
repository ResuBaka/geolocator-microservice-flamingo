package controller

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/magento-hackathon/geolocator-microservice-flamingo/src/app/domain"
	"net"
	"net/http"
)

type (
	// GeoLocationController is the app main controller
	GeoLocationController struct {
		responder         *web.Responder
		LocationProviders []domain.LocationProvider
	}
)

func (c *GeoLocationController) Inject(
	locationProviders []domain.LocationProvider,
) {
	c.LocationProviders = locationProviders
}

// GetGeoLocation returns a geolocation for a provided ipaddress param
func (c *GeoLocationController) GetGeoLocation(ctx context.Context, r *web.Request) web.Result {
	ipAddress, err := r.Query1("ipaddress")
	if err != nil {
		return c.responder.Data(nil).Status(http.StatusUnprocessableEntity)
	}

	validIP := net.ParseIP(ipAddress)
	if validIP == nil {
		return c.responder.Data(nil).Status(http.StatusUnprocessableEntity)
	}

	var results []*domain.LocationData

	for _, provider := range c.LocationProviders {
		result, err := provider.GetLocationByIP(validIP)

		if err != nil {
			continue
		}

		results = append(results, result)
	}

	res := c.responder.Data(results).Status(http.StatusOK)
	res.Header.Set("Content-Type", "application/json")

	return res
}
