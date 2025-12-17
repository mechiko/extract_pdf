package licenser

import (
	"golang.org/x/sys/windows/registry"
)

// Example
// root - registry.CURRENT_USER
// keyPath - `Software\MyApp`
// Check if a registry key exists.
func keyExists(root registry.Key, keyPath string) bool {
	k, err := registry.OpenKey(root, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	k.Close()
	return true
}

// createKey creates a new registry key or opens an existing one.
func createKey(root registry.Key, keyPath string) (registry.Key, error) {
	k, _, err := registry.CreateKey(root, keyPath, registry.ALL_ACCESS)
	return k, err
}

// readStringValueWithDefault reads a string value from the Windows Registry with a default value.
func readStringValueWithDefault(root registry.Key, keyPath, valueName, defaultValue string) (string, error) {
	k, err := registry.OpenKey(root, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return defaultValue, nil // Return the default value if the key or value doesn't exist
	}
	defer k.Close()

	value, _, err := k.GetStringValue(valueName)
	if err != nil {
		return defaultValue, nil // Return the default value if the value doesn't exist
	}

	return value, nil
}

// writeStringValue writes an string value to the Windows Registry.
func writeStringValue(root registry.Key, keyPath, valueName, data string) error {
	k, err := registry.OpenKey(root, keyPath, registry.WRITE)
	if err != nil {
		return err
	}
	defer k.Close()

	if err := k.SetStringValue(valueName, data); err != nil {
		return err
	}
	return nil
}

// Check if a registry value exists.
func valueExists(root registry.Key, keyPath, valueName string) bool {
	k, err := registry.OpenKey(root, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(valueName)
	return err == nil
}
