package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strings"
)

type FakeNUTMessage struct {
	Command string
	Options map[string]string
}

type FakeNUTResponse struct {
	Variable string
	Value    interface{}
}

// type FakeNUTCommandHandler func(device *FakeNUTDevice, options map[string]string) []FakeNUTResponse

// type FakeNUTCommand struct {
// 	Name    string
// 	Handler FakeNUTCommandHandler
// }

type FakeUPSCommand struct {
	Command string `json:"command"`
	Option  string `json:"option"`
	Arg     string `json:"arg"`
}

type FakeNUTServer struct {
	// Hostname or IP address to listen on. Defaults to "localhost".
	Host string

	// Port to listen on. Defaults to "3493".
	Port string

	// Map of devices, eg. "FakeUPS: FakeNUTDevice"
	Devices map[string]*FakeNUTDevice
	// Commands []FakeNUTCommand
}

// FakeNUTDevice represents a fake NUT device.
type FakeNUTDevice struct {
	// Battery charge in percent. Example: 100
	BatteryCharge int `json:"battery.charge"`

	// Battery charge low in percent. Example: 20
	BatteryChargeLow int `json:"battery.charge.low"`

	// Battery charge warning in percent. Example: 25
	BatteryChargeWarning int `json:"battery.charge.warning"`

	// Battery manufacturer date. Example: 1
	BatteryMfrDate string `json:"battery.mfr.date"`

	// Battery runtime in seconds. Example: 1620
	BatteryRuntime int `json:"battery.runtime"`

	// Battery runtime low in seconds. Example: 300
	BatteryRuntimeLow int `json:"battery.runtime.low"`

	// Battery type. Example: PbAcid
	BatteryType string `json:"battery.type"`

	// Battery voltage in volts. Example: 26
	BatteryVoltage int `json:"battery.voltage"`

	// Battery nominal voltage in volts. Example: 24
	BatteryVoltageNominal int `json:"battery.voltage.nominal"`

	// Device manufacturer. Example: 1
	DeviceMfr string `json:"device.mfr"`

	// Device model. Example: Powerwalker VI 2200 RLE
	DeviceModel string `json:"device.model"`

	// Device serial. Example: 000000000000
	DeviceSerial string `json:"device.serial"`

	// Device type. Example: ups
	DeviceType string `json:"device.type"`

	// Driver name. Example: usbhid-ups
	DriverName string `json:"driver.name"`

	// Driver parameter poll frequency in seconds. Example: 40
	DriverParameterPollFreq int `json:"driver.parameter.pollfreq"`

	// Driver parameter poll interval in seconds. Example: 2
	DriverParameterPollInterval int `json:"driver.parameter.pollinterval"`

	// Driver parameter port. Example: auto
	DriverParameterPort string `json:"driver.parameter.port"`

	// Driver parameter synchronous. Example: auto
	DriverParameterSynchronous string `json:"driver.parameter.synchronous"`

	// Driver version. Example: 2.8.0
	DriverVersion string `json:"driver.version"`

	// Driver version data. Example: CyberPower HID 0.6
	DriverVersionData string `json:"driver.version.data"`

	// Driver version internal. Example: 0.47
	DriverVersionInternal string `json:"driver.version.internal"`

	// Driver version USB. Example: libusb-1.0.0 (API: 0x1000102)
	DriverVersionUSB string `json:"driver.version.usb"`

	// Input frequency in hertz. Example: 50.0
	InputFrequency float64 `json:"input.frequency"`

	// Input transfer high in volts. Example: 290
	InputTransferHigh int `json:"input.transfer.high"`

	// Input transfer low in volts. Example: 165
	InputTransferLow int `json:"input.transfer.low"`

	// Input voltage in volts. Example: 232.6
	InputVoltage float64 `json:"input.voltage"`

	// Input nominal voltage in volts. Example: 230
	InputVoltageNominal int `json:"input.voltage.nominal"`

	// Output frequency in hertz. Example: 50.0
	OutputFrequency float64 `json:"output.frequency"`

	// Output voltage in volts. Example: 2.3
	OutputVoltage float64 `json:"output.voltage"`

	// UPS beeper status. Example: disabled
	UPSBeeperStatus string `json:"ups.beeper.status"`

	// UPS delay shutdown in seconds. Example: 20
	UPSDelayShutdown int `json:"ups.delay.shutdown"`

	// UPS delay start in seconds. Example: 30
	UPSDelayStart int `json:"ups.delay.start"`

	// UPS load in percent. Example: 12
	UPSLoad int `json:"ups.load"`

	// UPS manufacturer. Example: 1
	UPSMfr string `json:"ups.mfr"`

	// UPS model. Example: 2200R
	UPSModel string `json:"ups.model"`

	// UPS product ID. Example: 0601
	UPSProductID string `json:"ups.productid"`

	// UPS real power nominal in watts. Example: 1320
	UPSRealPowerNominal int `json:"ups.realpower.nominal"`

	// UPS serial. Example: 000000000000
	UPSSerial string `json:"ups.serial"`

	// UPS status. Example: OL
	UPSStatus string `json:"ups.status"`

	// UPS timer shutdown in seconds. Example: -60
	UPSTimerShutdown int `json:"ups.timer.shutdown"`

	// UPS timer start in seconds. Example: -60
	UPSTimerStart int `json:"ups.timer.start"`

	// UPS vendor ID. Example: 0764
	UPSVendorID string `json:"ups.vendorid"`
}

