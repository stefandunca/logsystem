package logsystem

import (
	"encoding/json"
	"fmt"
	"os"
)

const FileDriverID = "file"

type fileConfig struct {
	UserReadableTime bool   `json:"userReadableTime"`
	FilePath         string `json:"filePath"`
}

// FileDriverFactory implements DriverFactoryInterface
type FileDriverFactory struct {
}

func (f *FileDriverFactory) DriverID() DriverID {
	return DriverID(FileDriverID)
}

func (f *FileDriverFactory) CreateDriver(config json.RawMessage) (DriverInterface, error) {
	var fileConfig fileConfig
	err := json.Unmarshal(config, &fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file driver: %w", err)
	}

	file, err := openFile(fileConfig.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s; error: %w", fileConfig.FilePath, err)
	}

	return &FileDriver{
		file:   file,
		config: fileConfig,
	}, nil
}

// FileDriver implements DriverInterface
type FileDriver struct {
	config fileConfig
	file   *os.File
}

func (d *FileDriver) Log(data map[Param]string) {
	line := formatLine(data, d.config.UserReadableTime)
	d.file.WriteString(line + "\n")
}

func (d *FileDriver) BeginTx(id TxID, attr map[Param]string) {
	txData := make(map[Param]string)
	txData[TxIDParam] = id.String()
	message := fmt.Sprintf("TX Begin; Params: %v", attr)
	txData[MessageParam] = message
	txData[LevelParam] = string(Info)
	d.Log(txData)
}

func (d *FileDriver) EndTx(id TxID) {
	txData := make(map[Param]string)
	txData[TxIDParam] = id.String()
	txData[MessageParam] = "TX End"
	txData[LevelParam] = string(Info)
	d.Log(txData)
}

func (d *FileDriver) Stop() {
	if d.file != nil {
		d.file.Close()
	}
}

func openFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %s; %v\n", filePath, err)
		return nil, err
	}
	return file, nil
}
