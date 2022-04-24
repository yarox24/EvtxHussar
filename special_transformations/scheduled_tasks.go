package special_transformations

import (
	"encoding/xml"
	"fmt"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/evtxhussar/common"
	"golang.org/x/net/html/charset"
	"io"
	"strconv"
	"strings"
)

type XMLTask struct {
	XMLName          xml.Name `xml:"Task"`
	Version          string   `xml:"version,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	RegistrationInfo struct {
		Author             string `xml:"Author"`
		Description        string `xml:"Description"`
		SecurityDescriptor string `xml:"SecurityDescriptor"`
		URI                string `xml:"URI"`
		Version            string `xml:"Version"`
		Source             string `xml:"Source"`
		Date               string `xml:"Date"`
		Documentation      string `xml:"Documentation"`
	} `xml:"RegistrationInfo"`
	Principals struct {
		Principal struct {
			ID                  string `xml:"id,attr"`
			GroupId             string `xml:"GroupId"`
			UserId              string `xml:"UserId"`
			RunLevel            string `xml:"RunLevel"`
			DisplayName         string `xml:"DisplayName"`
			LogonType           string `xml:"LogonType"`
			ProcessTokenSidType string `xml:"ProcessTokenSidType"`
		} `xml:"Principal"`
	} `xml:"Principals"`
	Settings struct {
		DisallowStartIfOnBatteries      bool   `xml:"DisallowStartIfOnBatteries"`
		StopIfGoingOnBatteries          bool   `xml:"StopIfGoingOnBatteries"`
		Enabled                         bool   `xml:"Enabled"`
		MultipleInstancesPolicy         string `xml:"MultipleInstancesPolicy"`
		StartWhenAvailable              string `xml:"StartWhenAvailable"`
		AllowHardTerminate              bool   `xml:"AllowHardTerminate"`
		RunOnlyIfNetworkAvailable       bool   `xml:"RunOnlyIfNetworkAvailable"`
		AllowStartOnDemand              bool   `xml:"AllowStartOnDemand"`
		Hidden                          bool   `xml:"Hidden"`
		RunOnlyIfIdle                   bool   `xml:"RunOnlyIfIdle"`
		DisallowStartOnRemoteAppSession bool   `xml:"DisallowStartOnRemoteAppSession"`
		UseUnifiedSchedulingEngine      bool   `xml:"UseUnifiedSchedulingEngine"`
		WakeToRun                       bool   `xml:"WakeToRun"`
		ExecutionTimeLimit              string `xml:"ExecutionTimeLimit"`
		DeleteExpiredTaskAfter          string `xml:"DeleteExpiredTaskAfter"`
		Priority                        string `xml:"Priority"`
		NetworkProfileName              string `xml:"NetworkProfileName"`
		IdleSettings                    struct {
			Duration      string `xml:"Duration"`
			WaitTimeout   string `xml:"WaitTimeout"`
			StopOnIdleEnd string `xml:"StopOnIdleEnd"`
			RestartOnIdle string `xml:"RestartOnIdle"`
		} `xml:"IdleSettings"`
		RestartOnFailure struct {
			Interval string `xml:"Interval"`
			Count    string `xml:"Count"`
		} `xml:"RestartOnFailure"`
	} `xml:"Settings"`
	Triggers struct {
		LogonTrigger []struct {
			ID            string `xml:"id,attr"`
			StartBoundary string `xml:"StartBoundary"`
			EndBoundary   string `xml:"EndBoundary"`
			Delay         string `xml:"Delay"`
			Enabled       string `xml:"Enabled"`
			Repetition    struct {
				Interval string `xml:"Interval"`
			} `xml:"Repetition"`
		} `xml:"LogonTrigger"`
		CalendarTrigger []struct {
			ID            string `xml:"id,attr"`
			StartBoundary string `xml:"StartBoundary"`
			Repetition    struct {
				Interval string `xml:"Interval"`
				Duration string `xml:"Duration"`
			} `xml:"Repetition"`
			ScheduleByDay struct {
				DaysInterval string `xml:"DaysInterval"`
			} `xml:"ScheduleByDay"`
		} `xml:"CalendarTrigger"`
		EventTrigger []struct {
			Enabled            string `xml:"Enabled"`
			ExecutionTimeLimit string `xml:"ExecutionTimeLimit"`
			Delay              string `xml:"Delay"`
			Repetition         struct {
				Interval string `xml:"Interval"`
				Duration string `xml:"Duration"`
			} `xml:"Repetition"`
			Subscription string `xml:"Subscription"`
		} `xml:"EventTrigger"`
		TimeTrigger []struct {
			ID            string `xml:"id,attr"`
			StartBoundary string `xml:"StartBoundary"`
			EndBoundary   string `xml:"EndBoundary"`
			Enabled       bool   `xml:"Enabled"`
		} `xml:"TimeTrigger"`
		BootTrigger []struct {
			Enabled string `xml:"Enabled"`
			Delay   string `xml:"Delay"`
		} `xml:"BootTrigger"`
		RegistrationTrigger []struct {
			Delay string `xml:"Delay"`
		} `xml:"RegistrationTrigger"`
		IdleTrigger []struct {
			ID         string `xml:"id,attr"`
			Repetition struct {
				Interval string `xml:"Interval"`
			} `xml:"Repetition"`
		} `xml:"IdleTrigger"`
		SessionStateChangeTrigger []struct {
			StateChange string `xml:"StateChange"`
		} `xml:"SessionStateChangeTrigger"`
	} `xml:"Triggers"`
	Actions struct {
		Context string `xml:"Context,attr"`
		Exec    []struct {
			Command          string `xml:"Command"`
			Arguments        string `xml:"Arguments"`
			WorkingDirectory string `xml:"WorkingDirectory"`
		} `xml:"Exec"`
		ComHandler []struct {
			ClassId string `xml:"ClassId"`
			Data    string `xml:"Data"`
		} `xml:"ComHandler"`
		SendEmail []struct {
			Server  string `xml:"Server"`
			Subject string `xml:"Subject"`
			To      string `xml:"To"`
			Cc      string `xml:"Cc"`
			Bcc     string `xml:"Bcc"`
			ReplyTo string `xml:"ReplyTo"`
			From    string `xml:"From"`
			//HeaderFields string `xml:"HeaderFields"`
			Body string `xml:"Body"`
			//Attachments      string `xml:"Attachments"`
		} `xml:"SendEmail"`
		ShowMessage []struct {
			Title string `xml:"Title"`
			Body  string `xml:"Body"`
		} `xml:"ShowMessage"`
	} `xml:"Actions"`
}

func BypassReader(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func DecodeUtf16XML(r io.Reader, v interface{}) (err error) {
	nr, err := charset.NewReader(r, "utf-16")
	if err != nil {
		return
	}
	decoder := xml.NewDecoder(nr)
	decoder.CharsetReader = BypassReader
	err = decoder.Decode(v)
	return
}

func XMLScheduledTask(ord_map *ordereddict.Dict, options map[string]string) {
	Input_field := options["input_field"]
	Output_field := options["output_field"]
	Path := options["path"]

	if !common.KeyExistsInOrderedDict(ord_map, Input_field) {
		panic("Wrong Yaml - field_extra_transformations - input_field")
	}

	if !common.KeyExistsInOrderedDict(ord_map, Output_field) {
		panic("Wrong Yaml - field_extra_transformations - output_field")
	}

	input_val, _ := ord_map.GetString(Input_field)

	if len(input_val) > 0 {

		r := strings.NewReader(input_val)
		task := new(XMLTask)

		err := DecodeUtf16XML(r, task)

		// Unexpected EOF is normal when XML task is truncated
		if err != nil && !strings.HasSuffix(err.Error(), "unexpected EOF") {
			common.LogErrorWithError("When decoding XML Scheduled Task: ", err)
		}

		// Error in XML
		if task.Xmlns != "http://schemas.microsoft.com/windows/2004/02/mit/task" {
			return
		}

		// XML parsed properly
		xml_element := ""

		// Triggers
		count_logontrigger := len(task.Triggers.LogonTrigger)
		count_boottrigger := len(task.Triggers.BootTrigger)
		count_idletrigger := len(task.Triggers.IdleTrigger)
		count_eventtrigger := len(task.Triggers.EventTrigger)
		count_timetrigger := len(task.Triggers.TimeTrigger)
		count_calendartrigger := len(task.Triggers.CalendarTrigger)
		count_registrationtrigger := len(task.Triggers.RegistrationTrigger)
		count_sessionstatechangetrigger := len(task.Triggers.SessionStateChangeTrigger)

		// Actions
		count_showmessage := len(task.Actions.ShowMessage)
		count_exec := len(task.Actions.Exec)
		count_comhandler := len(task.Actions.ComHandler)
		count_sendemail := len(task.Actions.SendEmail)

		switch Path {
		case "author":
			xml_element = task.RegistrationInfo.Author
		case "task_version":
			xml_element = task.Version
		case "description":
			xml_element = task.RegistrationInfo.Description
		case "uri":
			xml_element = task.RegistrationInfo.URI
		case "version":
			xml_element = task.RegistrationInfo.Version
		case "source":
			xml_element = task.RegistrationInfo.Source
		case "date":
			xml_element = task.RegistrationInfo.Date
		case "principal_id":
			xml_element = task.Principals.Principal.ID
		case "principal_userid":
			xml_element = task.Principals.Principal.UserId
		case "principal_groupid":
			xml_element = task.Principals.Principal.GroupId
		case "principal_runlevel":
			xml_element = task.Principals.Principal.RunLevel
		case "principal_displayname":
			xml_element = task.Principals.Principal.DisplayName
		case "principal_logontype":
			xml_element = task.Principals.Principal.LogonType
		case "task_enabled":
			xml_element = strconv.FormatBool(task.Settings.Enabled)
		case "task_hidden":
			xml_element = strconv.FormatBool(task.Settings.Hidden)
		case "triggers_summary":
			if count_logontrigger > 0 {
				xml_element += fmt.Sprintf("LogonTrigger: %d ,", count_logontrigger)
			}

			if count_boottrigger > 0 {
				xml_element += fmt.Sprintf("BootTrigger: %d ,", count_boottrigger)
			}

			if count_idletrigger > 0 {
				xml_element += fmt.Sprintf("IdleTrigger: %d ,", count_idletrigger)
			}

			if count_eventtrigger > 0 {
				xml_element += fmt.Sprintf("EventTrigger: %d ,", count_eventtrigger)
			}

			if count_timetrigger > 0 {
				xml_element += fmt.Sprintf("TimeTrigger: %d ,", count_timetrigger)
			}

			if count_calendartrigger > 0 {
				xml_element += fmt.Sprintf("CalendarTrigger: %d ,", count_calendartrigger)
			}

			if count_registrationtrigger > 0 {
				xml_element += fmt.Sprintf("RegistrationTrigger: %d ,", count_registrationtrigger)
			}

			if count_sessionstatechangetrigger > 0 {
				xml_element += fmt.Sprintf("SessionStateChangeTrigger: %d ,", count_sessionstatechangetrigger)
			}

			xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))

		case "calendartrigger_startboundary":
			if count_calendartrigger > 0 {
				if count_calendartrigger == 1 {
					xml_element = task.Triggers.CalendarTrigger[0].StartBoundary
				} else {
					for nr, ct := range task.Triggers.CalendarTrigger {
						if len(ct.StartBoundary) > 0 {
							xml_element += fmt.Sprintf("StartBoundary[%d]: %s ,", nr, ct.StartBoundary)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}

			}
		case "timetrigger_startboundary":
			if count_timetrigger > 0 {
				if count_timetrigger == 1 {
					xml_element = task.Triggers.TimeTrigger[0].StartBoundary
				} else {
					for nr, tt := range task.Triggers.TimeTrigger {
						if len(tt.StartBoundary) > 0 {
							xml_element += fmt.Sprintf("StartBoundary[%d]: %s ,", nr, tt.StartBoundary)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "actions_summary":
			if count_showmessage > 0 {
				xml_element += fmt.Sprintf("ShowMessage: %d ,", count_showmessage)
			}

			if count_exec > 0 {
				xml_element += fmt.Sprintf("Exec: %d ,", count_exec)
			}

			if count_comhandler > 0 {
				xml_element += fmt.Sprintf("ComHandler: %d ,", count_comhandler)
			}

			if count_sendemail > 0 {
				xml_element += fmt.Sprintf("SendEmail: %d ,", count_sendemail)
			}

			xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
		case "actions_context":
			xml_element = task.Actions.Context
		case "exec_command_with_arguments":
			if count_exec > 0 {
				if count_exec == 1 {
					action0 := task.Actions.Exec[0]
					xml_element = fmt.Sprintf("%s %s", action0.Command, action0.Arguments)
				} else {
					for nr, ae := range task.Actions.Exec {
						if len(ae.Command) > 0 {
							xml_element += fmt.Sprintf("Command[%d]: %s %s,", nr, ae.Command, ae.Arguments)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "exec_workingdirectory":
			if count_exec > 0 {
				if count_exec == 1 {
					action0 := task.Actions.Exec[0]
					xml_element = action0.WorkingDirectory
				} else {
					for nr, ae := range task.Actions.Exec {
						if len(ae.WorkingDirectory) > 0 {
							xml_element += fmt.Sprintf("WorkDir[%d]: %s ,", nr, ae.WorkingDirectory)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "comhandler_classid":
			if count_comhandler > 0 {
				if count_comhandler == 1 {
					xml_element = task.Actions.ComHandler[0].ClassId
				} else {
					for nr, ac := range task.Actions.ComHandler {
						if len(ac.ClassId) > 0 {
							xml_element += fmt.Sprintf("ClassId[%d]: %s ,", nr, ac.ClassId)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "comhandler_data":
			if count_comhandler > 0 {
				if count_comhandler == 1 {
					xml_element = task.Actions.ComHandler[0].Data
				} else {
					for nr, ac := range task.Actions.ComHandler {
						if len(ac.Data) > 0 {
							xml_element += fmt.Sprintf("Data[%d]: %s ,", nr, ac.Data)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "sendemail_server":
			if count_sendemail > 0 {
				if count_sendemail == 1 {
					xml_element = task.Actions.SendEmail[0].Server
				} else {
					for nr, se := range task.Actions.SendEmail {
						if len(se.Server) > 0 {
							xml_element += fmt.Sprintf("Server[%d]: %s ,", nr, se.Server)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "sendemail_to":
			if count_sendemail > 0 {
				if count_sendemail == 1 {
					xml_element = task.Actions.SendEmail[0].To
				} else {
					for nr, se := range task.Actions.SendEmail {
						if len(se.Server) > 0 {
							xml_element += fmt.Sprintf("To[%d]: %s ,", nr, se.To)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "sendemail_from":
			if count_sendemail > 0 {
				if count_sendemail == 1 {
					xml_element = task.Actions.SendEmail[0].From
				} else {
					for nr, se := range task.Actions.SendEmail {
						if len(se.Server) > 0 {
							xml_element += fmt.Sprintf("From[%d]: %s ,", nr, se.From)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "sendemail_body":
			if count_sendemail > 0 {
				if count_sendemail == 1 {
					xml_element = task.Actions.SendEmail[0].Body
				} else {
					for nr, se := range task.Actions.SendEmail {
						if len(se.Server) > 0 {
							xml_element += fmt.Sprintf("Body[%d]: %s ,", nr, se.Body)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "showmessage_title":
			if count_showmessage > 0 {
				if count_showmessage == 1 {
					xml_element = task.Actions.ShowMessage[0].Title
				} else {
					for nr, sm := range task.Actions.ShowMessage {
						if len(sm.Title) > 0 {
							xml_element += fmt.Sprintf("Title[%d]: %s ,", nr, sm.Title)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		case "showmessage_body":
			if count_showmessage > 0 {
				if count_showmessage == 1 {
					xml_element = task.Actions.ShowMessage[0].Body
				} else {
					for nr, sm := range task.Actions.ShowMessage {
						if len(sm.Title) > 0 {
							xml_element += fmt.Sprintf("Body[%d]: %s ,", nr, sm.Body)
						}
					}
					xml_element = strings.TrimSpace(strings.TrimSuffix(xml_element, ","))
				}
			}
		default:
			panic("Wrong path name")
		}

		ord_map.Update(Output_field, xml_element)
	}

}