func NewFakeNUTServer() *FakeNUTServer {
	// Create a new fake NUT device.
	device := &FakeNUTDevice{
		BatteryCharge:               100,
		BatteryChargeLow:            20,
		BatteryChargeWarning:        25,
		BatteryMfrDate:              "1",
		BatteryRuntime:              1620,
		BatteryRuntimeLow:           300,
		BatteryType:                 "PbAcid",
		BatteryVoltage:              26,
		BatteryVoltageNominal:       24,
		DeviceMfr:                   "1",
		DeviceModel:                 "FakeNUT Server",
		DeviceSerial:                "000000000000",
		DeviceType:                  "ups",
		DriverName:                  "usbhid-ups",
		DriverParameterPollFreq:     40,
		DriverParameterPollInterval: 2,
		DriverParameterPort:         "auto",
		DriverParameterSynchronous:  "auto",
		DriverVersion:               "2.8.0",
		DriverVersionData:           "FakeNUT Server",
		DriverVersionInternal:       "0.47",
		DriverVersionUSB:            "libusb-1.0.0 (API: 0x1000102)",
		InputFrequency:              50.0,
		InputTransferHigh:           290,
		InputTransferLow:            165,
		InputVoltage:                232.6,
		InputVoltageNominal:         230,
		OutputFrequency:             50.0,
		OutputVoltage:               2.3,
		UPSBeeperStatus:             "disabled",
		UPSDelayShutdown:            20,
		UPSDelayStart:               30,
		UPSLoad:                     12,
		UPSMfr:                      "1",
		UPSModel:                    "2200R",
		UPSProductID:                "0601",
		UPSRealPowerNominal:         1320,
		UPSSerial:                   "000000000000",
		UPSStatus:                   "OL",
		UPSTimerShutdown:            -60,
		UPSTimerStart:               -60,
		UPSVendorID:                 "0764",
	}

	// Create a new fake NUT server.
	server := &FakeNUTServer{
		Host:    "localhost",
		Port:    "3493",
		Devices: map[string]*FakeNUTDevice{"FakeUPS": device},
	}

	// Set the server host and/or port if the appropriate environment variables are set.
	if host, ok := os.LookupEnv("NUT_SERVER"); ok {
		server.Host = host
	}
	if port, ok := os.LookupEnv("NUT_PORT"); ok {
		server.Port = port
	}

	// Return the server.
	return server
}

