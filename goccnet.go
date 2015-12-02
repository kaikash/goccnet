package goccnet
import(
  "github.com/tarm/serial"
  "github.com/joaojeronimo/go-crc16"
  "bytes"
  "log"
)

const sync = 0x02

type Device struct {
  isConnect bool
  busy bool
  // command *Command
  config *DeviceConfig
  serialPort *serial.Port
}
type DeviceConfig struct {
  DeviceType byte
  Path string
  Baud int
}

func NewDevice(config *DeviceConfig) *Device {
  // conf := &serial.Config{Name: path, Baud: baud}
  // s, err := serial.OpenPort(conf)
  // if err != nil {
    // log.Fatal(err)
  // }
  device := &Device{false, false, config, nil}
  return device
}

func (device *Device) Connect() {
  conf := &serial.Config{Name: device.config.Path, Baud: device.config.Baud}
  s, err := serial.OpenPort(conf)
  if err != nil {
    log.Fatal(err)
  }
  device.serialPort = s
  err = device.Reset()
  if err != nil {
    log.Fatal(err)
  }
}

func (device *Device) Reset() error {
  var code byte = 0x30
  _, err:=device.Execute(code, nil)
  return err
}

func (device *Device) Execute(code byte, data []byte) ([]byte, error) {
  cmd:=bytes.NewBuffer([]byte{sync, device.config.DeviceType})
  code_arr := []byte{code}
  cmd.Write([]byte{(byte)(len(code_arr)+5)})
  cmd.Write([]byte{code})
  // cmd.Write([]byte{cmd, crc16.Crc16(cmd)})
  res := bytes.NewBuffer(cmd.Bytes())
  crc := crc16.Crc16(cmd.Bytes())
  res.Write([]byte{(byte)(crc & 0xff), (byte)((crc>>8) & 0xff) })
  log.Println(res)
  err := device.serialPort.Write(res.Bytes())
  return nil, nil
}
