package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "slave",
				Aliases: []string{"s"},
				Value:   100,
				Usage:   "Slave ID",
			},
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Value:    "",
				Usage:    "/dev/usb....",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "turnon",
				Aliases: []string{"to"},
				Usage:   "Turn on",
				Action:  turnOn,
			},
			{
				Name:    "turnoff",
				Aliases: []string{"tf"},
				Usage:   "Turn off",
				Action:  turnOff,
			},
			{
				Name:    "powerdata",
				Aliases: []string{"pd"},
				Usage:   "Read power data",
				Action:  readData,
			},
			{
				Name:    "product",
				Aliases: []string{"po"},
				Usage:   "Read product information",
				Action:  readProduct,
			},
			{
				Name:    "operation",
				Aliases: []string{"op"},
				Usage:   "Operation parameter",
				Action:  readOpParameters,
			},
			{
				Name:    "runstatus",
				Aliases: []string{"rs"},
				Usage:   "Operation parameter",
				Action:  readRunStatus,
			},
			{
				Name:    "protectparameters",
				Aliases: []string{"pp"},
				Usage:   "Protect parameters",
				Action:  readProtectParameters,
			},
			{
				Name:    "setrecordnum",
				Aliases: []string{"sr"},
				Usage:   "setRecordNumber --recordtype=type, Record type --recordnum Record No.",
				Action:  setRecordNumber,
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "recordtype", Usage: "--recordtype"},
					&cli.IntFlag{Name: "recordnum", Usage: "--recordnum"},
				},
			},
			{
				Name:    "readrecord",
				Aliases: []string{"rr"},
				Usage:   "record data",
				Action:  readRecord,
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "recordtype", Usage: "--recordtype"},
				},
			},
			{
				Name:    "timerparameters",
				Aliases: []string{"tp"},
				Usage:   "Timer parameters",
				Action:  readTimerParameters,
			},
			{
				Name:    "settimerparameters",
				Aliases: []string{"st"},
				Usage:   "Json file path --jsonpath=file",
				Action:  setTimerParameters,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "jsonpath", Usage: "--jsonpath"},
				},
			},
			{
				Name:    "telecommand",
				Aliases: []string{"tc"},
				Usage:   "telecommand",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "summary1",
				Aliases: []string{"su1"},
				Usage:   "Summary data",
				Action:  readSummary1,
			},
			{
				Name:    "summary2",
				Aliases: []string{"su2"},
				Usage:   "Summary data",
				Action:  readSummary2,
			},
			{
				Name:    "summary3",
				Aliases: []string{"su3"},
				Usage:   "Summary data",
				Action:  readSummary3,
			},
			{
				Name:    "summary4",
				Aliases: []string{"su4"},
				Usage:   "Summary data",
				Action:  readSummary4,
			},
			{
				Name:    "logs",
				Aliases: []string{"lg"},
				Usage:   "readLogs --logtype=type --logindex=index",
				Action:  readLogs,
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "logtype", Usage: "--logtype"},
					&cli.IntFlag{Name: "logindex", Usage: "--logindex"},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		Logger.Fatal(err)
	}
}

/*
func main() {

	set := flag.NewFlagSet("flag", 0)
	set.String("address", "com3", "test")
	set.Int("slave", 100, "test")

	ctx := cli.NewContext(nil, set, nil)

	//readProduct(ctx)
	//readOpParameters(ctx)
	//readRunStatus(ctx)
	//readData(ctx)
	//readProtectParameters(ctx)
	//readSummary1(ctx)
	//readSummary2(ctx)
	//readSummary3(ctx)
	//readSummary4(ctx)
	setTimerToSwitch(ctx)
	readRemoteCmd(ctx)
}
*/
func turnOff(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}
	return SwitchBreaker(client, false)
}

func turnOn(c *cli.Context) error {
	Logger.Println("Turning on")
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}
	return SwitchBreaker(client, true)
}

func setTimerParameters(c *cli.Context) error {
	Logger.Println("set Timers to control switch")
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	jsonpath := c.String("jsonpath")

	return SetTimerParameters(client, jsonpath)
}

func readData(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadData(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readLogs(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	var addr uint16 = 0
	logtype := uint16(c.Int("logtype"))
	logindex := uint16(c.Int("logindex"))

	if logtype == FAULT_TYPE && logindex < MAX_FAULTRECORDLOG_NUM {
		addr = (FAULTRECORDLOG_ADDR + logindex*RECORD_LOG_LEN)

	} else if logtype == ALARM_TYPE && logindex < MAX_ALARMRECORDLOG_NUM {
		addr = (ALARMRECORDLOG_ADDR + logindex*RECORD_LOG_LEN)
	} else if logtype == SWITCH_TYPE && logindex < MAX_SWITCHRECORDLOG_NUM {
		addr = (SWITCHRECORDLOG_ADDR + logindex*RECORD_LOG_LEN)
	} else {
		return err
	}
	data, err := ReadLogs(client, addr)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readProtectParameters(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadProtectParameters(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func setRecordNumber(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	recordtype := c.Int("recordtype")
	recordnum := uint16(c.Int("recordnum"))
	if recordtype == FAULT_TYPE {
		return SetRecordNo(client, FAULTRECORD_NUM_ADDR, recordnum)
	} else if recordtype == ALARM_TYPE {
		return SetRecordNo(client, ALARMRECORD_NUM_ADDR, recordnum)
	} else if recordtype == SWITCH_TYPE {
		return SetRecordNo(client, SWITCHRECORD_NUM_ADDR, recordnum)
	}

	return err
}

func readRecord(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	recordtype := c.Int("recordtype")
	var addr uint16 = 0

	if recordtype == FAULT_TYPE {
		addr = FAULTRECORD_ADDR
	} else if recordtype == ALARM_TYPE {
		addr = ALARMRECORD_ADDR
	} else if recordtype == SWITCH_TYPE {
		addr = SWITCHRECORD_ADDR
	} else {
		return err
	}

	data, err := ReadRecord(client, addr)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readSummary1(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadSummary1(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readSummary2(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadSummary2(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readSummary3(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadSummary3(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readSummary4(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadSummary4(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readProduct(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadProduct(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readOpParameters(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadOpParameters(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readRunStatus(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadRunStatus(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func readTimerParameters(c *cli.Context) error {
	client, err := openConnection(c)
	defer client.CloseConnection()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	data, err := ReadTimerParameters(client)
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	return outputData(data)
}

func outputData(data JsonMarshal) error {
	jsonData, err := data.ToJson()
	if err != nil {
		Logger.Fatal(err)
		return err
	}

	Logger.Println(string(jsonData))
	return nil
}

func openConnection(c *cli.Context) (*ModbusClient, error) {

	slaveID := c.Int("slave")
	address := c.String("address")

	Logger.Println(fmt.Sprintf("Opening connection to address:%s,slave:%d", address, slaveID))

	return ConnectSlave(address, uint8(slaveID))

}
