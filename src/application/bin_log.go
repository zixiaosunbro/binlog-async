package application

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	log "github.com/sirupsen/logrus"

	"binlog-async/src"
)

func InitBinlogSvc() {
	canalIns, err := canal.NewCanal(src.CanalCfg)
	if err != nil {
		panic(fmt.Sprintf("create canal ins fail:%s", err.Error()))
	}
	// register event handler with canal instance
	canalIns.SetEventHandler(&eventHandler{})
	log.Info("start binlog receiver")
	if err = canalIns.Run(); err != nil {
		panic(fmt.Sprintf("start canal with failure:%s", err.Error()))
	}
	// wait for SIGINT or SIGTERM signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	canalIns.Close()
	log.Println("binlog service Canal stopped")
}

// Define a custom event handler to process binlog events
type eventHandler struct{}

func (h *eventHandler) String() string {
	//TODO implement me
	panic("implement me")
}

func (h *eventHandler) OnRotate(header *replication.EventHeader, r *replication.RotateEvent) error {
	// Do nothing
	return nil
}

func (h *eventHandler) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	// Do nothing
	return nil
}

func (h *eventHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	// Print the DDL statement to the console
	log.Printf("DDL statement: %v", string(queryEvent.Query))

	return nil
}

func (h *eventHandler) OnRow(e *canal.RowsEvent) error {
	// Print the row event to the console
	log.Printf("Row event: %v", e)

	return nil
}

func (h *eventHandler) OnGTID(*replication.EventHeader, mysql.GTIDSet) error {
	// Do nothing
	return nil
}

func (h *eventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	// Do nothing
	return nil
}

func (h *eventHandler) OnXID(*replication.EventHeader, mysql.Position) error {
	// Do nothing
	return nil
}

func (h *eventHandler) OnUnmarshal(data []byte) (interface{}, error) {
	// Do nothing
	return nil, nil
}

func (h *eventHandler) OnRawEvent(event *replication.BinlogEvent) error {
	// Do nothing
	return nil
}
