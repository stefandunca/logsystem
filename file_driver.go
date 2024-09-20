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

func (f *FileDriverFactory) driverID() DriverID {
	return DriverID(FileDriverID)
}

func (f *FileDriverFactory) createDriver(config json.RawMessage) DriverInterface {
	var fileConfig fileConfig
	err := json.Unmarshal(config, &fileConfig)
	if err != nil {
		fmt.Printf("Failed to unmarshal file driver config: %v\n", err)
		return nil
	}

	file, err := openFile(fileConfig.FilePath)
	if err != nil {
		return nil
	}

	return &FileDriver{
		file:   file,
		config: fileConfig,
	}
}

// FileDriver implements DriverInterface
type FileDriver struct {
	config fileConfig
	file   *os.File
}

func (d *FileDriver) log(data map[Param]string) {
	line := formatLine(data, d.config.UserReadableTime)
	d.file.WriteString(line + "\n")
}

func (d *FileDriver) beginTx(id TxID, attr map[Param]string) {
	txData := make(map[Param]string)
	txData[TxIDParam] = id.String()
	message := fmt.Sprintf("TX Begin; Params: %v", attr)
	txData[MessageParam] = message
	txData[LevelParam] = string(Info)
	d.log(txData)
}

func (d *FileDriver) endTx(id TxID) {
	txData := make(map[Param]string)
	txData[TxIDParam] = id.String()
	txData[MessageParam] = "TX End"
	txData[LevelParam] = string(Info)
	d.log(txData)
}

func (d *FileDriver) stop() {
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