func (fakeNUTServer *FakeNUTServer) handleUPSCommand(conn net.Conn, command string) {
	// // Parse the command
	// parts := strings.Split(cmd, " ")
	// if len(parts) < 2 {
	// 	return fmt.Errorf("invalid command")
	// }
	// var upsCmd, upsVar string
	// upsCmd, upsVar = parts[0], parts[1]

	args := strings.Split(command, " ")
	cmd := args[0]
	subCmd := ""
	subCmdVal := ""
	subCmdVar := ""
	if len(args) > 1 {
		subCmd = args[1]
	}
	if len(args) > 2 {
		subCmdVal = args[2]
	}
	if len(args) > 3 {
		subCmdVar = args[3]
	}

	// log.Println("Received command:", cmd, subCmd, subCmdVal, subCmdVar)

	// The first argument is the command name.
	// The second argument is either a variable or a subcommand, and is optional.
	// The third argument is the value to set the variable to, and is optional.

	switch cmd {
	case "HELP":
		// Handle HELP command
		fmt.Fprint(conn, "Commands: ")
		fmt.Fprint(conn, "HELP ")
		fmt.Fprint(conn, "VER ")
		fmt.Fprint(conn, "GET ")
		fmt.Fprint(conn, "LIST ")
		fmt.Fprint(conn, "SET ")
		fmt.Fprint(conn, "INSTCMD ")
		fmt.Fprint(conn, "LOGIN ")
		fmt.Fprint(conn, "LOGOUT ")
		fmt.Fprint(conn, "USERNAME ")
		fmt.Fprint(conn, "PASSWORD ")
		fmt.Fprint(conn, "STARTTLS\n")
		log.Println("Sent HELP response")
	case "VER":
		// Handle VER command
		fmt.Fprintln(conn, "Fake UPS Server")
		// log.Println("Sent VER response")
	case "GET":
		// Handle GET command
		switch subCmd {
		case "NUMLOGINS":
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			fmt.Fprintf(conn, "NUMLOGINS %s 1\n", subCmdVal)
			// log.Println("Sent NUMLOGINS response")
		case "CMDDESC":
			// Handle GET CMDDESC command
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			if subCmdVar == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			fmt.Fprintf(conn, "CMDDESC %s %s \"Description unavailable\"\n", subCmdVal, subCmdVar)
			// log.Println("Sent CMDDESC response")
		case "UPSDESC":
			// Handle GET UPSDESC command
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			fmt.Fprintf(conn, "UPSDESC %s \"Fake UPS Device\"\n", subCmdVal)
			// log.Println("Sent UPSDESC response")
		default:
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			fmt.Fprintf(conn, "NUMLOGINS %s 1\n", subCmdVal)
			// log.Println("Sent NUMLOGINS response")
		}
	case "LIST":
		// TODO: Handle LIST command
		switch subCmd {
		case "UPS":
			// Handle LIST UPS command
			fmt.Fprintln(conn, "BEGIN LIST UPS")
			for upsName := range fakeNUTServer.Devices {
				fmt.Fprintf(conn, "UPS %s \"Description unavailable\"\n", upsName)
			}
			fmt.Fprintln(conn, "END LIST UPS")
			// log.Println("Sent LIST UPS response")
		case "CLIENT":
			// Handle LIST CLIENT command
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			// TODO: Ensure that this doesn't just use a hardcoded "UPS" but use the actual UPS name?!
			fmt.Fprintf(conn, "BEGIN LIST CLIENT %s\n", subCmdVal)
			fmt.Fprintf(conn, "CLIENT %s 127.0.0.1\n", subCmdVal)
			fmt.Fprintf(conn, "END LIST CLIENT %s\n", subCmdVal)
			// log.Println("Sent LIST CLIENT response")
		case "CMD":
			// Handle LIST CMD command
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			fmt.Fprintf(conn, "BEGIN LIST CMD %s\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s beeper.disable\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s beeper.enable\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s beeper.mute\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s beeper.off\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s beeper.on\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s load.off\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s load.off.delay\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s load.on\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s load.on.delay\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s shutdown.return\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s shutdown.stayoff\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s shutdown.stop\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s test.battery.start.deep\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s test.battery.start.quick\n", subCmdVal)
			fmt.Fprintf(conn, "CMD %s test.battery.stop\n", subCmdVal)
			fmt.Fprintf(conn, "END LIST CMD %s\n", subCmdVal)
			// log.Println("Sent LIST CMD response")
		case "VAR":
			// Handle LIST VAR command
			if subCmdVal == "" {
				fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
				// log.Println("Sent ERR response")
				return
			}
			fmt.Fprintf(conn, "BEGIN LIST VAR %s\n", subCmdVal)
			// log.Println("Sent LIST VAR header")

			// Get the device based on the UPS command name
			var fakeNUTDevice *FakeNUTDevice
			var deviceOk bool
			if fakeNUTDevice, deviceOk = fakeNUTServer.Devices[subCmdVal]; !deviceOk {
				fmt.Fprintf(conn, "ERR NOSUCHCMD %s\n", subCmdVal)
				// log.Println("Sent ERR response")
			}

			// Use reflection to get and print all field values from the fakeNUTDevice
			fakeNUTDeviceValue := reflect.ValueOf(fakeNUTDevice).Elem()
			for i := 0; i < fakeNUTDeviceValue.NumField(); i++ {
				field := fakeNUTDeviceValue.Type().Field(i)

				// Get the JSON tag value from the field
				jsonTag := field.Tag.Get("json")

				// Get the field value from the fakeNUTDevice
				fieldValue := fakeNUTDeviceValue.Field(i)

				fmt.Fprintf(conn, "VAR %s %s \"%s\"\n", subCmdVal, jsonTag, fieldValue)
				// log.Println("Sent VAR response: VAR", subCmdVal, jsonTag, fieldValue)
			}

			fmt.Fprintf(conn, "END LIST VAR %s\n", subCmdVal)
			// log.Println("Sent LIST VAR footer")

			// // Print all the variables for the device in the format required by NUT
			// for _, variable := range fakeNUTDevice.Variables {
			// 	fmt.Fprintf(conn, "VAR %s %s %s", subCmdVal, variable.Name, variable.Type)
			// }
			// fmt.Fprintln(conn, "END LIST VAR %s", subCmdVal)

			// // Use reflection to get the field value from the fakeNUTDevice
			// fieldValue := reflect.ValueOf(&fakeNUTDevice).Elem().FieldByName(upsCommand.Arg)
			// if !fieldValue.IsValid() {
			// 	fmt.Fprintf(conn, "ERR BADVAR %s\n", upsCommand.Arg)
			// 	return
			// }

			// // Compare the requested value to the actual value using reflection
			// requestedValue := reflect.ValueOf(upsCommand.Arg)
			// if requestedValue.Type() != fieldValue.Type() {
			// 	fmt.Fprintf(conn, "ERR BADTYPE %s\n", upsCommand.Arg)
			// 	return
			// }
			// if requestedValue.Interface() != fieldValue.Interface() {
			// 	fmt.Fprintf(conn, "ERR REJECTED %s\n", upsCommand.Arg)
			// 	return
			// }

			// // Return the value of the requested variable
			// fmt.Fprintf(conn, "OK %s %v\n", upsCommand.Arg, fieldValue.Interface())
		default:
			fmt.Fprintln(conn, "ERR INVALID-ARGUMENT")
			// log.Println("Sent ERR response")
		}
	case "SET":
		// TODO: Handle SET command
		fmt.Fprintln(conn, "ERR INVALID ARGUMENT")
		// log.Println("Sent ERR response")
	case "INSTCMD":
		// Handle INSTCMD command
		fmt.Fprintln(conn, "ERR USERNAME-REQUIRED")
		// log.Println("Sent ERR response")
	case "LOGIN":
		// Handle LOGIN command
		fmt.Fprintln(conn, "OK")
		// log.Println("Sent OK response")
	case "LOGOUT":
		// Handle LOGOUT command
		fmt.Fprintln(conn, "OK")
		// log.Println("Sent OK response")
	case "USERNAME":
		// Handle USERNAME command
		fmt.Fprintln(conn, "OK")
		// log.Println("Sent OK response")
	case "PASSWORD":
		// Handle PASSWORD command
		fmt.Fprintln(conn, "OK")
		// log.Println("Sent OK response")
	case "STARTTLS":
		// Handle STARTTLS command
		fmt.Fprintln(conn, "ERR FEATURE-NOT-CONFIGURED")
		// log.Println("STARTTLS not supported")
	default:
		fmt.Fprintln(conn, "ERR UNKNOWN-COMMAND")
		// log.Println("Unknown command:", command)
	}

	// // Parse the command into a UPSCommand struct
	// upsCommand := FakeUPSCommand{}
	// err := json.Unmarshal([]byte(command), &upsCommand)
	// if err != nil {
	// 	fmt.Fprintf(conn, "ERR BADCMD %s\n", err)
	// }

	// // Get the device based on the UPS command name
	// var fakeNUTDevice *FakeNUTDevice
	// var deviceOk bool
	// if fakeNUTDevice, deviceOk = fakeNUTServer.Devices[upsCommand.Command]; !deviceOk {
	// 	fmt.Fprintf(conn, "ERR NOSUCHCMD %s\n", upsCommand.Command)
	// }

	// // Use reflection to get the field value from the fakeNUTDevice
	// fieldValue := reflect.ValueOf(&fakeNUTDevice).Elem().FieldByName(upsCommand.Arg)
	// if !fieldValue.IsValid() {
	// 	fmt.Fprintf(conn, "ERR BADVAR %s\n", upsCommand.Arg)
	// 	return
	// }

	// // Compare the requested value to the actual value using reflection
	// requestedValue := reflect.ValueOf(upsCommand.Arg)
	// if requestedValue.Type() != fieldValue.Type() {
	// 	fmt.Fprintf(conn, "ERR BADTYPE %s\n", upsCommand.Arg)
	// 	return
	// }
	// if requestedValue.Interface() != fieldValue.Interface() {
	// 	fmt.Fprintf(conn, "ERR REJECTED %s\n", upsCommand.Arg)
	// 	return
	// }

	// // Return the value of the requested variable
	// fmt.Fprintf(conn, "OK %s %v\n", upsCommand.Arg, fieldValue.Interface())

	// // Get the corresponding field in the struct
	// var field reflect.Value
	// var ok bool
	// if field, ok = reflect.ValueOf(&fakeNUTDevice).Elem().Type().FieldByNameFunc(func(fieldName string) bool {
	// 	return strings.ToLower(fieldName) == strings.ToLower(upsVar)
	// }); !ok {
	// 	return fmt.Errorf("Fake NUT server invalid variable %s", upsVar)
	// }

	// // Get the value of the field
	// val := reflect.ValueOf(fakeNUTDevice).FieldByName(field.Name)

	// // Format the response using the JSON tags in the struct
	// res := fmt.Sprintf("%s %s\n", upsCmd, field.Tag.Get("json")+": "+val.String())

	// // Send the response
	// _, err := conn.Write([]byte(res))
	// return err
}

