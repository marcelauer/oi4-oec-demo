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

// ==================================
// Debugging
// ==================================
// Use values from according IODD file of your sensor (https://ioddfinder.io-link.com/)
var DebugParameter = []sensor.IOLinkParameter{
	{
		ID:           "DebugParameter",        // identifier for the parameter
		Index:        123,                     // Index in the IODD
		AccessRights: "ro",                    // Access rights, e.g. "rw" for read/write, "ro" for read-only
		DataType:     "UIntegerT(16)",         // Data type, e.g. "UIntegerT(8)", "Float32T", "StringT(16)"
		DefaultValue: "",                      // Default value, e.g. "0" for integers, "0.0" for floats, or "" for strings
		Name:         "Device type",           // Name of the parameter
		Description:  "Shows the device type", // Description of the parameter
	},
}

// ==================================

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

	unit := "°C"
	response, rErr := app.sensorService.GetSensorProcessData(&unit)

	if rErr != nil {
		app.logger.Warn("Failed to get sensor data:", rErr)

		return nil
	}

	pv := api.NewOi4Data(response)
	app.applicationSource.UpdateData(pv, "eh_values")
	app.applicationSource.UpdateHealth(api.Health{Health: api.Health_Normal, HealthScore: 100})

	// Get sensor parameters
	// TODO: should be tested with a real sensor before adding all parameters
	//response_parametersSensor, rErr := app.sensorService.GetSensorParameterData(parametersSensor, "Sensor Parameters")
	response_parametersSensor, rErr := app.sensorService.GetSensorParameterData(DebugParameter, "Debug Parameter")
	if rErr != nil {
		app.logger.Warn("Failed to get sensor parameters:", rErr)

		return nil
	}
	parSensor := api.NewOi4Data(response_parametersSensor)

	return []api.Data{pv, parSensor}
}

func newDataPublication(application api.Oi4Application, oi4Source api.BaseSource) *publication.IntervalPublicationImpl {
	return publication.NewIntervalBuilder(application, 1*time.Minute). //
										Oi4Source(oi4Source).                                      //
										Resource(api.ResourceData).                                //
										Filter(api.NewFilter("Oi4Data")).                          //
										PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
										Build()
}
