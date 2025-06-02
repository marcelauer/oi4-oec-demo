package application

import (
	"IOLinkConnect/internal/sensor"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/application"
	"github.com/OI4/oi4-oec-service-go/service/application/publication"
	"github.com/OI4/oi4-oec-service-go/service/application/source"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"go.uber.org/zap"
)

type SensorApplication struct {
	*application.Oi4ApplicationImpl
	applicationSource *source.ApplicationSourceImpl
	mam               api.MasterAssetModel
	assets            map[string]AssetsEntry
	storage           *container.Storage
	sensorService     *sensor.Service
	logger            *zap.SugaredLogger
}

func NewSensorApplication(
	mam api.MasterAssetModel,
	storage *container.Storage,
	sensorService *sensor.Service,
	logger *zap.SugaredLogger,
) *SensorApplication {
	applicationSource := source.NewApplicationSourceImpl(mam)
	oi4Application, err := application.CreateNewApplication(api.ServiceTypeOTConnector, applicationSource, logger)

	if err != nil {
		logger.Fatal("Failed to create application:", err)
		panic(err)
	}

	assets := make(map[string]AssetsEntry)

	return &SensorApplication{
		oi4Application,
		applicationSource,
		mam,
		assets,
		storage,
		sensorService,
		logger,
	}
}

func (app *SensorApplication) AddAssets(assetList []Asset) {
	for _, asset := range assetList {
		app.AddAsset(asset)
	}
}

func (app *SensorApplication) AddAsset(asset Asset) {
	key := asset.ToOi4Identifier().ToString()

	option := source.WithDataFn(
		func(_ api.BaseSource, filter *api.Filter) []api.Data {
			return app.getSensorData(asset, filter)
		},
	)

	assetSource := source.NewAssetSourceImpl(asset.MasterAssetModel, option)
	oi4Asset := application.CreateNewAsset(assetSource, app.Oi4ApplicationImpl)
	dataAssetPublication := newDataPublication(app.Oi4ApplicationImpl, assetSource)
	err := oi4Asset.RegisterPublication(dataAssetPublication)

	if err != nil {
		app.logger.Error("Failed to register publication:", err)
	}

	metaDataAssetPublication := publication.NewResourcePublication(
		app.Oi4ApplicationImpl,
		assetSource,
		api.ResourceMetadata)

	err = oi4Asset.RegisterPublication(metaDataAssetPublication)

	if err != nil {
		app.logger.Error("Failed to register publication:", err)
	}

	app.RegisterAsset(oi4Asset)

	assetEntry := AssetsEntry{
		assetSource,
		oi4Asset,
		asset,
	}

	app.assets[key] = assetEntry
}

func (app *SensorApplication) getSensorData(asset Asset, filter *api.Filter) []api.Data {
	if filter != nil && !api.FilterEquals(filter, api.NewFilter("Oi4Data")) {
		return nil
	}

	response, rErr := app.sensorService.GetSensorData()

	if rErr != nil {
		app.logger.Warn("Failed to get sensor data:", rErr)

		return nil
	}

	data := api.NewOi4Data(response)

	return []api.Data{data}
}

func newDataPublication(application api.Oi4Application, oi4Source api.BaseSource) *publication.IntervalPublicationImpl {
	return publication.NewIntervalBuilder(application, 1*time.Minute). //
										Oi4Source(oi4Source).                                      //
										Resource(api.ResourceData).                                //
										Filter(api.NewFilter("Oi4Data")).                          //
										PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
										Build()
}