func (fakeNUTServer *FakeNUTServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", fakeNUTServer.Host, fakeNUTServer.Port))
	if err != nil {
		return fmt.Errorf("Fake NUT server failed to start server: %w", err)
	}
	defer listener.Close()

	log.Printf("Fake NUT server listening on %s:%s", fakeNUTServer.Host, fakeNUTServer.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Fake NUT server error accepting connection: %v", err)
			continue
		}

		// log.Printf("Fake NUT server accepted connection from %s", conn.RemoteAddr())

		go func() {
			defer conn.Close()

			reader := bufio.NewReader(conn)

			for {
				command, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("Fake NUT server error reading from connection: %v", err)
					}
					return
				}

				command = strings.TrimSpace(command)

				if config.Verbose {
					log.Printf("Fake NUT server received command from %s: %s", conn.RemoteAddr(), command)
				}

				fakeNUTServer.handleUPSCommand(conn, command)

				// if strings.HasPrefix(command, "UPS") {
				// 	fakeNUTServer.handleUPSCommand(conn, command)
				// } else {
				// 	response := "ERR Unknown command\n"
				// 	_, err := conn.Write([]byte(response))
				// 	if err != nil {
				// 		log.Printf("Error sending response: %v", err)
				// 	}
				// }
			}
		}()
	}
}

