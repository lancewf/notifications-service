package config

import (
	"io/ioutil"
	"os"

	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	config      aggregateConfig
	configFile  string
	updateQueue chan<- update
}

type update struct {
	fun  func(aggregateConfig) (aggregateConfig, error)
	errc chan error
}

type aggregateConfig struct {
	InspecConfig InspecConfig `toml:"inspec_config"`
	InfraConfig  InfraConfig  `toml:"infra_config"`
}

type InspecConfig struct {
	State string `toml:"state"`
}

type InfraConfig struct {
	State string `toml:"state"`
}

func defaultConfig() aggregateConfig {
	return aggregateConfig{
		InspecConfig: InspecConfig{
			State: "success",
		},
		InfraConfig: InfraConfig{
			State: "success",
		},
	}
}

// NewManager - create a new config. There should only be one config for the service.
func NewManager(configFile string) (*Manager, error) {
	storedConfig, err := readConfigFromFile(configFile, defaultConfig())
	if err != nil {
		return &Manager{}, err
	}

	updateQueue := make(chan update, 100)
	manager := &Manager{
		config:      storedConfig,
		configFile:  configFile,
		updateQueue: updateQueue,
	}

	// Single goroutine that updates the Config data and saves the data to the config file.
	go manager.configUpdater(updateQueue)

	// Initial Testing
	err = manager.updateConfig(func(config aggregateConfig) (aggregateConfig, error) {
		return storedConfig, nil
	})
	if err != nil {
		return &Manager{}, err
	}

	return manager, nil
}

func (manager *Manager) GetInspecConfig() InspecConfig {
	return manager.config.InspecConfig
}

// UpdateInspecConfig - update the InSpec config
func (manager *Manager) UpdateInspecConfig(inspecConfig InspecConfig) error {
	return manager.updateConfig(func(config aggregateConfig) (aggregateConfig, error) {
		config.InspecConfig = inspecConfig
		return config, nil
	})
}

func (manager *Manager) GetInfraConfig() InfraConfig {
	return manager.config.InfraConfig
}

// UpdateInfraConfig - update the Infra config
func (manager *Manager) UpdateInfraConfig(infraConfig InfraConfig) error {
	return manager.updateConfig(func(config aggregateConfig) (aggregateConfig, error) {
		config.InfraConfig = infraConfig
		return config, nil
	})
}

// Close - to close out the channel for this object. This should only be called when the service is being shutdown
func (manager *Manager) Close() {
	// closes the updateQueue channel and ends that update goroutine
	close(manager.updateQueue)
}

// UpdateConfig - update the config
func (manager *Manager) updateConfig(updateFunc func(aggregateConfig) (aggregateConfig, error)) error {
	errc := make(chan error)
	manager.updateQueue <- update{fun: updateFunc, errc: errc}

	// Wait for the function to run
	err := <-errc
	close(errc)
	return err
}

func (manager *Manager) configUpdater(updateQueue <-chan update) {
	for update := range updateQueue {
		// Update the config object
		c, err := update.fun(manager.config)
		if err != nil {
			update.errc <- err
			continue
		}

		manager.config = c

		update.errc <- manager.saveToFile(c)
	}
}

func (manager *Manager) saveToFile(config aggregateConfig) error {
	log.WithFields(log.Fields{
		"config_file": manager.configFile,
	}).Debug("Saving Config File")

	tomlData, err := toml.Marshal(config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error Marshaling Config struct")
		return err
	}

	// create a file with the permissions of the running user can read/write other can only read.
	var permissions os.FileMode = 0644
	err = ioutil.WriteFile(manager.configFile, tomlData, permissions)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"config_file": manager.configFile,
		}).Error("Error writing config file")
		return err
	}

	return err
}

func readConfigFromFile(configFile string, defaultConfig aggregateConfig) (aggregateConfig, error) {
	config := defaultConfig

	tomlData, err := ioutil.ReadFile(configFile)
	if os.IsNotExist(err) {
		// config file does not exists use the default config
		return config, nil
	} else if err != nil {
		log.WithFields(log.Fields{
			"config_file": configFile,
		}).WithError(err).Error("Unable to read config file")

		return defaultConfig, err
	}

	err = toml.Unmarshal(tomlData, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"config_file": configFile,
		}).WithError(err).Error("Unable to load manager configuration")

		// Could not load data from config file using the default config.
		return defaultConfig, nil
	}

	return config, nil
}
