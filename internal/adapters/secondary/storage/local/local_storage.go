package local_storage

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"time"
)

const (
	LastIpAddressFile = "last_known_ip_address"
)

type LocalStorage struct {
	Directory string
}

func New(directory string) *LocalStorage {
	return &LocalStorage{Directory: directory}
}

func (l *LocalStorage) GetLastKnownIPAddress(ctx context.Context) (net.IP, error) {
	filename := l.Directory + "/" + LastIpAddressFile + ".json"

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, nil
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data LastKnownIPAddressData
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return nil, err
	}

	return data.IPAddress, nil
}

func (l *LocalStorage) SaveIPAddress(ctx context.Context, ip net.IP) error {
	filename := l.Directory + "/" + LastIpAddressFile + ".json"

	data := LastKnownIPAddressData{
		IPAddress: ip,
		CheckedAt: time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