// FIXME: Implement this!
func (fakeNUTServer *FakeNUTServer) Stop() error {
	log.Println("NOTICE: Fake NUT server graceful shutdown not implemented yet!")
	return nil
}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"strings"
// )

// // NUTCommand represents a NUT command.
// type NUTCommand struct {
// 	Name   string
// 	Args   []string
// 	String string
// }

// // NUTResponse represents a NUT response.
// type NUTResponse string

// const (
// 	// NUTResponseOK is the OK response.
// 	NUTResponseOK NUTResponse = "OK"

// 	// NUTResponseVersion is the version response.
// 	NUTResponseVersion NUTResponse = "OK 1.7"

// 	// NUTResponseUnknownCmd is the unknown command response.
// 	NUTResponseUnknownCmd NUTResponse = "ERR UNKNOWNCMD"

// 	// NUTResponseInvalidArg is the invalid argument response.
// 	NUTResponseInvalidArg NUTResponse = "ERR INVALIDARG"
// )

// // NUTResponseData represents a NUT response with data.
// type NUTResponseData struct {
// 	Response NUTResponse
// 	Data     string
// }

// // Stringer for NUTResponse which adds a new line to the end.
// func (r NUTResponse) String() string {
// 	return string(r) + "\n"
// }

// // ParseNUTCommand parses a NUT command from a string.
// func ParseNUTCommand(s string) (*NUTCommand, error) {
// 	// Split the command into parts.
// 	parts := strings.Split(strings.TrimSpace(s), " ")

// 	// Make sure we have at least one part.
// 	if len(parts) < 1 {
// 		return nil, fmt.Errorf("Invalid command: %q", s)
// 	}

// 	// Return the parsed NUT command.
// 	return &NUTCommand{
// 		Name:   strings.ToUpper(parts[0]),
// 		Args:   parts[1:],
// 		String: s,
// 	}, nil
// }

// // FakeNUTServer represents a fake NUT server.
// type FakeNUTServer struct{}

// // FakeNUTDevice represents a fake NUT device.
// type FakeNUTDevice struct {
// 	// Battery charge in percent. Example: 100
// 	BatteryCharge int `json:"battery.charge"`

// 	// Battery charge low in percent. Example: 20
// 	BatteryChargeLow int `json:"battery.charge.low"`

// 	// Battery charge warning in percent. Example: 25
// 	BatteryChargeWarning int `json:"battery.charge.warning"`

// 	// Battery manufacturer date. Example: 1
// 	BatteryMfrDate string `json:"battery.mfr.date"`

// 	// Battery runtime in seconds. Example: 1620
// 	BatteryRuntime int `json:"battery.runtime"`

// 	// Battery runtime low in seconds. Example: 300
// 	BatteryRuntimeLow int `json:"battery.runtime.low"`

// 	// Battery type. Example: PbAcid
// 	BatteryType string `json:"battery.type"`

// 	// Battery voltage in volts. Example: 26
// 	BatteryVoltage int `json:"battery.voltage"`

// 	// Battery nominal voltage in volts. Example: 24
// 	BatteryVoltageNominal int `json:"battery.voltage.nominal"`

// 	// Device manufacturer. Example: 1
// 	DeviceMfr string `json:"device.mfr"`

// 	// Device model. Example: Powerwalker VI 2200 RLE
// 	DeviceModel string `json:"device.model"`

// 	// Device serial. Example: 000000000000
// 	DeviceSerial string `json:"device.serial"`

