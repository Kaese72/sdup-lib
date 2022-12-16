package sdupcache

import (
	"errors"

	"github.com/Kaese72/sdup-lib/devicestoretemplates"
	"github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sdupcache/filters"
	"github.com/Kaese72/sdup-lib/sduptemplates"
)

type SDUPCacheVolImpl struct {
	devices    map[string]sduptemplates.DeviceSpec
	groups     map[sduptemplates.DeviceGroupID]sduptemplates.DeviceGroupSpec
	target     sduptemplates.SDUPTarget
	updateChan chan sduptemplates.Update
}

func NewSDUPCache(target sduptemplates.SDUPTarget) SDUPCache {
	return &SDUPCacheVolImpl{
		devices:    map[string]sduptemplates.DeviceSpec{},
		groups:     map[sduptemplates.DeviceGroupID]sduptemplates.DeviceGroupSpec{},
		target:     target,
		updateChan: make(chan sduptemplates.Update, 10),
	}
}

func (cache *SDUPCacheVolImpl) Initialize() (chan sduptemplates.Update, error) {
	upstreamUpdates, err := cache.target.Initialize()
	if err != nil {
		return nil, err
	}

	go func() {
		for update := range upstreamUpdates {
			if dUpdate, err := update.GetDeviceUpdate(); err == nil {
				relevantUpdate := cache.updateDevice(dUpdate)
				if relevantUpdate.Relevant() {
					logging.Info("Device update considered relevant", map[string]string{"DeviceID": string(relevantUpdate.ID)})
					// Update is relevant, and should be passed on
					cache.updateChan <- sduptemplates.UpdateFromDeviceUpdate(relevantUpdate)

				}

			} else if gUpdate, err := update.GetDeviceGroupUpdate(); err == nil {
				if cache.updateGroup(gUpdate) {
					// Update is relevant, and should be passed on
					logging.Info("Group update", map[string]string{"DeviceID": string(gUpdate.GroupID)})
					cache.updateChan <- update

				}

			} else {
				panic("Could not identify update type")
			}
		}
	}()

	return cache.updateChan, nil
}

func (cache SDUPCacheVolImpl) Devices() ([]sduptemplates.DeviceSpec, error) {
	specs := []sduptemplates.DeviceSpec{}
	for _, spec := range cache.devices {
		specs = append(specs, spec)
	}

	return specs, nil
}

func (cache SDUPCacheVolImpl) Groups() ([]sduptemplates.DeviceGroupSpec, error) {
	specs := []sduptemplates.DeviceGroupSpec{}
	for _, spec := range cache.groups {
		specs = append(specs, spec)
	}

	return specs, nil
}

func (cache SDUPCacheVolImpl) Device(deviceID string) (sduptemplates.DeviceSpec, error) {
	if item, ok := cache.devices[deviceID]; ok {
		return item, nil
	}
	return sduptemplates.DeviceSpec{}, sduptemplates.NoSuchDevice
}

func (cache SDUPCacheVolImpl) TriggerCapability(deviceID string, capKey devicestoretemplates.CapabilityKey, capArg devicestoretemplates.CapabilityArgs) error {
	return cache.target.TriggerCapability(deviceID, capKey, capArg)
}

func (cache SDUPCacheVolImpl) GTriggerCapability(dgid sduptemplates.DeviceGroupID, ck devicestoretemplates.CapabilityKey, ca devicestoretemplates.CapabilityArgs) error {
	return cache.target.GTriggerCapability(dgid, ck, ca)
}

func (cache *SDUPCacheVolImpl) updateDevice(update sduptemplates.DeviceUpdate) sduptemplates.DeviceUpdate {
	var relevantUpdate sduptemplates.DeviceUpdate
	if device, ok := cache.devices[update.ID]; ok {
		cache.devices[update.ID], relevantUpdate = device.ApplyUpdate(update)

	} else {
		cache.devices[update.ID] = update.UpdateToDevice()
		relevantUpdate = update
	}
	return relevantUpdate
}

func (cache *SDUPCacheVolImpl) updateGroup(update sduptemplates.DeviceGroupUpdate) bool {
	newGroup := update.UpdateToDeviceGroup()
	if group, ok := cache.groups[update.GroupID]; ok {
		if !newGroup.Equal(group) {
			// Group is an update
			cache.groups[update.GroupID] = newGroup
			return true

		} else {
			// Group is not an update
			return false
		}

	} else {
		// Group is new, and should be passed along
		cache.groups[update.GroupID] = newGroup
		return true
	}
}

func deviceMatchesFilters(device sduptemplates.DeviceSpec, filters filters.AttributeFilters) (match bool, err error) {
	for _, filter := range filters {
		operator, err := filter.GetOperator()
		if err != nil {
			// Invalid operators lead to wacky scenarios
			return false, err
		}

		if _, _, err := filter.Key.KeyValKeys(); err == nil {
			// Composite key, we should use keyval
			return false, errors.New("keyval currently not supported")

		} else {
			if _, ok := device.Attributes[devicestoretemplates.AttributeKey(filter.Key)]; !ok {
				// Not having the attribute counts as false
				return false, nil
			}
			// Simple key
			// Get value based on what type the comparator is
			switch comp := filter.Value.(type) {
			case int:
				return matchNumericComparison(device.Attributes[devicestoretemplates.AttributeKey(filter.Key)].Numeric, float32(comp), operator)

			case float32:
				return matchNumericComparison(device.Attributes[devicestoretemplates.AttributeKey(filter.Key)].Numeric, comp, operator)

			case string:
				return matchStringComparison(device.Attributes[devicestoretemplates.AttributeKey(filter.Key)].Text, comp, operator)

			case bool:
				return matchBooleanComparison(device.Attributes[devicestoretemplates.AttributeKey(filter.Key)].Boolean, comp, operator)

			default:
				// FIXME log better
				return false, errors.New("unsupported filter type")
			}
		}

	}
	return true, nil
}

func matchBooleanComparison(attrVal *bool, compVal bool, operator filters.Operator) (bool, error) {
	if attrVal == nil {
		// Not having the value set is considered
		return false, nil
	}

	switch operator {
	case filters.Equal:
		return *attrVal == compVal, nil
	default:
		return false, errors.New("not a supported operand")
	}
}

func matchStringComparison(attrVal *string, compVal string, operator filters.Operator) (bool, error) {
	if attrVal == nil {
		// Not having the value set is considered
		return false, nil
	}

	switch operator {
	case filters.Equal:
		return *attrVal == compVal, nil
	default:
		return false, errors.New("not a supported operand")
	}
}

func matchNumericComparison(attrVal *float32, compVal float32, operator filters.Operator) (bool, error) {
	if attrVal == nil {
		// Not having the value set is considered
		return false, nil
	}

	switch operator {
	case filters.Equal:
		return *attrVal == compVal, nil
	default:
		return false, errors.New("not a supported operand")
	}
}

func (store *SDUPCacheVolImpl) FilteredDevices(attrFilters filters.AttributeFilters) ([]sduptemplates.DeviceSpec, error) {
	specs := []sduptemplates.DeviceSpec{}
	for _, device := range store.devices {
		match, err := deviceMatchesFilters(device, attrFilters)
		if err != nil {
			return nil, err
		}

		if match {
			specs = append(specs, device)
		}
	}
	return specs, nil
}
