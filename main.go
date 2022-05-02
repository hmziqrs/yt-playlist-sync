package main

import (
	"log"
	"sync"
)

func main() {
	err := InitDirs()
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	youtubeClient, err := YoutubeClient()

	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	items, total := GetPlayList(youtubeClient)

	db, err := SavePlaylistToJSON(items, total)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	driveClient, err := DriveClient()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	InitRootDir(driveClient)

	total = 1

	GetDownloadedFiles(driveClient)

	batches, perBatch, batchesSize := CalculateBatches(total)
	var wg sync.WaitGroup
	wg.Add(batchesSize)

	for index, value := range batches {
		multiplier := (index * perBatch)
		start := multiplier + 1
		end := multiplier + value
		go func(index int) {
			WorkerShell(index, start, end, db, driveClient)
			wg.Done()
		}(index)
	}

	wg.Wait()
	SyncDB(db)
}