// 	// Device type. Example: ups
// 	DeviceType string `json:"device.type"`

// 	// Driver name. Example: usbhid-ups
// 	DriverName string `json:"driver.name"`

// 	// Driver parameter poll frequency in seconds. Example: 40
// 	DriverParameterPollFreq int `json:"driver.parameter.pollfreq"`

// 	// Driver parameter poll interval in seconds. Example: 2
// 	DriverParameterPollInterval int `json:"driver.parameter.pollinterval"`

// 	// Driver parameter port. Example: auto
// 	DriverParameterPort string `json:"driver.parameter.port"`

// 	// Driver parameter synchronous. Example: auto
// 	DriverParameterSynchronous string `json:"driver.parameter.synchronous"`

// 	// Driver version. Example: 2.8.0
// 	DriverVersion string `json:"driver.version"`

// 	// Driver version data. Example: CyberPower HID 0.6
// 	DriverVersionData string `json:"driver.version.data"`

// 	// Driver version internal. Example: 0.47
// 	DriverVersionInternal string `json:"driver.version.internal"`

// 	// Driver version USB. Example: libusb-1.0.0 (API: 0x1000102)
// 	DriverVersionUSB string `json:"driver.version.usb"`

// 	// Input frequency in hertz. Example: 50.0
// 	InputFrequency float64 `json:"input.frequency"`

// 	// Input transfer high in volts. Example: 290
// 	InputTransferHigh int `json:"input.transfer.high"`

// 	// Input transfer low in volts. Example: 165
// 	InputTransferLow int `json:"input.transfer.low"`

// 	// Input voltage in volts. Example: 232.6
// 	InputVoltage float64 `json:"input.voltage"`

// 	// Input nominal voltage in volts. Example: 230
// 	InputVoltageNominal int `json:"input.voltage.nominal"`

// 	// Output frequency in hertz. Example: 50.0
// 	OutputFrequency float64 `json:"output.frequency"`

// 	// Output voltage in volts. Example: 2.3
// 	OutputVoltage float64 `json:"output.voltage"`

// 	// UPS beeper status. Example: disabled
// 	UPSBeeperStatus string `json:"ups.beeper.status"`

// 	// UPS delay shutdown in seconds. Example: 20
// 	UPSDelayShutdown int `json:"ups.delay.shutdown"`

// 	// UPS delay start in seconds. Example: 30
// 	UPSDelayStart int `json:"ups.delay.start"`

// 	// UPS load in percent. Example: 12
// 	UPSLoad int `json:"ups.load"`

// 	// UPS manufacturer. Example: 1
// 	UPSMfr string `json:"ups.mfr"`

// 	// UPS model. Example: 2200R
// 	UPSModel string `json:"ups.model"`

// 	// UPS product ID. Example: 0601
// 	UPSProductID string `json:"ups.productid"`

// 	// UPS real power nominal in watts. Example: 1320
// 	UPSRealPowerNominal int `json:"ups.realpower.nominal"`

// 	// UPS serial. Example: 000000000000
// 	UPSSerial string `json:"ups.serial"`

// 	// UPS status. Example: OL
// 	UPSStatus string `json:"ups.status"`

// 	// UPS timer shutdown in seconds. Example: -60
// 	UPSTimerShutdown int `json:"ups.timer.shutdown"`

// 	// UPS timer start in seconds. Example: -60
// 	UPSTimerStart int `json:"ups.timer.start"`

// 	// UPS vendor ID. Example: 0764
// 	UPSVendorID string `json:"ups.vendorid"`
// }

// // HandleCommand handles a NUT command and returns a response.
// func (s *FakeNUTServer) HandleCommand(cmd *NUTCommand) string {
// 	// Handle the NUT command.
// 	switch cmd.Name {
// 	case "GET":
// 		return s.handleGetCommand(cmd)
// 	case "SET":
// 		return s.handleSetCommand(cmd)
// 	case "USERNAME":
// 		return "OK\n"
// 	case "PASSWORD":
// 		return "OK\n"
// 	case "LIST":
// 		return s.handleListCommand(cmd)
// 	case "UPS":
// 		return s.handleUPSCommand(cmd)
// 	default:
// 		return NUTResponseUnknownCmd.String()
// 	}
// }

// // handleGetCommand handles a GET command.
// func (s *FakeNUTServer) handleGetCommand(cmd *NUTCommand) string {
// 	return "OK battery.charge: 90\n"
// }

// // handleSetCommand handles a SET command.
// func (s *FakeNUTServer) handleSetCommand(cmd *NUTCommand) string {
// 	return "OK\n"
// }

// // handleListCommand handles a LIST command.
// func (s *FakeNUTServer) handleListCommand(cmd *NUTCommand) string {
// 	if len(cmd.Args) != 1 || cmd.Args[0] != "VAR" {
// 		return "ERR INVALIDARG\n"
// 	}
// 	return "VAR battery.charge 0 100 \"%\"\n"
// }

