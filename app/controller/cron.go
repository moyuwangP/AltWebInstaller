package controller

import (
	db2 "AltWebServer/app/model/db"
	"AltWebServer/app/service"
	"AltWebServer/app/util"
	"context"
	"sync"
	"time"
)

type node struct {
	scheduledAt time.Time
	udid        string
}

type queue struct {
	nodes []node
}

func (q *queue) enqueue(n node) {
	q.nodes = append(q.nodes, n)
}

func (q *queue) dequeue() node {
	n := q.peek()
	q.nodes = q.nodes[1:]
	return n
}

func (q *queue) peek() node {
	return q.nodes[0]
}

func (q *queue) empty() bool {
	return len(q.nodes) == 0
}

var refreshQueue queue
var mutex sync.Mutex

func ScheduleRefresh() {
	mutex.Lock()
	defer mutex.Unlock()

	refreshQueue = queue{}
	devices, err := service.Device.ListPairedDevices()
	if err != nil {
		util.LogErrorf("unable to find paired devices; %s", err.Error())
		return
	}

	for _, device := range devices {
		refreshQueue.enqueue(node{
			scheduledAt: time.Now(),
			udid:        device.UDID,
		})
	}
}

func DoRefresh() {
	mutex.Lock()
	defer mutex.Unlock()

	for !refreshQueue.empty() {
		if refreshQueue.peek().scheduledAt.Unix() > time.Now().Unix() {
			break
		}

		udid := refreshQueue.dequeue().udid
		if err := refreshDevice(udid); err != nil {
			util.LogErrorf("refresh %s failed: %s", udid, err.Error())
			refreshQueue.enqueue(node{
				scheduledAt: time.Now().Add(10 * time.Minute),
				udid:        udid,
			})
		}
	}
}

func refreshDevice(udid string) error {
	ctx := context.Background()
	packages, err := service.Device.AppsInstalled(ctx, udid)
	if err != nil {
		return err
	}

	for _, p := range packages {
		if p.RefreshedAt.Add(12*time.Hour).Unix() <= time.Now().Unix() {
			if err = refreshPackage(ctx, udid, p.IPAHash); err != nil {
				return err
			}
		}
	}

	return nil
}

func refreshPackage(ctx context.Context, udid string, hash string) error {
	installation := db2.Installation{}
	if err := util.DB().
		Where("udid", udid).Where("md5", hash).
		First(&installation).Error; err != nil {
		return err
	}
	return service.Device.InstallIPA(ctx, udid, hash, installation.RemovePlugIns)
}
