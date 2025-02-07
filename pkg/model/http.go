package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kdudkov/goatak/pkg/cot"
	"github.com/kdudkov/goatak/pkg/cotproto"
)

type WebUnit struct {
	Uid            string    `json:"uid"`
	Callsign       string    `json:"callsign"`
	Category       string    `json:"category"`
	Team           string    `json:"team"`
	Role           string    `json:"role"`
	Time           time.Time `json:"time"`
	LastSeen       time.Time `json:"last_seen"`
	StaleTime      time.Time `json:"stale_time"`
	StartTime      time.Time `json:"start_time"`
	SendTime       time.Time `json:"send_time"`
	Type           string    `json:"type"`
	Lat            float64   `json:"lat"`
	Lon            float64   `json:"lon"`
	Hae            float64   `json:"hae"`
	Speed          float64   `json:"speed"`
	Course         float64   `json:"course"`
	Sidc           string    `json:"sidc"`
	TakVersion     string    `json:"tak_version"`
	Status         string    `json:"status"`
	Text           string    `json:"text"`
	Color          string    `json:"color"`
	Icon           string    `json:"icon"`
	ParentCallsign string    `json:"parent_callsign"`
	ParentUid      string    `json:"parent_uid"`
	Local          bool      `json:"local"`
	Send           bool      `json:"send"`
}

type Contact struct {
	Uid          string `json:"uid"`
	Callsign     string `json:"callsign"`
	Team         string `json:"team"`
	Role         string `json:"role"`
	Takv         string `json:"takv"`
	Notes        string `json:"notes"`
	FilterGroups string `json:"filterGroups"`
}

type DigitalPointer struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Name string  `json:"name"`
}

func (i *Item) ToWeb() *WebUnit {
	evt := i.msg.TakMessage.CotEvent

	i.mx.RLock()
	defer i.mx.RUnlock()

	w := &WebUnit{
		Uid:            i.uid,
		Category:       i.class,
		Callsign:       i.callsign,
		Time:           cot.TimeFromMillis(evt.SendTime),
		LastSeen:       i.lastSeen,
		StaleTime:      i.staleTime,
		StartTime:      i.startTime,
		SendTime:       i.sendTime,
		Type:           i.cottype,
		Lat:            evt.Lat,
		Lon:            evt.Lon,
		Hae:            evt.Hae,
		Speed:          evt.GetDetail().GetTrack().GetSpeed(),
		Course:         evt.GetDetail().GetTrack().GetCourse(),
		Team:           evt.GetDetail().GetGroup().GetName(),
		Role:           evt.GetDetail().GetGroup().GetRole(),
		Sidc:           getSIDC(i.cottype),
		ParentUid:      i.parentUid,
		ParentCallsign: i.parentCallsign,
		Color:          fmt.Sprintf("#%.6x", i.color&0xffffff),
		Icon:           i.icon,
		Local:          i.local,
		Send:           i.send,
	}

	if i.class == CONTACT {
		if i.online {
			w.Status = "Online"
		} else {
			w.Status = "Offline"
		}

		if v := i.msg.TakMessage.GetCotEvent().GetDetail().GetTakv(); v != nil {
			w.TakVersion = strings.Trim(fmt.Sprintf("%s %s on %s", v.Platform, v.Version, v.Device), " ")
		}
	}

	w.Text = i.msg.Detail.GetFirst("remarks").GetText()
	return w
}

func (w *WebUnit) ToMsg() *cot.CotMessage {
	msg := &cotproto.TakMessage{
		CotEvent: &cotproto.CotEvent{
			Type:      w.Type,
			Uid:       w.Uid,
			SendTime:  cot.TimeToMillis(w.SendTime),
			StartTime: cot.TimeToMillis(w.StartTime),
			StaleTime: cot.TimeToMillis(w.StaleTime),
			How:       "h-g-i-g-o",
			Lat:       w.Lat,
			Lon:       w.Lon,
			Hae:       w.Hae,
			Ce:        9999999,
			Le:        9999999,
			Detail: &cotproto.Detail{
				Contact: &cotproto.Contact{Callsign: w.Callsign},
				PrecisionLocation: &cotproto.PrecisionLocation{
					Geopointsrc: "DTED0",
					Altsrc:      "USER",
				},
			},
		},
	}

	xd := cot.NewXmlDetails()
	if w.ParentUid != "" {
		m := map[string]string{"uid": w.ParentUid, "relation": "p-p"}
		if w.ParentCallsign != "" {
			m["parent_callsign"] = w.ParentCallsign
		}
		xd.AddChild("link", m, "")
	}
	xd.AddChild("status", map[string]string{"readiness": "true"}, "")
	if w.Text != "" {
		xd.AddChild("remarks", nil, w.Text)
	}
	msg.GetCotEvent().Detail.XmlDetail = xd.AsXMLString()

	zero := time.Unix(0, 0)
	if msg.CotEvent.Uid == "" {
		msg.CotEvent.Uid = uuid.New().String()
	}

	if w.StartTime.Before(zero) {
		msg.CotEvent.StartTime = cot.TimeToMillis(time.Now())
	}

	if w.SendTime.Before(zero) {
		msg.CotEvent.SendTime = cot.TimeToMillis(time.Now())
	}

	if w.StaleTime.Before(zero) {
		msg.CotEvent.StaleTime = cot.TimeToMillis(time.Now().Add(time.Hour * 24))
	}

	return &cot.CotMessage{
		From:       "",
		TakMessage: msg,
		Detail:     xd,
	}
}

func getSIDC(fn string) string {
	if !strings.HasPrefix(fn, "a-") {
		return ""
	}

	tokens := strings.Split(fn, "-")

	sidc := "S" + tokens[1]

	if len(tokens) > 2 {
		sidc += tokens[2] + "P"
	} else {
		sidc += "-P"
	}

	if len(tokens) > 3 {
		for _, c := range tokens[3:] {
			if len(c) > 1 {
				break
			}
			sidc += c
		}
	}

	if len(sidc) < 12 {
		sidc += strings.Repeat("-", 10-len(sidc))
	}
	return strings.ToUpper(sidc)
}