// // handleUPSCommand handles a UPS command.
// func (s *FakeNUTServer) handleUPSCommand(cmd *NUTCommand) string {
// 	if len(cmd.Args) != 0 {
// 		return "ERR INVALIDARG\n"
// 	}
// 	return fmt.Sprintf("UPS %s\n", "fakeups")
// }

// func fakeNutStart() {
// 	listener, err := net.Listen("tcp", "127.0.0.1:3493")
// 	if err != nil {
// 		fmt.Println("Fake NUT -> Error listening:", err.Error())
// 		return
// 	}
// 	defer listener.Close()

// 	if config.Verbose {
// 		fmt.Println("Fake NUT -> Server started and listening on localhost:3493")
// 	}

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Fake NUT -> Error accepting:", err.Error())
// 			return
// 		}

// 		if config.Verbose {
// 			fmt.Println("Fake NUT -> Client connected.")
// 		}

// 		go fakeNutHandleConnection(conn)
// 	}
// }

// func fakeNutHandleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)

// 	for {
// 		// Read the data sent from the client
// 		data, err := reader.ReadString('\n')
// 		if err != nil {
// 			fmt.Println("Fake NUT -> Error reading:", err.Error())
// 			return
// 		}

// 		// Trim any whitespace from the received data
// 		data = strings.TrimSpace(data)

// 		// Print the received data if verbose mode is enabled
// 		if config.Verbose {
// 			fmt.Println("Received data:", data)
// 		}

// 		// Default to an error response
// 		response := errResponse

