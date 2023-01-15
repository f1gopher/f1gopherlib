// F1GopherLib - Copyright (C) 2022 f1gopher
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connection

import (
	"context"
	"fmt"
	"github.com/f1gopher/f1gopherlib/f1log"
	"github.com/f1gopher/signalr/v2"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
	"time"
)

type live struct {
	log     *f1log.F1GopherLibLog
	archive *os.File
	ctx     context.Context
	wg      *sync.WaitGroup
	c2      *signalr.Conn
	client  *signalr.Client

	dataFeed chan Payload
}

func CreateLive(ctx context.Context, wg *sync.WaitGroup, log *f1log.F1GopherLibLog) *live {
	return &live{
		ctx:      ctx,
		wg:       wg,
		log:      log,
		dataFeed: make(chan Payload, 1000),
		archive:  nil,
	}
}

func CreateArchivingLive(ctx context.Context, archiveFile string) (*live, error) {
	archive, err := os.Create(fmt.Sprintf("%s_%d.txt", archiveFile, time.Now().UnixMilli()))
	if err != nil {
		return nil, err
	}

	return &live{
		ctx:      ctx,
		dataFeed: make(chan Payload, 1000),
		archive:  archive,
	}, nil
}

func (l *live) Connect() (error, <-chan Payload) {
	var err error

	// Prepare a SignalR client.
	l.c2, err = signalr.Dial(
		l.ctx,
		"https://livetiming.formula1.com/signalr",
		`[{"name":"streaming"}]`,
	)
	if err != nil {
		l.log.Errorf("Connect to live failed: %v", err)
		return err, nil
	}

	l.client = signalr.NewClient("streaming", l.c2)

	errg, ctx := errgroup.WithContext(l.ctx)
	errg.Go(func() error { return l.client.Run(ctx) })
	errg.Go(func() error {
		l.wg.Add(1)
		defer l.wg.Done()

		stream, err1 := l.client.Callback(ctx, "feed")
		if err1 != nil {
			return err1
		}
		defer stream.Close()

		l.log.Info("Waiting for live data...")

		for {
			res := stream.ReadRaw()

			if res.Args != nil {
				if len(res.Args) == 3 {
					data := Payload{}
					abc, _ := res.Args[0].MarshalJSON()
					data.Name = string(abc[1 : len(abc)-1])
					abc, _ = res.Args[1].MarshalJSON()
					if abc[0] == '"' {
						data.Data = abc[1 : len(abc)-1]
					} else {
						data.Data = abc
					}
					abc, _ = res.Args[2].MarshalJSON()
					data.Timestamp = string(abc[1 : len(abc)-1])

					if l.archive != nil {
						l.archive.WriteString(data.Name + "\r\n")
						l.archive.Write(data.Data)
						l.archive.WriteString("\r\n" + data.Timestamp + "\r\n")
					}

					l.dataFeed <- data
				} else if len(res.Args) == 1 {
					data := Payload{
						Name:      CatchupFile,
						Timestamp: "",
					}
					abc, _ := res.Args[0].MarshalJSON()
					data.Data = abc

					if l.archive != nil {
						l.archive.WriteString(data.Name + "\r\n")
						l.archive.Write(data.Data)
						l.archive.WriteString("\r\n" + data.Timestamp + "\r\n")
					}

					l.dataFeed <- data
				} else {
					l.log.Errorf("There is an unhandled number of arguments for live data: %d, dropping data", len(res.Args))
				}
			}
		}
	})
	err = l.client.Invoke(ctx, "Subscribe", OrderedFiles).Exec()
	if err != nil {
		l.log.Errorf("Live connection subscribe failed: %v", err)
		return err, nil
	}

	l.log.Info("Connected to live")

	return nil, l.dataFeed
}

// Can't do anything because this is live data
func (l *live) IncrementTime(amount time.Duration) {}