// 		// Parse the NUT command and generate the appropriate response
// 		if strings.HasPrefix(data, "SET ") {
// 			// Respond with a fake OK for SET commands
// 			response = "OK\n"
// 		} else if strings.HasPrefix(data, "USERNAME ") {
// 			// Respond with a fake OK for username commands
// 			response = "OK\n"
// 		} else if strings.HasPrefix(data, "PASSWORD ") {
// 			// Respond with a fake OK for password commands
// 			response = "OK\n"
// 		} else if strings.HasPrefix(data, "VER") {
// 			// Respond with a fake OK for version commands
// 			response = "OK 1.7\n"
// 		} else if strings.HasPrefix(data, "NETVER") {
// 			// Respond with a fake OK for version commands
// 			response = "OK 1.7\n"
// 		} else if strings.HasPrefix(data, "LOGOUT") {
// 			// Respond with a fake OK for logout commands
// 			response = "OK\n"
// 		} else if strings.HasPrefix(data, "LIST UPS") {
// 			// Respond with a fake UPS for list ups commands
// 			response = "BEGIN LIST UPS\nUPS UPS \"Description unavailable\"\nEND LIST UPS\n"
// 		} else if strings.HasPrefix(data, "LIST CLIENT ") {
// 			// Respond with a fake client for list client commands
// 			response = "BEGIN LIST CLIENT UPS\nCLIENT UPS 127.0.0.1\nEND LIST CLIENT UPS\n"
// 		} else if strings.HasPrefix(data, "UPS ") {
// 			// Generate fake data for UPS <name> commands
// 			name := strings.TrimPrefix(data, "UPS ")
// 			response = fmt.Sprintf("UPS %s\n", name)
// 		} else if strings.HasPrefix(data, "HELP") {
// 			// Respond with a fake OK for help commands
// 			response = "OK\n"
// 		} else if strings.HasPrefix(data, "LIST CMD ") {
// 			// Responsd with a fake list of commands for list cmd commands
// 			response = "BEGIN LIST CMD UPS\nCMD UPS beeper.disable\nCMD UPS beeper.enable\nCMD UPS beeper.mute\nCMD UPS beeper.off\nCMD UPS beeper.on\nCMD UPS load.off\nCMD UPS load.off.delay\nCMD UPS load.on\nCMD UPS load.on.delay\nCMD UPS shutdown.return\nCMD UPS shutdown.stayoff\nCMD UPS shutdown.stop\nCMD UPS test.battery.start.deep\nCMD UPS test.battery.start.quick\nCMD UPS test.battery.stop\nEND LIST CMD UPS\n"
// 		} else if strings.HasPrefix(data, "GET CMDDESC ") {
// 			// Get the device name as the first parameter
// 			// and the command name as the second parameter
// 			params := strings.Split(strings.TrimPrefix(data, "GET CMDDESC "), " ")
// 			if len(params) == 2 {
// 				// Respond with a fake description for the command
// 				response = fmt.Sprintf("CMDDESC %s %s \"Description unavailable\"\n", params[0], params[1])
// 			}
// 		} else if strings.HasPrefix(data, "GET UPSDESC ") {
// 			// Get the device name as the first parameter
// 			params := strings.Split(strings.TrimPrefix(data, "GET UPSDESC "), " ")
// 			if len(params) == 1 {
// 				// Respond with a fake description for the device
// 				response = fmt.Sprintf("UPSDESC %s \"Description unavailable\"\n", params[0])
// 			}
// 		} else if strings.HasPrefix(data, "GET NUMLOGINS ") {
// 			// Get the device name as the first parameter
// 			params := strings.Split(strings.TrimPrefix(data, "GET NUMLOGINS "), " ")
// 			if len(params) == 1 {
// 				// Respond with a fake number of logins for the device
// 				response = fmt.Sprintf("NUMLOGINS %s 1\n", params[0])
// 			}
// 		} else if strings.HasPrefix(data, "LIST VAR ") {
// 			// Get the device name as the first parameter
// 			params := strings.Split(strings.TrimPrefix(data, "LIST VAR "), " ")
// 			if len(params) == 1 {
// 				// Generate fake variables
// 				vars := "VAR UPS battery.charge \"100\"\nVAR UPS battery.charge.low \"20\"\nVAR UPS battery.charge.warning \"25\"\nVAR UPS battery.mfr.date \"1 \"\nVAR UPS battery.runtime \"1620\"\nVAR UPS battery.runtime.low \"300\"\nVAR UPS battery.type \"PbAcid\"\nVAR UPS battery.voltage \"26\"\nVAR UPS battery.voltage.nominal \"24\"\nVAR UPS device.mfr \"1 \"\nVAR UPS device.model \"Powerwalker VI 2200 RLE\"\nVAR UPS device.serial \"000000000000\"\nVAR UPS device.type \"ups\"\nVAR UPS driver.name \"usbhid-ups\"\nVAR UPS driver.parameter.pollfreq \"40\"\nVAR UPS driver.parameter.pollinterval \"2\"\nVAR UPS driver.parameter.port \"auto\"\nVAR UPS driver.parameter.synchronous \"auto\"\nVAR UPS driver.version \"2.8.0\"\nVAR UPS driver.version.data \"CyberPower HID 0.6\"\nVAR UPS driver.version.internal \"0.47\"\nVAR UPS driver.version.usb \"libusb-1.0.0 (API: 0x1000102)\"\nVAR UPS input.frequency \"50.0\"\nVAR UPS input.transfer.high \"290\"\nVAR UPS input.transfer.low \"165\"\nVAR UPS input.voltage \"232.6\"\nVAR UPS input.voltage.nominal \"230\"\nVAR UPS output.frequency \"50.0\"\nVAR UPS output.voltage \"2.3\"\nVAR UPS ups.beeper.status \"disabled\"\nVAR UPS ups.delay.shutdown \"20\"\nVAR UPS ups.delay.start \"30\"\nVAR UPS ups.load \"12\"\nVAR UPS ups.mfr \"1 \"\nVAR UPS ups.model \"2200R\"\nVAR UPS ups.productid \"0601\"\nVAR UPS ups.realpower.nominal \"1320\"\nVAR UPS ups.serial \"000000000000\"\nVAR UPS ups.status \"OL\"\nVAR UPS ups.timer.shutdown \"-60\"\nVAR UPS ups.timer.start \"-60\"\nVAR UPS ups.vendorid \"0764\"\n"

// 				// Respond with the fake variables
// 				response = fmt.Sprintf("BEGIN LIST VAR %s\n%sEND LIST VAR %s\n", params[0], vars, params[0])
// 			}
// 		} else if strings.HasPrefix(data, "GET DESC ") {
// 			// Get the device name as the first parameter
// 			// and the variable name as the second parameter
// 			params := strings.Split(strings.TrimPrefix(data, "GET DESC "), " ")
// 			if len(params) == 2 {
// 				// Respond with a fake description for the variable
// 				response = fmt.Sprintf("DESC %s %s \"Description unavailable\"\n", params[0], params[1])
// 			}
// 		} else if strings.HasPrefix(data, "GET TYPE ") {
// 			// Get the device name as the first parameter
// 			// and the variable name as the second parameter
// 			params := strings.Split(strings.TrimPrefix(data, "GET TYPE "), " ")
// 			if len(params) == 2 {
// 				// Respond with a fake type for the variable
// 				response = fmt.Sprintf("TYPE %s %s NUMBER\n", params[0], params[1])
// 			}
// 		}

// 		// Send the response back to the client
// 		_, err = conn.Write([]byte(response))
// 		if err != nil {
// 			fmt.Println("Fake NUT -> Error writing:", err.Error())
// 			return
// 		}

// 		if config.Verbose {
// 			fmt.Println("Fake NUT -> Sent response:", response)
// 		}
// 	}
// }
